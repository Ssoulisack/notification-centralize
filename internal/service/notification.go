package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/notification-center/internal/job"
	"github.com/your-org/notification-center/internal/model"
	"github.com/your-org/notification-center/internal/store"
	"github.com/your-org/notification-center/internal/template"
)

// NotificationService handles the business logic for sending notifications.
type NotificationService struct {
	notifStore  store.NotificationStore
	deviceStore store.DeviceStore
	prefStore   store.PreferenceStore
	tmplStore   store.TemplateStore
	queue       job.Queue
	tmplEngine  *template.Engine
	logger      *slog.Logger
}

// NewNotificationService creates a new service instance.
func NewNotificationService(
	notifStore store.NotificationStore,
	deviceStore store.DeviceStore,
	prefStore store.PreferenceStore,
	tmplStore store.TemplateStore,
	q job.Queue,
	tmplEngine *template.Engine,
	logger *slog.Logger,
) *NotificationService {
	return &NotificationService{
		notifStore:  notifStore,
		deviceStore: deviceStore,
		prefStore:   prefStore,
		tmplStore:   tmplStore,
		queue:       q,
		tmplEngine:  tmplEngine,
		logger:      logger,
	}
}

// Send handles both direct and event-based notification requests.
func (s *NotificationService) Send(ctx context.Context, projectID string, req *model.SendRequest) (*model.SendResponse, error) {
	// Event-based mode: resolve user → channels → recipients
	if req.UserID != "" && req.Event != "" {
		return s.sendByEvent(ctx, projectID, req)
	}

	// Direct mode: send to specific channel + recipient
	if req.Channel != "" && req.Recipient != "" {
		return s.sendDirect(ctx, projectID, req)
	}

	return nil, fmt.Errorf("provide either (channel + recipient) or (user_id + event)")
}

// sendDirect sends a single notification to a specific channel + recipient.
func (s *NotificationService) sendDirect(ctx context.Context, projectID string, req *model.SendRequest) (*model.SendResponse, error) {
	n := &model.Notification{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		UserID:    req.UserID,
		Channel:   req.Channel,
		Recipient: req.Recipient,
		Subject:   req.Subject,
		Body:      req.Body,
		Priority:  req.Priority,
		Status:    model.StatusQueued,
		Metadata:  req.Metadata,
		CreatedAt: time.Now(),
	}

	if n.Priority == "" {
		n.Priority = model.PriorityNormal
	}

	// Render template if provided
	if req.TemplateID != "" {
		if err := s.applyTemplate(ctx, projectID, n, req.TemplateID, req.Data); err != nil {
			return nil, fmt.Errorf("apply template: %w", err)
		}
	}

	// Persist
	if err := s.notifStore.Create(ctx, n); err != nil {
		return nil, fmt.Errorf("persist notification: %w", err)
	}

	// Enqueue
	if err := s.queue.Publish(n); err != nil {
		return nil, fmt.Errorf("enqueue notification: %w", err)
	}

	s.logger.Info("notification queued",
		"id", n.ID, "project_id", projectID, "channel", n.Channel, "recipient", n.Recipient)

	return &model.SendResponse{
		ID:     n.ID,
		Status: n.Status,
	}, nil
}

// sendByEvent resolves user preferences and fans out to enabled channels.
func (s *NotificationService) sendByEvent(ctx context.Context, projectID string, req *model.SendRequest) (*model.SendResponse, error) {
	// Check user preferences
	pref, err := s.prefStore.Get(ctx, projectID, req.UserID)
	if err != nil {
		s.logger.Warn("no preferences found, using defaults", "user_id", req.UserID)
		pref = &model.UserPreference{
			ProjectID:       projectID,
			UserID:          req.UserID,
			EnabledChannels: []model.Channel{model.ChannelEmail},
		}
	}

	// Check opt-out
	if pref.IsEventOptedOut(req.Event) {
		return &model.SendResponse{
			ID:      "",
			Status:  model.StatusCancelled,
			Message: "user opted out of this event",
		}, nil
	}

	// Fan out to each enabled channel
	var lastID string
	for _, ch := range pref.EnabledChannels {
		recipient, err := s.resolveRecipient(ctx, projectID, req.UserID, ch)
		if err != nil || recipient == "" {
			s.logger.Warn("could not resolve recipient",
				"user_id", req.UserID, "channel", ch, "error", err)
			continue
		}

		channelReq := &model.SendRequest{
			Channel:    ch,
			Recipient:  recipient,
			UserID:     req.UserID,
			Subject:    req.Subject,
			Body:       req.Body,
			TemplateID: req.TemplateID,
			Data:       req.Data,
			Priority:   req.Priority,
			Metadata:   req.Metadata,
		}

		resp, err := s.sendDirect(ctx, projectID, channelReq)
		if err != nil {
			s.logger.Error("failed to queue for channel",
				"channel", ch, "error", err)
			continue
		}
		lastID = resp.ID
	}

	if lastID == "" {
		return nil, fmt.Errorf("could not send to any channel for user %s", req.UserID)
	}

	return &model.SendResponse{
		ID:      lastID,
		Status:  model.StatusQueued,
		Message: fmt.Sprintf("queued for %d channels", len(pref.EnabledChannels)),
	}, nil
}

// resolveRecipient looks up the address for a user + channel combination.
func (s *NotificationService) resolveRecipient(ctx context.Context, projectID, userID string, ch model.Channel) (string, error) {
	switch ch {
	case model.ChannelPush:
		devices, err := s.deviceStore.GetByUser(ctx, projectID, userID)
		if err != nil || len(devices) == 0 {
			return "", fmt.Errorf("no device tokens for user %s", userID)
		}
		return devices[0].Token, nil // Send to most recent device

	case model.ChannelEmail, model.ChannelSMS, model.ChannelSlack, model.ChannelTelegram, model.ChannelLine:
		// TODO: look up from user service or contacts table
		// For now, check metadata or return empty
		return "", fmt.Errorf("recipient lookup not implemented for channel %s", ch)

	default:
		return "", fmt.Errorf("unknown channel: %s", ch)
	}
}

// applyTemplate renders a template and populates the notification fields.
func (s *NotificationService) applyTemplate(ctx context.Context, projectID string, n *model.Notification, templateID string, data map[string]string) error {
	tmpl, err := s.tmplStore.GetByID(ctx, projectID, templateID)
	if err != nil {
		return err
	}

	subject, err := s.tmplEngine.RenderSubject(tmpl, data)
	if err != nil {
		return err
	}
	n.Subject = subject

	body, err := s.tmplEngine.RenderBody(tmpl, data)
	if err != nil {
		return err
	}
	n.Body = body
	n.TemplateID = templateID

	return nil
}

// GetStatus returns the current status of a notification.
func (s *NotificationService) GetStatus(ctx context.Context, projectID, id string) (*model.Notification, error) {
	return s.notifStore.GetByID(ctx, projectID, id)
}

// ListByUser returns notification history for a user.
func (s *NotificationService) ListByUser(ctx context.Context, projectID, userID string, limit, offset int) ([]*model.Notification, error) {
	return s.notifStore.ListByUser(ctx, projectID, userID, limit, offset)
}

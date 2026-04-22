package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"github.com/your-org/notification-center/bootstrap/messaging"
	"github.com/your-org/notification-center/internal/data/repository"
	"github.com/your-org/notification-center/internal/domain/constants"
	"github.com/your-org/notification-center/internal/domain/dto"
	"github.com/your-org/notification-center/internal/domain/mapper"
	"github.com/your-org/notification-center/internal/domain/models"
)

type NotificationService interface {
	Send(ctx context.Context, projectID, senderID uuid.UUID, req *dto.SendRequest) (*dto.NotificationDTO, error)
	GetNotification(ctx context.Context, projectID, notificationID uuid.UUID) (*dto.NotificationDTO, error)
	ListNotifications(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*dto.NotificationDTO, int64, error)
	MarkAsRead(ctx context.Context, notificationID, userID uuid.UUID) error
	GetUnreadCount(ctx context.Context, projectID, userID uuid.UUID) (int64, error)
	GetUserInbox(ctx context.Context, projectID, userID uuid.UUID, limit, offset int) ([]*dto.NotificationDTO, int64, error)
}

type notificationService struct {
	repo     repository.NotificationRepository
	rabbitmq *messaging.Client
	logger   *slog.Logger
}

func NewNotificationService(repo repository.NotificationRepository, rabbitmq *messaging.Client, logger *slog.Logger) NotificationService {
	return &notificationService{
		repo:     repo,
		rabbitmq: rabbitmq,
		logger:   logger,
	}
}

func (s *notificationService) Send(ctx context.Context, projectID, senderID uuid.UUID, req *dto.SendRequest) (*dto.NotificationDTO, error) {
	if req.Priority == "" {
		req.Priority = constants.PriorityNormal
	}

	n := &models.Notification{
		ID:          uuid.New(),
		ProjectID:   projectID,
		SenderID:    &senderID,
		TemplateID:  req.TemplateID,
		Title:       req.Title,
		Body:        req.Body,
		Data:        datatypes.JSON(req.Data),
		Priority:    req.Priority,
		ScheduledAt: req.ScheduledAt,
		ExpiresAt:   req.ExpiresAt,
		CreatedAt:   time.Now(),
	}

	var recipients []models.NotificationRecipient

	for _, rec := range req.Recipients {
		for _, channel := range rec.Channels {
			address, err := s.resolveAddress(ctx, rec.UserID, channel)
			if err != nil {
				s.logger.Warn("failed to resolve recipient address",
					"user_id", rec.UserID, "channel", channel, "error", err)
				continue
			}

			recipientID := uuid.New()
			recipients = append(recipients, models.NotificationRecipient{
				ID:               recipientID,
				NotificationID:   n.ID,
				UserID:           rec.UserID,
				Channel:          channel,
				RecipientAddress: address,
				Status:           constants.StatusPending,
				CreatedAt:        time.Now(),
			})
		}
	}

	if err := s.repo.CreateWithRecipients(ctx, n, recipients); err != nil {
		return nil, fmt.Errorf("failed to persist notification: %w", err)
	}

	for i := range recipients {
		msg := &models.QueueMessage{
			NotificationID: n.ID,
			RecipientID:    recipients[i].ID,
			Channel:        recipients[i].Channel,
			Recipient:      recipients[i].RecipientAddress,
			Title:          n.Title,
			Body:           n.Body,
			Data:           n.Data,
			Priority:       n.Priority,
		}
		if err := s.rabbitmq.Publish(ctx, string(recipients[i].Channel), msg); err != nil {
			s.logger.Error("failed to queue message",
				"notification_id", n.ID, "channel", recipients[i].Channel, "error", err)
		}
	}

	s.logger.Info("notification sent",
		"notification_id", n.ID, "project_id", projectID, "recipient_count", len(recipients))

	n.Recipients = recipients
	return mapper.ToNotificationDTO(n), nil
}

func (s *notificationService) GetNotification(ctx context.Context, projectID, notificationID uuid.UUID) (*dto.NotificationDTO, error) {
	n, err := s.repo.GetByID(ctx, projectID, notificationID)
	if err != nil {
		return nil, err
	}
	return mapper.ToNotificationDTO(n), nil
}

func (s *notificationService) ListNotifications(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*dto.NotificationDTO, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	notifications, total, err := s.repo.List(ctx, projectID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*dto.NotificationDTO, len(notifications))
	for i := range notifications {
		result[i] = mapper.ToNotificationDTO(&notifications[i])
	}
	return result, total, nil
}

func (s *notificationService) MarkAsRead(ctx context.Context, notificationID, userID uuid.UUID) error {
	return s.repo.MarkAsRead(ctx, notificationID, userID)
}

func (s *notificationService) GetUnreadCount(ctx context.Context, projectID, userID uuid.UUID) (int64, error) {
	return s.repo.GetUnreadCount(ctx, projectID, userID)
}

func (s *notificationService) GetUserInbox(ctx context.Context, projectID, userID uuid.UUID, limit, offset int) ([]*dto.NotificationDTO, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	notifications, total, err := s.repo.GetUserInbox(ctx, projectID, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*dto.NotificationDTO, len(notifications))
	for i := range notifications {
		result[i] = mapper.ToNotificationDTO(&notifications[i])
	}
	return result, total, nil
}

// resolveAddress determines the delivery address for a user on a given channel.
func (s *notificationService) resolveAddress(ctx context.Context, userID uuid.UUID, channel constants.Channel) (string, error) {
	switch channel {
	case constants.ChannelEmail:
		return s.repo.GetUserEmail(ctx, userID)
	case constants.ChannelPush:
		return s.repo.GetActiveDeviceToken(ctx, userID)
	case constants.ChannelInApp:
		return userID.String(), nil
	default:
		return "", fmt.Errorf("unsupported channel: %s", channel)
	}
}

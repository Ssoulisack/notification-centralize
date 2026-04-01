package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-org/notification-center/bootstrap/messaging"
	"github.com/your-org/notification-center/internal/models"
)

// NotificationService handles notification business logic.
type NotificationService struct {
	db       *pgxpool.Pool
	rabbitmq *messaging.Client
	logger   *slog.Logger
}

// NewNotificationService creates a new notification service.
func NewNotificationService(db *pgxpool.Pool, rabbitmq *messaging.Client, logger *slog.Logger) *NotificationService {
	return &NotificationService{
		db:       db,
		rabbitmq: rabbitmq,
		logger:   logger,
	}
}

// SendRequest represents a notification send request.
type SendRequest struct {
	TemplateID  *uuid.UUID      `json:"template_id,omitempty"`
	Title       string          `json:"title" binding:"required"`
	Body        string          `json:"body" binding:"required"`
	Data        json.RawMessage `json:"data,omitempty"`
	Priority    models.Priority `json:"priority,omitempty"`
	Recipients  []RecipientReq  `json:"recipients" binding:"required,min=1"`
	ScheduledAt *time.Time      `json:"scheduled_at,omitempty"`
	ExpiresAt   *time.Time      `json:"expires_at,omitempty"`
}

// RecipientReq represents a recipient in the send request.
type RecipientReq struct {
	UserID   uuid.UUID        `json:"user_id" binding:"required"`
	Channels []models.Channel `json:"channels" binding:"required,min=1"`
}

// Send creates and queues a notification for delivery.
func (s *NotificationService) Send(ctx context.Context, projectID, senderID uuid.UUID, req *SendRequest) (*models.Notification, error) {
	// Set defaults
	if req.Priority == "" {
		req.Priority = models.PriorityNormal
	}

	// Create notification
	notification := &models.Notification{
		ID:          uuid.New(),
		ProjectID:   projectID,
		SenderID:    &senderID,
		TemplateID:  req.TemplateID,
		Title:       req.Title,
		Body:        req.Body,
		Data:        req.Data,
		Priority:    req.Priority,
		ScheduledAt: req.ScheduledAt,
		ExpiresAt:   req.ExpiresAt,
		CreatedAt:   time.Now(),
	}

	// Begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert notification
	insertNotificationQuery := `
		INSERT INTO notifications (
			id, project_id, template_id, sender_id, title, body, data,
			priority, scheduled_at, expires_at, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = tx.Exec(ctx, insertNotificationQuery,
		notification.ID, notification.ProjectID, notification.TemplateID,
		notification.SenderID, notification.Title, notification.Body,
		notification.Data, notification.Priority, notification.ScheduledAt,
		notification.ExpiresAt, notification.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert notification: %w", err)
	}

	// Insert recipients
	insertRecipientQuery := `
		INSERT INTO notification_recipients (
			id, notification_id, user_id, channel, recipient_address, status, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	for _, rec := range req.Recipients {
		for _, channel := range rec.Channels {
			recipientID := uuid.New()

			// Get recipient address based on channel
			address, err := s.getRecipientAddress(ctx, rec.UserID, channel)
			if err != nil {
				s.logger.Warn("failed to get recipient address",
					"user_id", rec.UserID,
					"channel", channel,
					"error", err,
				)
				continue
			}

			_, err = tx.Exec(ctx, insertRecipientQuery,
				recipientID, notification.ID, rec.UserID, channel,
				address, models.StatusPending, time.Now(),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to insert recipient: %w", err)
			}

			// Queue message for delivery
			queueMsg := &models.QueueMessage{
				NotificationID: notification.ID,
				RecipientID:    recipientID,
				Channel:        channel,
				Recipient:      address,
				Title:          notification.Title,
				Body:           notification.Body,
				Data:           notification.Data,
				Priority:       notification.Priority,
				RetryCount:     0,
			}

			routingKey := string(channel)
			if err := s.rabbitmq.Publish(ctx, routingKey, queueMsg); err != nil {
				s.logger.Error("failed to queue message",
					"notification_id", notification.ID,
					"channel", channel,
					"error", err,
				)
			}

			notification.Recipients = append(notification.Recipients, models.NotificationRecipient{
				ID:               recipientID,
				NotificationID:   notification.ID,
				UserID:           rec.UserID,
				Channel:          channel,
				RecipientAddress: address,
				Status:           models.StatusPending,
				CreatedAt:        time.Now(),
			})
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("notification sent",
		"notification_id", notification.ID,
		"project_id", projectID,
		"recipient_count", len(notification.Recipients),
	)

	return notification, nil
}

// getRecipientAddress looks up the recipient address for a channel.
func (s *NotificationService) getRecipientAddress(ctx context.Context, userID uuid.UUID, channel models.Channel) (string, error) {
	switch channel {
	case models.ChannelEmail:
		var email string
		err := s.db.QueryRow(ctx, "SELECT email FROM users WHERE id = $1", userID).Scan(&email)
		return email, err

	case models.ChannelPush:
		var token string
		err := s.db.QueryRow(ctx,
			"SELECT token FROM device_tokens WHERE user_id = $1 AND is_active = true LIMIT 1",
			userID,
		).Scan(&token)
		return token, err

	case models.ChannelInApp:
		return userID.String(), nil

	default:
		return "", fmt.Errorf("unsupported channel: %s", channel)
	}
}

// GetNotification retrieves a notification by ID.
func (s *NotificationService) GetNotification(ctx context.Context, projectID, notificationID uuid.UUID) (*models.Notification, error) {
	var notification models.Notification

	query := `
		SELECT id, project_id, template_id, sender_id, title, body, data,
		       priority, scheduled_at, expires_at, created_at
		FROM notifications
		WHERE id = $1 AND project_id = $2
	`

	err := s.db.QueryRow(ctx, query, notificationID, projectID).Scan(
		&notification.ID, &notification.ProjectID, &notification.TemplateID,
		&notification.SenderID, &notification.Title, &notification.Body,
		&notification.Data, &notification.Priority, &notification.ScheduledAt,
		&notification.ExpiresAt, &notification.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get recipients
	recipientQuery := `
		SELECT id, notification_id, user_id, channel, recipient_address, status,
		       sent_at, delivered_at, read_at, failed_at, error_message, retry_count, created_at
		FROM notification_recipients
		WHERE notification_id = $1
	`

	rows, err := s.db.Query(ctx, recipientQuery, notificationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var recipient models.NotificationRecipient
		err := rows.Scan(
			&recipient.ID, &recipient.NotificationID, &recipient.UserID,
			&recipient.Channel, &recipient.RecipientAddress, &recipient.Status,
			&recipient.SentAt, &recipient.DeliveredAt, &recipient.ReadAt,
			&recipient.FailedAt, &recipient.ErrorMessage, &recipient.RetryCount,
			&recipient.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notification.Recipients = append(notification.Recipients, recipient)
	}

	return &notification, nil
}

// ListNotifications retrieves notifications for a project.
func (s *NotificationService) ListNotifications(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]models.Notification, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM notifications WHERE project_id = $1`
	if err := s.db.QueryRow(ctx, countQuery, projectID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get notifications
	query := `
		SELECT id, project_id, template_id, sender_id, title, body, data,
		       priority, scheduled_at, expires_at, created_at
		FROM notifications
		WHERE project_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(ctx, query, projectID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		err := rows.Scan(
			&n.ID, &n.ProjectID, &n.TemplateID, &n.SenderID,
			&n.Title, &n.Body, &n.Data, &n.Priority,
			&n.ScheduledAt, &n.ExpiresAt, &n.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		notifications = append(notifications, n)
	}

	return notifications, total, nil
}

// MarkAsRead marks a notification as read for a user.
func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID, userID uuid.UUID) error {
	query := `
		UPDATE notification_recipients
		SET status = $1, read_at = NOW()
		WHERE notification_id = $2 AND user_id = $3 AND status != $1
	`

	_, err := s.db.Exec(ctx, query, models.StatusRead, notificationID, userID)
	return err
}

// GetUnreadCount returns the count of unread notifications for a user in a project.
func (s *NotificationService) GetUnreadCount(ctx context.Context, projectID, userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM notification_recipients nr
		JOIN notifications n ON n.id = nr.notification_id
		WHERE n.project_id = $1 AND nr.user_id = $2 AND nr.status != $3 AND nr.channel = $4
	`

	var count int
	err := s.db.QueryRow(ctx, query, projectID, userID, models.StatusRead, models.ChannelInApp).Scan(&count)
	return count, err
}

// GetUserInbox retrieves the in-app notification inbox for a user.
func (s *NotificationService) GetUserInbox(ctx context.Context, projectID, userID uuid.UUID, limit, offset int) ([]models.Notification, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM notification_recipients nr
		JOIN notifications n ON n.id = nr.notification_id
		WHERE n.project_id = $1 AND nr.user_id = $2 AND nr.channel = $3
	`
	if err := s.db.QueryRow(ctx, countQuery, projectID, userID, models.ChannelInApp).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get notifications
	query := `
		SELECT n.id, n.project_id, n.title, n.body, n.data, n.priority, n.created_at,
		       nr.id, nr.status, nr.read_at
		FROM notification_recipients nr
		JOIN notifications n ON n.id = nr.notification_id
		WHERE n.project_id = $1 AND nr.user_id = $2 AND nr.channel = $3
		ORDER BY n.created_at DESC
		LIMIT $4 OFFSET $5
	`

	rows, err := s.db.Query(ctx, query, projectID, userID, models.ChannelInApp, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		var r models.NotificationRecipient

		err := rows.Scan(
			&n.ID, &n.ProjectID, &n.Title, &n.Body, &n.Data, &n.Priority, &n.CreatedAt,
			&r.ID, &r.Status, &r.ReadAt,
		)
		if err != nil {
			return nil, 0, err
		}

		r.NotificationID = n.ID
		r.UserID = userID
		r.Channel = models.ChannelInApp
		n.Recipients = []models.NotificationRecipient{r}

		notifications = append(notifications, n)
	}

	return notifications, total, nil
}

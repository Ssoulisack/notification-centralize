package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/notification-center/internal/model"
)

type NotificationStore struct {
	pool *pgxpool.Pool
}

func NewNotificationStore(pool *pgxpool.Pool) *NotificationStore {
	return &NotificationStore{pool: pool}
}

func (s *NotificationStore) Create(ctx context.Context, n *model.Notification) error {
	query := `
		INSERT INTO notifications (id, project_id, user_id, channel, recipient, subject, body,
			template_id, priority, status, retry_count, metadata, error_message, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err := s.pool.Exec(ctx, query,
		n.ID, n.ProjectID, n.UserID, n.Channel, n.Recipient, n.Subject, n.Body,
		n.TemplateID, n.Priority, n.Status, n.RetryCount, n.Metadata,
		n.ErrorMessage, n.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func (s *NotificationStore) GetByID(ctx context.Context, projectID, id string) (*model.Notification, error) {
	query := `
		SELECT id, project_id, user_id, channel, recipient, subject, body, template_id,
			priority, status, retry_count, metadata, error_message, created_at, sent_at
		FROM notifications WHERE project_id = $1 AND id = $2`

	var n model.Notification
	err := s.pool.QueryRow(ctx, query, projectID, id).Scan(
		&n.ID, &n.ProjectID, &n.UserID, &n.Channel, &n.Recipient, &n.Subject, &n.Body,
		&n.TemplateID, &n.Priority, &n.Status, &n.RetryCount, &n.Metadata,
		&n.ErrorMessage, &n.CreatedAt, &n.SentAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get notification %s: %w", id, err)
	}
	return &n, nil
}

func (s *NotificationStore) UpdateStatus(ctx context.Context, projectID, id string, status model.Status, errMsg string) error {
	query := `UPDATE notifications SET status = $3, error_message = $4`

	if status == model.StatusSent {
		now := time.Now()
		query += `, sent_at = $5 WHERE project_id = $1 AND id = $2`
		_, err := s.pool.Exec(ctx, query, projectID, id, status, errMsg, now)
		return err
	}

	query += ` WHERE project_id = $1 AND id = $2`
	_, err := s.pool.Exec(ctx, query, projectID, id, status, errMsg)
	if err != nil {
		return fmt.Errorf("update notification status: %w", err)
	}
	return nil
}

func (s *NotificationStore) ListByUser(ctx context.Context, projectID, userID string, limit, offset int) ([]*model.Notification, error) {
	query := `
		SELECT id, project_id, user_id, channel, recipient, subject, body, template_id,
			priority, status, retry_count, metadata, error_message, created_at, sent_at
		FROM notifications
		WHERE project_id = $1 AND user_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := s.pool.Query(ctx, query, projectID, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*model.Notification
	for rows.Next() {
		var n model.Notification
		if err := rows.Scan(
			&n.ID, &n.ProjectID, &n.UserID, &n.Channel, &n.Recipient, &n.Subject, &n.Body,
			&n.TemplateID, &n.Priority, &n.Status, &n.RetryCount, &n.Metadata,
			&n.ErrorMessage, &n.CreatedAt, &n.SentAt,
		); err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}
		notifications = append(notifications, &n)
	}
	return notifications, nil
}

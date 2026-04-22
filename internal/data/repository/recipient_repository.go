package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/notification-center/internal/domain/constants"
	"github.com/your-org/notification-center/internal/domain/models"
)

type RecipientRepository interface {
	UpdateStatus(ctx context.Context, recipientID uuid.UUID, status constants.Status, errorMsg string) error
	IncrementRetryCount(ctx context.Context, recipientID uuid.UUID) error
}

type recipientRepository struct {
	db *gorm.DB
}

func NewRecipientRepository(db *gorm.DB) RecipientRepository {
	return &recipientRepository{db: db}
}

func (r *recipientRepository) UpdateStatus(ctx context.Context, recipientID uuid.UUID, status constants.Status, errorMsg string) error {
	updates := map[string]any{"status": status}

	switch status {
	case constants.StatusSent:
		now := time.Now()
		updates["sent_at"] = now
	case constants.StatusDelivered:
		now := time.Now()
		updates["delivered_at"] = now
	case constants.StatusFailed:
		now := time.Now()
		updates["failed_at"] = now
		updates["error_message"] = errorMsg
	}

	return r.db.WithContext(ctx).
		Model(&models.NotificationRecipient{}).
		Where("id = ?", recipientID).
		Updates(updates).Error
}

func (r *recipientRepository) IncrementRetryCount(ctx context.Context, recipientID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.NotificationRecipient{}).
		Where("id = ?", recipientID).
		UpdateColumn("retry_count", gorm.Expr("retry_count + 1")).Error
}

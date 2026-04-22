package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/notification-center/internal/domain/constants"
	"github.com/your-org/notification-center/internal/domain/models"
)

type NotificationRepository interface {
	CreateWithRecipients(ctx context.Context, n *models.Notification, recipients []models.NotificationRecipient) error
	GetByID(ctx context.Context, projectID, notificationID uuid.UUID) (*models.Notification, error)
	List(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]models.Notification, int64, error)
	GetUserInbox(ctx context.Context, projectID, userID uuid.UUID, limit, offset int) ([]models.Notification, int64, error)
	MarkAsRead(ctx context.Context, notificationID, userID uuid.UUID) error
	GetUnreadCount(ctx context.Context, projectID, userID uuid.UUID) (int64, error)
	GetUserEmail(ctx context.Context, userID uuid.UUID) (string, error)
	GetActiveDeviceToken(ctx context.Context, userID uuid.UUID) (string, error)
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

type notificationRepository struct {
	db *gorm.DB
}

func (r *notificationRepository) CreateWithRecipients(ctx context.Context, n *models.Notification, recipients []models.NotificationRecipient) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(n).Error; err != nil {
			return err
		}
		if len(recipients) > 0 {
			return tx.Create(&recipients).Error
		}
		return nil
	})
}

func (r *notificationRepository) GetByID(ctx context.Context, projectID, notificationID uuid.UUID) (*models.Notification, error) {
	var n models.Notification
	err := r.db.WithContext(ctx).
		Preload("Recipients").
		Where("id = ? AND project_id = ?", notificationID, projectID).
		First(&n).Error
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *notificationRepository) List(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]models.Notification, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&models.Notification{}).Where("project_id = ?", projectID).Count(&total)

	var notifications []models.Notification
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&notifications).Error

	return notifications, total, err
}

func (r *notificationRepository) GetUserInbox(ctx context.Context, projectID, userID uuid.UUID, limit, offset int) ([]models.Notification, int64, error) {
	var total int64
	r.db.WithContext(ctx).
		Model(&models.NotificationRecipient{}).
		Joins("JOIN notifications n ON n.id = notification_recipients.notification_id").
		Where("n.project_id = ? AND notification_recipients.user_id = ? AND notification_recipients.channel = ?",
			projectID, userID, constants.ChannelInApp).
		Count(&total)

	var recipients []models.NotificationRecipient
	err := r.db.WithContext(ctx).
		Preload("Notification").
		Joins("JOIN notifications n ON n.id = notification_recipients.notification_id").
		Where("n.project_id = ? AND notification_recipients.user_id = ? AND notification_recipients.channel = ?",
			projectID, userID, constants.ChannelInApp).
		Order("n.created_at DESC").
		Limit(limit).Offset(offset).
		Find(&recipients).Error
	if err != nil {
		return nil, 0, err
	}

	notifications := make([]models.Notification, 0, len(recipients))
	for i := range recipients {
		n := recipients[i].Notification
		n.Recipients = []models.NotificationRecipient{recipients[i]}
		notifications = append(notifications, n)
	}

	return notifications, total, nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, notificationID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.NotificationRecipient{}).
		Where("notification_id = ? AND user_id = ? AND status != ?", notificationID, userID, constants.StatusRead).
		Updates(map[string]any{"status": constants.StatusRead, "read_at": r.db.NowFunc()}).Error
}

func (r *notificationRepository) GetUnreadCount(ctx context.Context, projectID, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.NotificationRecipient{}).
		Joins("JOIN notifications n ON n.id = notification_recipients.notification_id").
		Where("n.project_id = ? AND notification_recipients.user_id = ? AND notification_recipients.status != ? AND notification_recipients.channel = ?",
			projectID, userID, constants.StatusRead, constants.ChannelInApp).
		Count(&count).Error
	return count, err
}

func (r *notificationRepository) GetUserEmail(ctx context.Context, userID uuid.UUID) (string, error) {
	var email string
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Select("email").
		Where("id = ?", userID).
		Scan(&email).Error
	return email, err
}

func (r *notificationRepository) GetActiveDeviceToken(ctx context.Context, userID uuid.UUID) (string, error) {
	var token string
	err := r.db.WithContext(ctx).
		Model(&models.DeviceToken{}).
		Select("token").
		Where("user_id = ? AND is_active = true", userID).
		Limit(1).
		Scan(&token).Error
	return token, err
}

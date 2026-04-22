package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/your-org/notification-center/internal/domain/constants"
)

type UserNotificationPreference struct {
	ID              uuid.UUID         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID          uuid.UUID         `gorm:"type:uuid;not null;index"`
	User            User              `gorm:"foreignKey:UserID"`
	ProjectID       uuid.UUID         `gorm:"type:uuid;not null"`
	Project         Project           `gorm:"foreignKey:ProjectID"`
	Channel         constants.Channel `gorm:"type:notification_channel;not null"`
	Enabled         bool              `gorm:"default:true"`
	QuietHoursStart string            `gorm:"type:time"`
	QuietHoursEnd   string            `gorm:"type:time"`
	Frequency       string            `gorm:"default:'instant';size:50"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (UserNotificationPreference) TableName() string { return "user_notification_preferences" }

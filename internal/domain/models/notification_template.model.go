package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"github.com/your-org/notification-center/internal/domain/constants"
)

type NotificationTemplate struct {
	ID              uuid.UUID           `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID       uuid.UUID           `gorm:"type:uuid;not null;index"`
	Project         Project             `gorm:"foreignKey:ProjectID"`
	Name            string              `gorm:"not null;size:255"`
	Slug            string              `gorm:"not null;size:255"`
	Channel         constants.Channel   `gorm:"type:notification_channel;not null"`
	SubjectTemplate string
	BodyTemplate    string              `gorm:"not null"`
	Variables       datatypes.JSON      `gorm:"type:jsonb;default:'[]'"`
	IsActive        bool                `gorm:"default:true"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (NotificationTemplate) TableName() string { return "notification_templates" }

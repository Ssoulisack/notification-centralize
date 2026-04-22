package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"github.com/your-org/notification-center/internal/domain/constants"
)

type WebhookEndpoint struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID uuid.UUID      `gorm:"type:uuid;not null;index"`
	Project   Project        `gorm:"foreignKey:ProjectID"`
	URL       string         `gorm:"not null"`
	Secret    string         `gorm:"not null;size:255"`
	Events    datatypes.JSON `gorm:"type:jsonb;default:'[\"notification.sent\",\"notification.delivered\",\"notification.failed\"]'"`
	IsActive  bool           `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (WebhookEndpoint) TableName() string { return "webhook_endpoints" }

type WebhookDelivery struct {
	ID             uuid.UUID               `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	EndpointID     uuid.UUID               `gorm:"type:uuid;not null;index"`
	Endpoint       WebhookEndpoint         `gorm:"foreignKey:EndpointID"`
	EventType      string                  `gorm:"not null;size:50"`
	Payload        datatypes.JSON          `gorm:"type:jsonb;not null"`
	ResponseStatus int
	ResponseBody   string
	Status         constants.WebhookStatus `gorm:"type:webhook_status;default:'pending';index"`
	Attempts       int                     `gorm:"default:0"`
	NextRetryAt    *time.Time
	CreatedAt      time.Time
	CompletedAt    *time.Time
}

func (WebhookDelivery) TableName() string { return "webhook_deliveries" }

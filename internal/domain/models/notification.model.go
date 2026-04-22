package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"github.com/your-org/notification-center/internal/domain/constants"
)

type Notification struct {
	ID          uuid.UUID             `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID   uuid.UUID             `gorm:"type:uuid;not null;index"`
	Project     Project               `gorm:"foreignKey:ProjectID"`
	TemplateID  *uuid.UUID            `gorm:"type:uuid"`
	Template    *NotificationTemplate `gorm:"foreignKey:TemplateID"`
	SenderID    *uuid.UUID            `gorm:"type:uuid"`
	Sender      *User                 `gorm:"foreignKey:SenderID"`
	Title       string                `gorm:"not null;size:500"`
	Body        string                `gorm:"not null"`
	Data        datatypes.JSON        `gorm:"type:jsonb;default:'{}'"`
	Priority    constants.Priority    `gorm:"type:notification_priority;default:'normal'"`
	ScheduledAt *time.Time
	ExpiresAt   *time.Time
	CreatedAt   time.Time

	Recipients []NotificationRecipient `gorm:"foreignKey:NotificationID"`
}

func (Notification) TableName() string { return "notifications" }

type NotificationRecipient struct {
	ID               uuid.UUID          `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	NotificationID   uuid.UUID          `gorm:"type:uuid;not null;index"`
	Notification     Notification       `gorm:"foreignKey:NotificationID"`
	UserID           uuid.UUID          `gorm:"type:uuid;not null;index"`
	User             User               `gorm:"foreignKey:UserID"`
	Channel          constants.Channel  `gorm:"type:notification_channel;not null"`
	RecipientAddress string             `gorm:"size:500"`
	Status           constants.Status   `gorm:"type:notification_status;default:'pending';index"`
	SentAt           *time.Time
	DeliveredAt      *time.Time
	ReadAt           *time.Time
	FailedAt         *time.Time
	ErrorMessage     string
	RetryCount       int            `gorm:"default:0"`
	Metadata         datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	CreatedAt        time.Time
}

func (NotificationRecipient) TableName() string { return "notification_recipients" }

type NotificationEvent struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	NotificationID uuid.UUID      `gorm:"type:uuid;not null;index"`
	Notification   Notification   `gorm:"foreignKey:NotificationID"`
	RecipientID    *uuid.UUID     `gorm:"type:uuid"`
	EventType      string         `gorm:"not null;size:50;index"`
	EventData      datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	CreatedAt      time.Time
}

func (NotificationEvent) TableName() string { return "notification_events" }

// QueueMessage is not a DB model — used for async notification processing.
type QueueMessage struct {
	NotificationID uuid.UUID          `json:"notification_id"`
	RecipientID    uuid.UUID          `json:"recipient_id"`
	Channel        constants.Channel  `json:"channel"`
	Recipient      string             `json:"recipient"`
	Title          string             `json:"title"`
	Body           string             `json:"body"`
	Data           any                `json:"data,omitempty"`
	Priority       constants.Priority `json:"priority"`
	RetryCount     int                `json:"retry_count"`
}

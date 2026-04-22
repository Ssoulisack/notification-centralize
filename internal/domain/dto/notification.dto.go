package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"github.com/your-org/notification-center/internal/domain/constants"
)

type NotificationDTO struct {
	ID          uuid.UUID                  `json:"id"`
	ProjectID   uuid.UUID                  `json:"project_id"`
	TemplateID  *uuid.UUID                 `json:"template_id,omitempty"`
	SenderID    *uuid.UUID                 `json:"sender_id,omitempty"`
	Title       string                     `json:"title"`
	Body        string                     `json:"body"`
	Data        datatypes.JSON             `json:"data,omitempty"`
	Priority    constants.Priority         `json:"priority"`
	ScheduledAt *time.Time                 `json:"scheduled_at,omitempty"`
	ExpiresAt   *time.Time                 `json:"expires_at,omitempty"`
	CreatedAt   time.Time                  `json:"created_at"`
	Recipients  []NotificationRecipientDTO `json:"recipients,omitempty"`
}

type NotificationRecipientDTO struct {
	ID               uuid.UUID         `json:"id"`
	NotificationID   uuid.UUID         `json:"notification_id"`
	UserID           uuid.UUID         `json:"user_id"`
	Channel          constants.Channel `json:"channel"`
	RecipientAddress string            `json:"recipient_address,omitempty"`
	Status           constants.Status  `json:"status"`
	SentAt           *time.Time        `json:"sent_at,omitempty"`
	DeliveredAt      *time.Time        `json:"delivered_at,omitempty"`
	ReadAt           *time.Time        `json:"read_at,omitempty"`
	FailedAt         *time.Time        `json:"failed_at,omitempty"`
	ErrorMessage     string            `json:"error_message,omitempty"`
	RetryCount       int               `json:"retry_count"`
	CreatedAt        time.Time         `json:"created_at"`
}

type SendRequest struct {
	TemplateID  *uuid.UUID         `json:"template_id,omitempty"`
	Title       string             `json:"title" binding:"required"`
	Body        string             `json:"body" binding:"required"`
	Data        json.RawMessage    `json:"data,omitempty"`
	Priority    constants.Priority `json:"priority,omitempty"`
	Recipients  []RecipientReq     `json:"recipients" binding:"required,min=1"`
	ScheduledAt *time.Time         `json:"scheduled_at,omitempty"`
	ExpiresAt   *time.Time         `json:"expires_at,omitempty"`
}

type RecipientReq struct {
	UserID   uuid.UUID           `json:"user_id" binding:"required"`
	Channels []constants.Channel `json:"channels" binding:"required,min=1"`
}

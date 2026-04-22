package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"github.com/your-org/notification-center/internal/domain/constants"
)

type WebhookEndpointDTO struct {
	ID        uuid.UUID      `json:"id"`
	ProjectID uuid.UUID      `json:"project_id"`
	URL       string         `json:"url"`
	Events    datatypes.JSON `json:"events"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type WebhookDeliveryDTO struct {
	ID             uuid.UUID               `json:"id"`
	EndpointID     uuid.UUID               `json:"endpoint_id"`
	EventType      string                  `json:"event_type"`
	Payload        datatypes.JSON          `json:"payload"`
	ResponseStatus int                     `json:"response_status,omitempty"`
	ResponseBody   string                  `json:"response_body,omitempty"`
	Status         constants.WebhookStatus `json:"status"`
	Attempts       int                     `json:"attempts"`
	NextRetryAt    *time.Time              `json:"next_retry_at,omitempty"`
	CreatedAt      time.Time               `json:"created_at"`
	CompletedAt    *time.Time              `json:"completed_at,omitempty"`
}

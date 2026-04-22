package dto

import (
	"time"

	"github.com/google/uuid"
)

type DeviceTokenDTO struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	ProjectID  uuid.UUID  `json:"project_id"`
	Platform   string     `json:"platform"`
	DeviceName string     `json:"device_name,omitempty"`
	AppVersion string     `json:"app_version,omitempty"`
	IsActive   bool       `json:"is_active"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

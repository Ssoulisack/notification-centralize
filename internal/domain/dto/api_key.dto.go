package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type APIKeyDTO struct {
	ID         uuid.UUID      `json:"id"`
	ProjectID  uuid.UUID      `json:"project_id"`
	UserID     uuid.UUID      `json:"user_id"`
	Name       string         `json:"name"`
	KeyPrefix  string         `json:"key_prefix"`
	Scopes     datatypes.JSON `json:"scopes"`
	LastUsedAt *time.Time     `json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time     `json:"expires_at,omitempty"`
	IsActive   bool           `json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	RevokedAt  *time.Time     `json:"revoked_at,omitempty"`
}

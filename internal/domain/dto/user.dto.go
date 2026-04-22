package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserDTO struct {
	ID          uuid.UUID  `json:"id"`
	KeycloakID  string     `json:"keycloak_id"`
	Email       string     `json:"email"`
	Username    string     `json:"username"`
	FirstName   string     `json:"first_name,omitempty"`
	LastName    string     `json:"last_name,omitempty"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

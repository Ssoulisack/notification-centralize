package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ProjectDTO struct {
	ID          uuid.UUID      `json:"id"`
	OwnerID     uuid.UUID      `json:"owner_id"`
	Owner       *UserDTO       `json:"owner,omitempty"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Slug        string         `json:"slug"`
	Settings    datatypes.JSON `json:"settings,omitempty"`
	IsActive    bool           `json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"github.com/your-org/notification-center/internal/domain/constants"
)

type NotificationTemplateDTO struct {
	ID              uuid.UUID         `json:"id"`
	ProjectID       uuid.UUID         `json:"project_id"`
	Name            string            `json:"name"`
	Slug            string            `json:"slug"`
	Channel         constants.Channel `json:"channel"`
	SubjectTemplate string            `json:"subject_template,omitempty"`
	BodyTemplate    string            `json:"body_template"`
	Variables       datatypes.JSON    `json:"variables,omitempty"`
	IsActive        bool              `json:"is_active"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

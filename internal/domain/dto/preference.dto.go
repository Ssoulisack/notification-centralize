package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/your-org/notification-center/internal/domain/constants"
)

type UserNotificationPreferenceDTO struct {
	ID              uuid.UUID         `json:"id"`
	UserID          uuid.UUID         `json:"user_id"`
	ProjectID       uuid.UUID         `json:"project_id"`
	Channel         constants.Channel `json:"channel"`
	Enabled         bool              `json:"enabled"`
	QuietHoursStart string            `json:"quiet_hours_start,omitempty"`
	QuietHoursEnd   string            `json:"quiet_hours_end,omitempty"`
	Frequency       string            `json:"frequency"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

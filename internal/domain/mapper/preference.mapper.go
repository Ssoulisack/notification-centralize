package mapper

import (
	"github.com/your-org/notification-center/internal/domain/dto"
	"github.com/your-org/notification-center/internal/domain/models"
)

func ToUserNotificationPreferenceDTO(m *models.UserNotificationPreference) *dto.UserNotificationPreferenceDTO {
	if m == nil {
		return nil
	}
	return &dto.UserNotificationPreferenceDTO{
		ID:              m.ID,
		UserID:          m.UserID,
		ProjectID:       m.ProjectID,
		Channel:         m.Channel,
		Enabled:         m.Enabled,
		QuietHoursStart: m.QuietHoursStart,
		QuietHoursEnd:   m.QuietHoursEnd,
		Frequency:       m.Frequency,
		UpdatedAt:       m.UpdatedAt,
	}
}

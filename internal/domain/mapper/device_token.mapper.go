package mapper

import (
	"github.com/your-org/notification-center/internal/domain/dto"
	"github.com/your-org/notification-center/internal/domain/models"
)

func ToDeviceTokenDTO(m *models.DeviceToken) *dto.DeviceTokenDTO {
	if m == nil {
		return nil
	}
	return &dto.DeviceTokenDTO{
		ID:         m.ID,
		UserID:     m.UserID,
		ProjectID:  m.ProjectID,
		Platform:   m.Platform,
		DeviceName: m.DeviceName,
		AppVersion: m.AppVersion,
		IsActive:   m.IsActive,
		LastUsedAt: m.LastUsedAt,
		CreatedAt:  m.CreatedAt,
	}
}


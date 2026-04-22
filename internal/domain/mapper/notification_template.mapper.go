package mapper

import (
	"github.com/your-org/notification-center/internal/domain/dto"
	"github.com/your-org/notification-center/internal/domain/models"
)

func ToNotificationTemplateDTO(m *models.NotificationTemplate) *dto.NotificationTemplateDTO {
	if m == nil {
		return nil
	}
	return &dto.NotificationTemplateDTO{
		ID:              m.ID,
		ProjectID:       m.ProjectID,
		Name:            m.Name,
		Slug:            m.Slug,
		Channel:         m.Channel,
		SubjectTemplate: m.SubjectTemplate,
		BodyTemplate:    m.BodyTemplate,
		Variables:       m.Variables,
		IsActive:        m.IsActive,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

package mapper

import (
	"github.com/your-org/notification-center/internal/domain/dto"
	"github.com/your-org/notification-center/internal/domain/models"
)

func ToNotificationDTO(m *models.Notification) *dto.NotificationDTO {
	if m == nil {
		return nil
	}
	d := &dto.NotificationDTO{
		ID:          m.ID,
		ProjectID:   m.ProjectID,
		TemplateID:  m.TemplateID,
		SenderID:    m.SenderID,
		Title:       m.Title,
		Body:        m.Body,
		Data:        m.Data,
		Priority:    m.Priority,
		ScheduledAt: m.ScheduledAt,
		ExpiresAt:   m.ExpiresAt,
		CreatedAt:   m.CreatedAt,
	}
	for i := range m.Recipients {
		d.Recipients = append(d.Recipients, *ToNotificationRecipientDTO(&m.Recipients[i]))
	}
	return d
}

func ToNotificationRecipientDTO(m *models.NotificationRecipient) *dto.NotificationRecipientDTO {
	if m == nil {
		return nil
	}
	return &dto.NotificationRecipientDTO{
		ID:               m.ID,
		NotificationID:   m.NotificationID,
		UserID:           m.UserID,
		Channel:          m.Channel,
		RecipientAddress: m.RecipientAddress,
		Status:           m.Status,
		SentAt:           m.SentAt,
		DeliveredAt:      m.DeliveredAt,
		ReadAt:           m.ReadAt,
		FailedAt:         m.FailedAt,
		ErrorMessage:     m.ErrorMessage,
		RetryCount:       m.RetryCount,
		CreatedAt:        m.CreatedAt,
	}
}

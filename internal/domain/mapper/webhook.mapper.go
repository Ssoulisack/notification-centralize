package mapper

import (
	"github.com/your-org/notification-center/internal/domain/dto"
	"github.com/your-org/notification-center/internal/domain/models"
)

func ToWebhookEndpointDTO(m *models.WebhookEndpoint) *dto.WebhookEndpointDTO {
	if m == nil {
		return nil
	}
	return &dto.WebhookEndpointDTO{
		ID:        m.ID,
		ProjectID: m.ProjectID,
		URL:       m.URL,
		Events:    m.Events,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToWebhookDeliveryDTO(m *models.WebhookDelivery) *dto.WebhookDeliveryDTO {
	if m == nil {
		return nil
	}
	return &dto.WebhookDeliveryDTO{
		ID:             m.ID,
		EndpointID:     m.EndpointID,
		EventType:      m.EventType,
		Payload:        m.Payload,
		ResponseStatus: m.ResponseStatus,
		ResponseBody:   m.ResponseBody,
		Status:         m.Status,
		Attempts:       m.Attempts,
		NextRetryAt:    m.NextRetryAt,
		CreatedAt:      m.CreatedAt,
		CompletedAt:    m.CompletedAt,
	}
}

package mapper

import (
	"github.com/your-org/notification-center/internal/domain/dto"
	"github.com/your-org/notification-center/internal/domain/models"
)

func ToAPIKeyDTO(m *models.APIKey) *dto.APIKeyDTO {
	if m == nil {
		return nil
	}
	return &dto.APIKeyDTO{
		ID:         m.ID,
		ProjectID:  m.ProjectID,
		UserID:     m.UserID,
		Name:       m.Name,
		KeyPrefix:  m.KeyPrefix,
		Scopes:     m.Scopes,
		LastUsedAt: m.LastUsedAt,
		ExpiresAt:  m.ExpiresAt,
		IsActive:   m.IsActive,
		CreatedAt:  m.CreatedAt,
		RevokedAt:  m.RevokedAt,
	}
}

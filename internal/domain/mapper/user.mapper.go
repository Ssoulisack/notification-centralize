package mapper

import (
	"github.com/your-org/notification-center/internal/domain/dto"
	"github.com/your-org/notification-center/internal/domain/models"
)

func ToUserDTO(m *models.User) *dto.UserDTO {
	if m == nil {
		return nil
	}
	return &dto.UserDTO{
		ID:          m.ID,
		KeycloakID:  m.KeycloakID,
		Email:       m.Email,
		Username:    m.Username,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		IsActive:    m.IsActive,
		LastLoginAt: m.LastLoginAt,
		CreatedAt:   m.CreatedAt,
	}
}

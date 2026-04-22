package mapper

import (
	"github.com/google/uuid"

	"github.com/your-org/notification-center/internal/domain/dto"
	"github.com/your-org/notification-center/internal/domain/models"
)

func ToProjectDTO(m *models.Project) *dto.ProjectDTO {
	if m == nil {
		return nil
	}
	d := &dto.ProjectDTO{
		ID:          m.ID,
		OwnerID:     m.OwnerID,
		Name:        m.Name,
		Description: m.Description,
		Slug:        m.Slug,
		Settings:    m.Settings,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	if m.Owner.ID != uuid.Nil {
		d.Owner = ToUserDTO(&m.Owner)
	}
	return d
}

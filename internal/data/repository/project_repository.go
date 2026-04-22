package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/notification-center/internal/domain/models"
)

var ErrNotFound = gorm.ErrRecordNotFound

type projectRepository struct {
	db *gorm.DB
}

type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
	GetBySlug(ctx context.Context, slug string) (*models.Project, error)
	List(ctx context.Context, limit, offset int) ([]models.Project, int64, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID, limit, offset int) ([]models.Project, int64, error)
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.db.WithContext(ctx).Preload("Owner").Where("id = ?", id).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) GetBySlug(ctx context.Context, slug string) (*models.Project, error) {
	var project models.Project
	err := r.db.WithContext(ctx).Preload("Owner").Where("slug = ?", slug).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) List(ctx context.Context, limit, offset int) ([]models.Project, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&models.Project{}).Where("is_active = true").Count(&total)

	var projects []models.Project
	err := r.db.WithContext(ctx).
		Where("is_active = true").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&projects).Error

	return projects, total, err
}

func (r *projectRepository) ListByOwner(ctx context.Context, ownerID uuid.UUID, limit, offset int) ([]models.Project, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&models.Project{}).Where("owner_id = ? AND is_active = true", ownerID).Count(&total)

	var projects []models.Project
	err := r.db.WithContext(ctx).
		Where("owner_id = ? AND is_active = true", ownerID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&projects).Error

	return projects, total, err
}

func (r *projectRepository) Update(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Model(&models.Project{}).
		Where("id = ?", id).
		Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

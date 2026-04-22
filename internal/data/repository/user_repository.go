package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/notification-center/internal/domain/models"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByKeycloakID(ctx context.Context, keycloakID string) (*models.User, error)
	Upsert(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByKeycloakID(ctx context.Context, keycloakID string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("keycloak_id = ?", keycloakID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Upsert(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).
		Where(models.User{KeycloakID: user.KeycloakID}).
		Assign(models.User{
			Email:       user.Email,
			Username:    user.Username,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			LastLoginAt: user.LastLoginAt,
		}).
		FirstOrCreate(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

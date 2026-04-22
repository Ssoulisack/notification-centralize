package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/your-org/notification-center/internal/data/repository"
	"github.com/your-org/notification-center/internal/domain/dto"
	"github.com/your-org/notification-center/internal/domain/mapper"
	"github.com/your-org/notification-center/internal/domain/models"
	"github.com/your-org/notification-center/pkg/middleware"
)

type UserSyncService interface {
	Sync(ctx context.Context, claims *middleware.KeycloakClaims) (*dto.UserDTO, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.UserDTO, error)
	GetUserDTOByKeycloakID(ctx context.Context, keycloakID string) (*dto.UserDTO, error)
	GetUserByKeycloakID(ctx context.Context, keycloakID string) (*models.User, error)
}

type userSyncService struct {
	repo   repository.UserRepository
	logger *slog.Logger
}

func NewUserSyncService(repo repository.UserRepository, logger *slog.Logger) UserSyncService {
	return &userSyncService{repo: repo, logger: logger}
}

// Sync upserts a user from Keycloak claims and returns a DTO.
func (s *userSyncService) Sync(ctx context.Context, claims *middleware.KeycloakClaims) (*dto.UserDTO, error) {
	now := time.Now()
	user := &models.User{
		KeycloakID:  claims.Subject,
		Email:       claims.Email,
		Username:    claims.PreferredUsername,
		FirstName:   claims.GivenName,
		LastName:    claims.FamilyName,
		IsActive:    true,
		LastLoginAt: &now,
	}

	if err := s.repo.Upsert(ctx, user); err != nil {
		s.logger.Error("failed to sync user", "keycloak_id", claims.Subject, "error", err)
		return nil, err
	}

	return mapper.ToUserDTO(user), nil
}

// GetByID returns a user DTO by database ID.
func (s *userSyncService) GetByID(ctx context.Context, id uuid.UUID) (*dto.UserDTO, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.ToUserDTO(user), nil
}

// GetUserDTOByKeycloakID returns a user DTO by Keycloak subject ID.
func (s *userSyncService) GetUserDTOByKeycloakID(ctx context.Context, keycloakID string) (*dto.UserDTO, error) {
	user, err := s.repo.GetByKeycloakID(ctx, keycloakID)
	if err != nil {
		return nil, err
	}
	return mapper.ToUserDTO(user), nil
}

// GetUserByKeycloakID returns the user by Keycloak subject ID.
func (s *userSyncService) GetUserByKeycloakID(ctx context.Context, keycloakID string) (*models.User, error) {
	return s.repo.GetByKeycloakID(ctx, keycloakID)
}

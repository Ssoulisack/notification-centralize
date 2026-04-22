package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/your-org/notification-center/internal/data/repository"
	"github.com/your-org/notification-center/internal/domain/dto"
	"github.com/your-org/notification-center/internal/domain/mapper"
	"github.com/your-org/notification-center/internal/domain/models"
)

type APIKeyService interface {
	List(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*dto.APIKeyDTO, int, error)
	Create(ctx context.Context, projectID uuid.UUID, keycloakID, name string) (*dto.APIKeyDTO, string, error)
	Revoke(ctx context.Context, keyID uuid.UUID) error
}

type apiKeyService struct {
	apiKeyRepo  repository.APIKeyRepository
	userService UserSyncService
}

func NewAPIKeyService(apiKeyRepo repository.APIKeyRepository, userService UserSyncService) APIKeyService {
	return &apiKeyService{
		apiKeyRepo:  apiKeyRepo,
		userService: userService,
	}
}

func (s *apiKeyService) List(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*dto.APIKeyDTO, int, error) {
	keys, total, err := s.apiKeyRepo.ListByProject(ctx, projectID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*dto.APIKeyDTO, len(keys))
	for i := range keys {
		result[i] = mapper.ToAPIKeyDTO(&keys[i])
	}
	return result, total, nil
}

func (s *apiKeyService) Create(ctx context.Context, projectID uuid.UUID, keycloakID, name string) (*dto.APIKeyDTO, string, error) {
	user, err := s.userService.GetUserByKeycloakID(ctx, keycloakID)
	if err != nil {
		return nil, "", err
	}

	apiKey := &models.APIKey{
		ProjectID: projectID,
		UserID:    user.ID,
		Name:      name,
	}

	fullKey, err := s.apiKeyRepo.Create(ctx, apiKey)
	if err != nil {
		return nil, "", err
	}

	return mapper.ToAPIKeyDTO(apiKey), fullKey, nil
}

func (s *apiKeyService) Revoke(ctx context.Context, keyID uuid.UUID) error {
	return s.apiKeyRepo.Revoke(ctx, keyID)
}

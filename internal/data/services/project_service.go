package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"

	"github.com/your-org/notification-center/internal/data/repository"
	"github.com/your-org/notification-center/internal/domain/dto"
	"github.com/your-org/notification-center/internal/domain/mapper"
	"github.com/your-org/notification-center/internal/domain/models"
)

type ProjectService interface {
	List(ctx context.Context, keycloakID string, limit, offset int) ([]*dto.ProjectDTO, int64, error)
	Create(ctx context.Context, keycloakID, name, description, slug string) (*dto.ProjectDTO, string, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.ProjectDTO, error)
	Update(ctx context.Context, id uuid.UUID, name, description string) (*dto.ProjectDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type projectService struct {
	projectRepo repository.ProjectRepository
	apiKeyRepo  repository.APIKeyRepository
	userService UserSyncService
}

func NewProjectService(
	projectRepo repository.ProjectRepository,
	apiKeyRepo repository.APIKeyRepository,
	userService UserSyncService,
) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		apiKeyRepo:  apiKeyRepo,
		userService: userService,
	}
}

func (s *projectService) List(ctx context.Context, keycloakID string, limit, offset int) ([]*dto.ProjectDTO, int64, error) {
	user, err := s.userService.GetUserByKeycloakID(ctx, keycloakID)
	if err != nil {
		return nil, 0, fmt.Errorf("user not found")
	}

	projects, total, err := s.projectRepo.ListByOwner(ctx, user.ID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*dto.ProjectDTO, len(projects))
	for i := range projects {
		result[i] = mapper.ToProjectDTO(&projects[i])
	}
	return result, total, nil
}

func (s *projectService) Create(ctx context.Context, keycloakID, name, description, slug string) (*dto.ProjectDTO, string, error) {
	user, err := s.userService.GetUserByKeycloakID(ctx, keycloakID)
	if err != nil {
		return nil, "", fmt.Errorf("user not found")
	}

	if slug == "" {
		slug = generateSlug(name)
	}

	project := &models.Project{
		OwnerID:     user.ID,
		Name:        name,
		Description: description,
		Slug:        slug,
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, "", err
	}

	apiKey := &models.APIKey{
		ProjectID: project.ID,
		UserID:    user.ID,
		Name:      "Default Key",
	}

	fullKey, err := s.apiKeyRepo.Create(ctx, apiKey)
	if err != nil {
		return nil, "", err
	}

	return mapper.ToProjectDTO(project), fullKey, nil
}

func (s *projectService) GetByID(ctx context.Context, id uuid.UUID) (*dto.ProjectDTO, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.ToProjectDTO(project), nil
}

func (s *projectService) Update(ctx context.Context, id uuid.UUID, name, description string) (*dto.ProjectDTO, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		project.Name = name
	}
	if description != "" {
		project.Description = description
	}

	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}

	return mapper.ToProjectDTO(project), nil
}

func (s *projectService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.projectRepo.Delete(ctx, id)
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(slug, "")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = uuid.New().String()[:8]
	}
	return slug
}

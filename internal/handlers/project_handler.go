package handlers

import (
	"log/slog"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/your-org/notification-center/internal/middleware"
	"github.com/your-org/notification-center/internal/models"
	"github.com/your-org/notification-center/internal/repository"
	"github.com/your-org/notification-center/internal/services"
)

// ProjectHandler handles project endpoints.
type ProjectHandler struct {
	projectRepo *repository.ProjectRepository
	apiKeyRepo  *repository.APIKeyRepository
	userService *services.UserSyncService
	logger      *slog.Logger
}

// NewProjectHandler creates a new project handler.
func NewProjectHandler(
	projectRepo *repository.ProjectRepository,
	apiKeyRepo *repository.APIKeyRepository,
	userService *services.UserSyncService,
	logger *slog.Logger,
) *ProjectHandler {
	return &ProjectHandler{
		projectRepo: projectRepo,
		apiKeyRepo:  apiKeyRepo,
		userService: userService,
		logger:      logger,
	}
}

// List returns all projects.
// GET /projects
func (h *ProjectHandler) List(c *gin.Context) {
	limit, offset := GetPaginationParams(c)

	projects, total, err := h.projectRepo.List(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("failed to list projects", "error", err)
		InternalError(c, "failed to list projects")
		return
	}

	Paginated(c, projects, total, limit, offset)
}

// Create creates a new project.
// POST /projects
func (h *ProjectHandler) Create(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		Unauthorized(c, "not authenticated")
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Slug        string `json:"slug"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	// Get user from database
	dbUser, err := h.userService.GetUserByKeycloakID(c.Request.Context(), user.KeycloakID)
	if err != nil {
		NotFound(c, "user not found")
		return
	}

	// Generate slug if not provided
	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Name)
	}

	project := &models.Project{
		Name:        req.Name,
		Description: req.Description,
		Slug:        slug,
	}

	if err := h.projectRepo.Create(c.Request.Context(), project); err != nil {
		h.logger.Error("failed to create project", "error", err)
		InternalError(c, "failed to create project")
		return
	}

	// Get owner role
	ownerRole, err := h.projectRepo.GetRoleByName(c.Request.Context(), models.RoleOwner)
	if err != nil {
		h.logger.Error("failed to get owner role", "error", err)
		InternalError(c, "failed to get owner role")
		return
	}

	// Add creator as owner
	member := &models.ProjectMember{
		ProjectID: project.ID,
		UserID:    dbUser.ID,
		RoleID:    ownerRole.ID,
	}

	if err := h.projectRepo.AddMember(c.Request.Context(), member); err != nil {
		h.logger.Error("failed to add member", "error", err)
		InternalError(c, "failed to add member")
		return
	}

	// Generate API key for owner
	apiKey := &models.APIKey{
		ProjectID: project.ID,
		UserID:    dbUser.ID,
		Name:      "Default Key",
	}

	fullKey, err := h.apiKeyRepo.Create(c.Request.Context(), apiKey)
	if err != nil {
		h.logger.Error("failed to create API key", "error", err)
	}

	Created(c, gin.H{
		"project": project,
		"api_key": fullKey,
	})
}

// Get returns a single project.
// GET /projects/:id
func (h *ProjectHandler) Get(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		BadRequest(c, "invalid project ID")
		return
	}

	project, err := h.projectRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		NotFound(c, "project not found")
		return
	}

	Success(c, project)
}

// Update updates a project.
// PATCH /projects/:id
func (h *ProjectHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		BadRequest(c, "invalid project ID")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	project, err := h.projectRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		NotFound(c, "project not found")
		return
	}

	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Description != "" {
		project.Description = req.Description
	}

	if err := h.projectRepo.Update(c.Request.Context(), project); err != nil {
		h.logger.Error("failed to update project", "error", err)
		InternalError(c, "failed to update project")
		return
	}

	Success(c, project)
}

// Delete deletes a project.
// DELETE /projects/:id
func (h *ProjectHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		BadRequest(c, "invalid project ID")
		return
	}

	if err := h.projectRepo.Delete(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete project", "error", err)
		InternalError(c, "failed to delete project")
		return
	}

	NoContent(c)
}

// ListMembers returns all members of a project.
// GET /projects/:id/members
func (h *ProjectHandler) ListMembers(c *gin.Context) {
	idParam := c.Param("id")
	projectID, err := uuid.Parse(idParam)
	if err != nil {
		BadRequest(c, "invalid project ID")
		return
	}

	limit, offset := GetPaginationParams(c)

	members, total, err := h.projectRepo.ListMembers(c.Request.Context(), projectID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list members", "error", err)
		InternalError(c, "failed to list members")
		return
	}

	Paginated(c, members, total, limit, offset)
}

// AddMember adds a member to a project.
// POST /projects/:id/members
func (h *ProjectHandler) AddMember(c *gin.Context) {
	idParam := c.Param("id")
	projectID, err := uuid.Parse(idParam)
	if err != nil {
		BadRequest(c, "invalid project ID")
		return
	}

	inviter, err := middleware.GetUserFromContext(c)
	if err != nil {
		Unauthorized(c, "not authenticated")
		return
	}

	inviterDB, err := h.userService.GetUserByKeycloakID(c.Request.Context(), inviter.KeycloakID)
	if err != nil {
		NotFound(c, "inviter not found")
		return
	}

	var req struct {
		UserID uuid.UUID         `json:"user_id" binding:"required"`
		Role   models.MemberRole `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	// Get role
	role, err := h.projectRepo.GetRoleByName(c.Request.Context(), req.Role)
	if err != nil {
		BadRequest(c, "invalid role")
		return
	}

	member := &models.ProjectMember{
		ProjectID: projectID,
		UserID:    req.UserID,
		RoleID:    role.ID,
		InvitedBy: &inviterDB.ID,
	}

	if err := h.projectRepo.AddMember(c.Request.Context(), member); err != nil {
		h.logger.Error("failed to add member", "error", err)
		InternalError(c, "failed to add member")
		return
	}

	// Generate API key for new member
	apiKey := &models.APIKey{
		ProjectID: projectID,
		UserID:    req.UserID,
		Name:      "Default Key",
	}

	fullKey, err := h.apiKeyRepo.Create(c.Request.Context(), apiKey)
	if err != nil {
		h.logger.Error("failed to create API key", "error", err)
	}

	Created(c, gin.H{
		"member":  member,
		"api_key": fullKey,
	})
}

// UpdateMember updates a member's role.
// PATCH /projects/:id/members/:memberId
func (h *ProjectHandler) UpdateMember(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequest(c, "invalid project ID")
		return
	}

	memberID, err := uuid.Parse(c.Param("memberId"))
	if err != nil {
		BadRequest(c, "invalid member ID")
		return
	}

	var req struct {
		Role models.MemberRole `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	role, err := h.projectRepo.GetRoleByName(c.Request.Context(), req.Role)
	if err != nil {
		BadRequest(c, "invalid role")
		return
	}

	if err := h.projectRepo.UpdateMemberRole(c.Request.Context(), projectID, memberID, role.ID); err != nil {
		h.logger.Error("failed to update member", "error", err)
		InternalError(c, "failed to update member")
		return
	}

	Success(c, gin.H{"message": "member updated"})
}

// RemoveMember removes a member from a project.
// DELETE /projects/:id/members/:memberId
func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequest(c, "invalid project ID")
		return
	}

	memberID, err := uuid.Parse(c.Param("memberId"))
	if err != nil {
		BadRequest(c, "invalid member ID")
		return
	}

	if err := h.projectRepo.RemoveMember(c.Request.Context(), projectID, memberID); err != nil {
		h.logger.Error("failed to remove member", "error", err)
		InternalError(c, "failed to remove member")
		return
	}

	NoContent(c)
}

// generateSlug creates a URL-friendly slug from a name.
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

package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/your-org/notification-center/internal/data/repository"
	"github.com/your-org/notification-center/internal/data/services"
	"github.com/your-org/notification-center/pkg/middleware"
	"github.com/your-org/notification-center/pkg/response"
)

// UserHandler handles user endpoints.
type UserHandler struct {
	userService services.UserSyncService
	projectRepo repository.ProjectRepository
	logger      *slog.Logger
}

// NewUserHandler creates a new user handler.
func NewUserHandler(userService services.UserSyncService, projectRepo repository.ProjectRepository, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		projectRepo: projectRepo,
		logger:      logger,
	}
}

// GetMe returns the current user profile.
// GET /users/me
func (h *UserHandler) GetMe(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		response.Unauthorized(c, "not authenticated")
		return
	}

	// Get full user from database
	res, err := h.userService.GetUserByKeycloakID(c.Request.Context(), user.KeycloakID)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	response.Success(c, res)
}

// UpdateMe updates the current user profile.
// PATCH /users/me
func (h *UserHandler) UpdateMe(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		response.Unauthorized(c, "not authenticated")
		return
	}

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get user from database
	res, err := h.userService.GetUserByKeycloakID(c.Request.Context(), user.KeycloakID)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	// Update fields
	if req.FirstName != "" {
		res.FirstName = req.FirstName
	}
	if req.LastName != "" {
		res.LastName = req.LastName
	}

	// Note: In a full implementation, you would save the updated user to DB

	response.Success(c, res)
}

// GetMyProjects returns projects the current user is a member of.
// GET /users/me/projects
func (h *UserHandler) GetMyProjects(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		response.Unauthorized(c, "not authenticated")
		return
	}

	// Get user from database
	res, err := h.userService.GetUserByKeycloakID(c.Request.Context(), user.KeycloakID)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	limit, offset := response.GetPaginationParams(c)

	projects, total, err := h.projectRepo.ListByOwner(c.Request.Context(), res.ID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list user projects", "error", err)
		response.InternalError(c, "failed to list projects")
		return
	}

	response.Paginated(c, projects, total, limit, offset)
}

package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/your-org/notification-center/internal/middleware"
	"github.com/your-org/notification-center/internal/repository"
	"github.com/your-org/notification-center/internal/services"
)

// UserHandler handles user endpoints.
type UserHandler struct {
	userService *services.UserSyncService
	projectRepo *repository.ProjectRepository
	logger      *slog.Logger
}

// NewUserHandler creates a new user handler.
func NewUserHandler(userService *services.UserSyncService, projectRepo *repository.ProjectRepository, logger *slog.Logger) *UserHandler {
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
		Unauthorized(c, "not authenticated")
		return
	}

	// Get full user from database
	dbUser, err := h.userService.GetUserByKeycloakID(c.Request.Context(), user.KeycloakID)
	if err != nil {
		NotFound(c, "user not found")
		return
	}

	Success(c, dbUser)
}

// UpdateMe updates the current user profile.
// PATCH /users/me
func (h *UserHandler) UpdateMe(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		Unauthorized(c, "not authenticated")
		return
	}

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		AvatarURL string `json:"avatar_url"`
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

	// Update fields
	if req.FirstName != "" {
		dbUser.FirstName = req.FirstName
	}
	if req.LastName != "" {
		dbUser.LastName = req.LastName
	}
	if req.AvatarURL != "" {
		dbUser.AvatarURL = req.AvatarURL
	}

	// Note: In a full implementation, you would save the updated user to DB

	Success(c, dbUser)
}

// GetMyProjects returns projects the current user is a member of.
// GET /users/me/projects
func (h *UserHandler) GetMyProjects(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		Unauthorized(c, "not authenticated")
		return
	}

	// Get user from database
	dbUser, err := h.userService.GetUserByKeycloakID(c.Request.Context(), user.KeycloakID)
	if err != nil {
		NotFound(c, "user not found")
		return
	}

	limit, offset := GetPaginationParams(c)

	projects, total, err := h.projectRepo.ListByUser(c.Request.Context(), dbUser.ID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list user projects", "error", err)
		InternalError(c, "failed to list projects")
		return
	}

	Paginated(c, projects, total, limit, offset)
}

package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/your-org/notification-center/internal/data/services"
	"github.com/your-org/notification-center/pkg/middleware"
	"github.com/your-org/notification-center/pkg/response"
)

type ProjectHandler struct {
	projectService services.ProjectService
	logger         *slog.Logger
}

func NewProjectHandler(projectService services.ProjectService, logger *slog.Logger) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
		logger:         logger,
	}
}

// GET /projects
func (h *ProjectHandler) List(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		response.Unauthorized(c, "not authenticated")
		return
	}

	limit, offset := response.GetPaginationParams(c)

	projects, total, err := h.projectService.List(c.Request.Context(), user.KeycloakID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list projects", "error", err)
		response.InternalError(c, "failed to list projects")
		return
	}

	response.Paginated(c, projects, total, limit, offset)
}

// POST /projects
func (h *ProjectHandler) Create(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		response.Unauthorized(c, "not authenticated")
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Slug        string `json:"slug"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	project, apiKey, err := h.projectService.Create(c.Request.Context(), user.KeycloakID, req.Name, req.Description, req.Slug)
	if err != nil {
		h.logger.Error("failed to create project", "error", err)
		response.InternalError(c, "failed to create project")
		return
	}

	response.Created(c, gin.H{
		"project": project,
		"api_key": apiKey,
	})
}

// GET /projects/:id
func (h *ProjectHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid project ID")
		return
	}

	project, err := h.projectService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "project not found")
		return
	}

	response.Success(c, project)
}

// PATCH /projects/:id
func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid project ID")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	project, err := h.projectService.Update(c.Request.Context(), id, req.Name, req.Description)
	if err != nil {
		h.logger.Error("failed to update project", "error", err)
		response.InternalError(c, "failed to update project")
		return
	}

	response.Success(c, project)
}

// DELETE /projects/:id
func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid project ID")
		return
	}

	if err := h.projectService.Delete(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete project", "error", err)
		response.InternalError(c, "failed to delete project")
		return
	}

	response.NoContent(c)
}

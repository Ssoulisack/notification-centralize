package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/your-org/notification-center/internal/middleware"
	"github.com/your-org/notification-center/internal/models"
	"github.com/your-org/notification-center/internal/repository"
	"github.com/your-org/notification-center/internal/services"
)

// APIKeyHandler handles API key endpoints.
type APIKeyHandler struct {
	apiKeyRepo  *repository.APIKeyRepository
	userService *services.UserSyncService
	logger      *slog.Logger
}

// NewAPIKeyHandler creates a new API key handler.
func NewAPIKeyHandler(
	apiKeyRepo *repository.APIKeyRepository,
	userService *services.UserSyncService,
	logger *slog.Logger,
) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyRepo:  apiKeyRepo,
		userService: userService,
		logger:      logger,
	}
}

// List returns all API keys for a project.
// GET /projects/:id/api-keys
func (h *APIKeyHandler) List(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequest(c, "invalid project ID")
		return
	}

	limit, offset := GetPaginationParams(c)

	keys, total, err := h.apiKeyRepo.ListByProject(c.Request.Context(), projectID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list API keys", "error", err)
		InternalError(c, "failed to list API keys")
		return
	}

	Paginated(c, keys, total, limit, offset)
}

// Create creates a new API key.
// POST /projects/:id/api-keys
func (h *APIKeyHandler) Create(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequest(c, "invalid project ID")
		return
	}

	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		Unauthorized(c, "not authenticated")
		return
	}

	dbUser, err := h.userService.GetUserByKeycloakID(c.Request.Context(), user.KeycloakID)
	if err != nil {
		NotFound(c, "user not found")
		return
	}

	var req struct {
		Name      string `json:"name" binding:"required"`
		ExpiresAt string `json:"expires_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	apiKey := &models.APIKey{
		ProjectID: projectID,
		UserID:    dbUser.ID,
		Name:      req.Name,
	}

	fullKey, err := h.apiKeyRepo.Create(c.Request.Context(), apiKey)
	if err != nil {
		h.logger.Error("failed to create API key", "error", err)
		InternalError(c, "failed to create API key")
		return
	}

	Created(c, gin.H{
		"id":         apiKey.ID,
		"key":        fullKey,
		"key_prefix": apiKey.KeyPrefix,
		"name":       apiKey.Name,
		"created_at": apiKey.CreatedAt,
		"message":    "Save this key securely. You won't be able to see it again.",
	})
}

// Revoke revokes an API key.
// DELETE /projects/:id/api-keys/:keyId
func (h *APIKeyHandler) Revoke(c *gin.Context) {
	keyID, err := uuid.Parse(c.Param("keyId"))
	if err != nil {
		BadRequest(c, "invalid key ID")
		return
	}

	if err := h.apiKeyRepo.Revoke(c.Request.Context(), keyID); err != nil {
		h.logger.Error("failed to revoke API key", "error", err)
		InternalError(c, "failed to revoke API key")
		return
	}

	NoContent(c)
}

package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/your-org/notification-center/internal/data/services"
	"github.com/your-org/notification-center/pkg/middleware"
	"github.com/your-org/notification-center/pkg/response"
)

type APIKeyHandler interface {
	List(c *gin.Context)
	Create(c *gin.Context)
	Revoke(c *gin.Context)
}

type apiKeyHandler struct {
	apiKeyService services.APIKeyService
	logger        *slog.Logger
}

func NewAPIKeyHandler(apiKeyService services.APIKeyService, logger *slog.Logger) APIKeyHandler {
	return &apiKeyHandler{
		apiKeyService: apiKeyService,
		logger:        logger,
	}
}

// GET /projects/:id/api-keys
func (h *apiKeyHandler) List(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid project ID")
		return
	}

	limit, offset := response.GetPaginationParams(c)

	keys, total, err := h.apiKeyService.List(c.Request.Context(), projectID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list API keys", "error", err)
		response.InternalError(c, "failed to list API keys")
		return
	}

	response.Paginated(c, keys, int64(total), limit, offset)
}

// POST /projects/:id/api-keys
func (h *apiKeyHandler) Create(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid project ID")
		return
	}

	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		response.Unauthorized(c, "not authenticated")
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	apiKey, fullKey, err := h.apiKeyService.Create(c.Request.Context(), projectID, user.KeycloakID, req.Name)
	if err != nil {
		h.logger.Error("failed to create API key", "error", err)
		response.InternalError(c, "failed to create API key")
		return
	}

	response.Created(c, gin.H{
		"id":         apiKey.ID,
		"key":        fullKey,
		"key_prefix": apiKey.KeyPrefix,
		"name":       apiKey.Name,
		"created_at": apiKey.CreatedAt,
		"message":    "Save this key securely. You won't be able to see it again.",
	})
}

// DELETE /projects/:id/api-keys/:keyId
func (h *apiKeyHandler) Revoke(c *gin.Context) {
	keyID, err := uuid.Parse(c.Param("keyId"))
	if err != nil {
		response.BadRequest(c, "invalid key ID")
		return
	}

	if err := h.apiKeyService.Revoke(c.Request.Context(), keyID); err != nil {
		h.logger.Error("failed to revoke API key", "error", err)
		response.InternalError(c, "failed to revoke API key")
		return
	}

	response.NoContent(c)
}

package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/middleware"
	"github.com/your-org/notification-center/internal/services"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	config      *config.KeycloakConfig
	userService *services.UserSyncService
	logger      *slog.Logger
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(cfg *config.KeycloakConfig, userService *services.UserSyncService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		config:      cfg,
		userService: userService,
		logger:      logger,
	}
}

// LoginRedirect redirects the user to Keycloak for authentication.
// POST /auth/login
func (h *AuthHandler) LoginRedirect(c *gin.Context) {
	// Build Keycloak authorization URL
	authURL := h.config.BaseURL + "/realms/" + h.config.Realm + "/protocol/openid-connect/auth"

	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		redirectURI = c.Request.Host + "/auth/callback"
	}

	params := "?client_id=" + h.config.ClientID +
		"&response_type=code" +
		"&scope=openid profile email" +
		"&redirect_uri=" + redirectURI

	c.Redirect(http.StatusFound, authURL+params)
}

// Callback handles the Keycloak OAuth callback.
// POST /auth/callback
func (h *AuthHandler) Callback(c *gin.Context) {
	// In a real implementation, this would:
	// 1. Exchange the authorization code for tokens
	// 2. Validate the tokens
	// 3. Create/sync user in database
	// 4. Return access token to client

	code := c.Query("code")
	if code == "" {
		BadRequest(c, "missing authorization code")
		return
	}

	// Token exchange would happen here
	// For now, return a placeholder response
	Success(c, gin.H{
		"message": "callback received",
		"code":    code,
	})
}

// Logout handles user logout.
// POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a real implementation, this would:
	// 1. Invalidate the user's session
	// 2. Optionally redirect to Keycloak logout endpoint

	logoutURL := h.config.BaseURL + "/realms/" + h.config.Realm + "/protocol/openid-connect/logout"

	Success(c, gin.H{
		"message":    "logged out",
		"logout_url": logoutURL,
	})
}

// Me returns the current authenticated user.
// GET /auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		Unauthorized(c, "not authenticated")
		return
	}

	// Sync user to database
	user, err := h.userService.SyncUser(c.Request.Context(), claims)
	if err != nil {
		h.logger.Error("failed to sync user", "error", err)
		InternalError(c, "failed to sync user")
		return
	}

	Success(c, user)
}

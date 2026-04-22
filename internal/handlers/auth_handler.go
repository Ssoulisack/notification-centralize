package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/data/services"
	"github.com/your-org/notification-center/pkg/middleware"
	"github.com/your-org/notification-center/pkg/response"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	config      *config.KeycloakConfig
	userService services.UserSyncService
	logger      *slog.Logger
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(cfg *config.KeycloakConfig, userService services.UserSyncService, logger *slog.Logger) *AuthHandler {
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
		response.BadRequest(c, "missing authorization code")
		return
	}

	// Token exchange would happen here
	// For now, return a placeholder response
	response.Success(c, gin.H{
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

	response.Success(c, gin.H{
		"message":    "logged out",
		"logout_url": logoutURL,
	})
}

// Token exchanges username/password directly for a Keycloak access token.
// POST /auth/token
func (h *AuthHandler) Token(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", h.config.BaseURL, h.config.Realm)

	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("client_id", h.config.ClientID)
	form.Set("client_secret", h.config.ClientSecret)
	form.Set("username", req.Username)
	form.Set("password", req.Password)
	form.Set("scope", "openid profile email")

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		h.logger.Error("failed to call keycloak token endpoint", "error", err)
		response.InternalError(c, "failed to authenticate")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var keycloakErr struct {
			ErrorDescription string `json:"error_description"`
		}
		_ = json.Unmarshal(body, &keycloakErr)
		response.Error(c, http.StatusUnauthorized, keycloakErr.ErrorDescription)
		return
	}

	var token map[string]any
	if err := json.Unmarshal(body, &token); err != nil {
		response.InternalError(c, "failed to parse token response")
		return
	}

	response.Success(c, token)
}

// Me returns the current authenticated user.
// GET /auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		response.Unauthorized(c, "not authenticated")
		return
	}

	// Sync user to database
	user, err := h.userService.Sync(c.Request.Context(), claims)
	if err != nil {
		h.logger.Error("failed to sync user", "error", err)
		response.InternalError(c, "failed to sync user")
		return
	}

	response.Success(c, user)
}

package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/models"
)

// Context keys for user information.
type contextKey string

const (
	UserContextKey   contextKey = "user"
	ClaimsContextKey contextKey = "claims"
)

// KeycloakClaims represents JWT claims from Keycloak.
type KeycloakClaims struct {
	jwt.RegisteredClaims
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	PreferredUsername string `json:"preferred_username"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Name              string `json:"name"`
}

// AuthMiddleware validates Keycloak JWT tokens.
type AuthMiddleware struct {
	jwks   *keyfunc.JWKS
	config *config.KeycloakConfig
	logger *slog.Logger
}

// NewAuthMiddleware creates a new Keycloak auth middleware.
func NewAuthMiddleware(cfg *config.KeycloakConfig, logger *slog.Logger) (*AuthMiddleware, error) {
	options := keyfunc.Options{
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
	}

	jwks, err := keyfunc.Get(cfg.JWKSURL(), options)
	if err != nil {
		return nil, err
	}

	return &AuthMiddleware{
		jwks:   jwks,
		config: cfg,
		logger: logger,
	}, nil
}

// Handler returns the Gin middleware handler.
func (m *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			return
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(tokenString, &KeycloakClaims{}, m.jwks.Keyfunc)
		if err != nil {
			m.logger.Warn("failed to parse token", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		claims, ok := token.Claims.(*KeycloakClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			return
		}

		// Create user from claims
		user := &models.User{
			KeycloakID:    claims.Subject,
			Email:         claims.Email,
			Username:      claims.PreferredUsername,
			FirstName:     claims.GivenName,
			LastName:      claims.FamilyName,
			EmailVerified: claims.EmailVerified,
		}

		// Store in context
		c.Set(string(UserContextKey), user)
		c.Set(string(ClaimsContextKey), claims)

		c.Next()
	}
}

// Close cleans up JWKS resources.
func (m *AuthMiddleware) Close() {
	if m.jwks != nil {
		m.jwks.EndBackground()
	}
}

// GetUserFromContext retrieves the user from the context.
func GetUserFromContext(c *gin.Context) (*models.User, error) {
	user, exists := c.Get(string(UserContextKey))
	if !exists {
		return nil, errors.New("user not found in context")
	}

	u, ok := user.(*models.User)
	if !ok {
		return nil, errors.New("invalid user type in context")
	}

	return u, nil
}

// GetClaimsFromContext retrieves the JWT claims from the context.
func GetClaimsFromContext(c *gin.Context) (*KeycloakClaims, error) {
	claims, exists := c.Get(string(ClaimsContextKey))
	if !exists {
		return nil, errors.New("claims not found in context")
	}

	cl, ok := claims.(*KeycloakClaims)
	if !ok {
		return nil, errors.New("invalid claims type in context")
	}

	return cl, nil
}

// GetUserIDFromContext retrieves the user ID from the context.
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	if c, ok := ctx.(*gin.Context); ok {
		user, err := GetUserFromContext(c)
		if err != nil {
			return uuid.Nil, err
		}
		return user.ID, nil
	}
	return uuid.Nil, errors.New("invalid context type")
}

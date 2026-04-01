package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-org/notification-center/internal/models"
)

const (
	// APIKeyHeader is the header name for API key authentication.
	APIKeyHeader = "X-API-Key"

	// APIKeyPrefix is the prefix for API keys.
	APIKeyPrefix = "nc_live_"

	// APIKeyContextKey is the context key for storing API key info.
	APIKeyContextKey contextKey = "api_key"

	// ProjectContextKey is the context key for storing project info.
	ProjectContextKey contextKey = "project"
)

// APIKeyMiddleware validates API keys and sets project context.
type APIKeyMiddleware struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

// NewAPIKeyMiddleware creates a new API key middleware.
func NewAPIKeyMiddleware(db *pgxpool.Pool, logger *slog.Logger) *APIKeyMiddleware {
	return &APIKeyMiddleware{
		db:     db,
		logger: logger,
	}
}

// Handler returns the Gin middleware handler.
func (m *APIKeyMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(APIKeyHeader)
		if apiKey == "" {
			// Fallback to query parameter
			apiKey = c.Query("api_key")
		}

		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing API key",
			})
			return
		}

		// Validate prefix
		if !strings.HasPrefix(apiKey, APIKeyPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid API key format",
			})
			return
		}

		// Hash the key for lookup
		keyHash := hashAPIKey(apiKey)

		// Lookup in database
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		var key models.APIKey
		var project models.Project

		query := `
			SELECT
				ak.id, ak.project_id, ak.user_id, ak.name, ak.key_prefix,
				ak.scopes, ak.last_used_at, ak.expires_at, ak.is_active, ak.created_at,
				p.id, p.name, p.description, p.slug, p.is_active
			FROM api_keys ak
			JOIN projects p ON p.id = ak.project_id
			WHERE ak.key_hash = $1 AND ak.is_active = true
		`

		err := m.db.QueryRow(ctx, query, keyHash).Scan(
			&key.ID, &key.ProjectID, &key.UserID, &key.Name, &key.KeyPrefix,
			&key.Scopes, &key.LastUsedAt, &key.ExpiresAt, &key.IsActive, &key.CreatedAt,
			&project.ID, &project.Name, &project.Description, &project.Slug, &project.IsActive,
		)

		if err != nil {
			m.logger.Warn("API key lookup failed", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid API key",
			})
			return
		}

		// Check expiration
		if key.ExpiresAt != nil && key.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "API key expired",
			})
			return
		}

		// Check project is active
		if !project.IsActive {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "project is inactive",
			})
			return
		}

		// Update last used timestamp (async)
		go m.updateLastUsed(context.Background(), key.ID)

		// Store in context
		c.Set(string(APIKeyContextKey), &key)
		c.Set(string(ProjectContextKey), &project)

		c.Next()
	}
}

// updateLastUsed updates the last_used_at timestamp.
func (m *APIKeyMiddleware) updateLastUsed(ctx context.Context, keyID uuid.UUID) {
	query := `UPDATE api_keys SET last_used_at = NOW() WHERE id = $1`
	_, err := m.db.Exec(ctx, query, keyID)
	if err != nil {
		m.logger.Error("failed to update API key last_used_at", "error", err)
	}
}

// hashAPIKey returns the SHA-256 hash of an API key.
func hashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

// GetAPIKeyFromContext retrieves the API key from the context.
func GetAPIKeyFromContext(c *gin.Context) (*models.APIKey, error) {
	key, exists := c.Get(string(APIKeyContextKey))
	if !exists {
		return nil, ErrAPIKeyNotFound
	}

	k, ok := key.(*models.APIKey)
	if !ok {
		return nil, ErrInvalidAPIKey
	}

	return k, nil
}

// GetProjectFromContext retrieves the project from the context.
func GetProjectFromContext(c *gin.Context) (*models.Project, error) {
	project, exists := c.Get(string(ProjectContextKey))
	if !exists {
		return nil, ErrProjectNotFound
	}

	p, ok := project.(*models.Project)
	if !ok {
		return nil, ErrInvalidProject
	}

	return p, nil
}

// GetProjectIDFromContext retrieves the project ID from the context.
func GetProjectIDFromContext(c *gin.Context) (uuid.UUID, error) {
	project, err := GetProjectFromContext(c)
	if err != nil {
		return uuid.Nil, err
	}
	return project.ID, nil
}

// Errors
var (
	ErrAPIKeyNotFound  = &ContextError{Message: "API key not found in context"}
	ErrInvalidAPIKey   = &ContextError{Message: "invalid API key type in context"}
	ErrProjectNotFound = &ContextError{Message: "project not found in context"}
	ErrInvalidProject  = &ContextError{Message: "invalid project type in context"}
)

// ContextError represents a context retrieval error.
type ContextError struct {
	Message string
}

func (e *ContextError) Error() string {
	return e.Message
}

// GenerateAPIKey creates a new API key.
func GenerateAPIKey() (key string, prefix string, hash string, err error) {
	// Generate 32 random bytes
	id := uuid.New()
	key = APIKeyPrefix + strings.ReplaceAll(id.String(), "-", "")
	prefix = key[:len(APIKeyPrefix)+8]
	hash = hashAPIKey(key)
	return key, prefix, hash, nil
}

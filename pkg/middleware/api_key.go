package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/notification-center/internal/domain/models"
)

const (
	APIKeyHeader      = "X-API-Key"
	APIKeyPrefix      = "nc_live_"
	APIKeyContextKey  contextKey = "api_key"
	ProjectContextKey contextKey = "project"
)

type APIKeyMiddleware struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewAPIKeyMiddleware(db *gorm.DB, logger *slog.Logger) *APIKeyMiddleware {
	return &APIKeyMiddleware{db: db, logger: logger}
}

func (m *APIKeyMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(APIKeyHeader)
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing API key"})
			return
		}
		if !strings.HasPrefix(apiKey, APIKeyPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid API key format"})
			return
		}

		keyHash := hashAPIKey(apiKey)

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		var key models.APIKey
		err := m.db.WithContext(ctx).Preload("Project").Where("key_hash = ? AND is_active = true", keyHash).First(&key).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			} else {
				m.logger.Warn("API key lookup failed", "error", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			}
			return
		}

		if key.ExpiresAt != nil && key.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key expired"})
			return
		}
		if !key.Project.IsActive {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "project is inactive"})
			return
		}

		go m.updateLastUsed(context.Background(), key.ID)

		c.Set(string(APIKeyContextKey), &key)
		c.Set(string(ProjectContextKey), &key.Project)

		c.Next()
	}
}

func (m *APIKeyMiddleware) updateLastUsed(ctx context.Context, keyID uuid.UUID) {
	if err := m.db.WithContext(ctx).Model(&models.APIKey{}).Where("id = ?", keyID).Update("last_used_at", time.Now()).Error; err != nil {
		m.logger.Error("failed to update API key last_used_at", "error", err)
	}
}

func hashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

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

func GetProjectIDFromContext(c *gin.Context) (uuid.UUID, error) {
	project, err := GetProjectFromContext(c)
	if err != nil {
		return uuid.Nil, err
	}
	return project.ID, nil
}

var (
	ErrAPIKeyNotFound  = &ContextError{Message: "API key not found in context"}
	ErrInvalidAPIKey   = &ContextError{Message: "invalid API key type in context"}
	ErrProjectNotFound = &ContextError{Message: "project not found in context"}
	ErrInvalidProject  = &ContextError{Message: "invalid project type in context"}
)

type ContextError struct {
	Message string
}

func (e *ContextError) Error() string { return e.Message }

func GenerateAPIKey() (key string, prefix string, hash string, err error) {
	id := uuid.New()
	key = APIKeyPrefix + strings.ReplaceAll(id.String(), "-", "")
	prefix = key[:len(APIKeyPrefix)+8]
	hash = hashAPIKey(key)
	return key, prefix, hash, nil
}

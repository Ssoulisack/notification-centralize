package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	apicontext "github.com/your-org/notification-center/internal/api/context"
	"github.com/your-org/notification-center/internal/store"
)

// ProjectAuth validates the X-API-Key header and injects the project into context.
func ProjectAuth(projectStore store.ProjectStore, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing X-API-Key header",
			})
			return
		}

		project, err := projectStore.GetByAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			logger.Warn("invalid API key", "error", err)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "invalid API key",
			})
			return
		}

		ctx := apicontext.WithProject(c.Request.Context(), project)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// RequestLogger logs each request with structured fields.
func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		logger.Info("request",
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}

// CORS adds permissive CORS headers for development.
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-API-Key")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

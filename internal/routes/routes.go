package routes

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/your-org/notification-center/bootstrap/messaging"
	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/pkg/middleware"
)

// Dependencies holds all dependencies needed for route setup.
type Dependencies struct {
	GormDB           *gorm.DB
	RabbitMQ         *messaging.Client
	Config           *config.Config
	Logger           *slog.Logger
	AuthMiddleware   *middleware.AuthMiddleware
	APIKeyMiddleware *middleware.APIKeyMiddleware
}

// Setup configures all application routes.
func Setup(router *gin.Engine, deps *Dependencies) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes (public)
	SetupAuthRoutes(router, deps)

	// Protected routes (JWT authenticated)
	protected := router.Group("/api/v1")
	protected.Use(deps.AuthMiddleware.Handler())
	{
		SetupUserRoutes(protected, deps)
		SetupProjectRoutes(protected, deps)
		SetupInboxRoutes(protected, deps)
	}

	// External API routes (API key authenticated)
	external := router.Group("/api/v1")
	external.Use(deps.APIKeyMiddleware.Handler())
	{
		SetupExternalRoutes(external, deps)
	}
}

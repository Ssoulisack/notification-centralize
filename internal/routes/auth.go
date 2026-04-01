package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/your-org/notification-center/internal/handlers"
	"github.com/your-org/notification-center/internal/services"
)

// SetupAuthRoutes configures authentication routes.
func SetupAuthRoutes(router *gin.Engine, deps *Dependencies) {
	// Initialize dependencies
	userService := services.NewUserSyncService(deps.DB, deps.Logger)
	handler := handlers.NewAuthHandler(&deps.Config.Keycloak, userService, deps.Logger)

	// Auth route group
	auth := router.Group("/auth")
	{
		// Public routes
		auth.POST("/login", handler.LoginRedirect)
		auth.POST("/callback", handler.Callback)
		auth.POST("/logout", handler.Logout)

		// Protected route
		auth.GET("/me", deps.AuthMiddleware.Handler(), handler.Me)
	}
}

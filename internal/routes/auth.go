package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/your-org/notification-center/internal/data/repository"
	"github.com/your-org/notification-center/internal/data/services"
	"github.com/your-org/notification-center/internal/handlers"
)

// SetupAuthRoutes configures authentication routes.
func SetupAuthRoutes(router *gin.Engine, deps *Dependencies) {
	userRepo := repository.NewUserRepository(deps.GormDB)
	userService := services.NewUserSyncService(userRepo, deps.Logger)
	handler := handlers.NewAuthHandler(&deps.Config.Keycloak, userService, deps.Logger)

	// Auth route group
	auth := router.Group("/auth")
	{
		// Public routes
		auth.POST("/login", handler.LoginRedirect)
		auth.POST("/token", handler.Token)
		auth.POST("/callback", handler.Callback)
		auth.POST("/logout", handler.Logout)

		// Protected route
		auth.GET("/me", deps.AuthMiddleware.Handler(), handler.Me)
	}
}

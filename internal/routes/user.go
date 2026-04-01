package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/your-org/notification-center/internal/handlers"
	"github.com/your-org/notification-center/internal/repository"
	"github.com/your-org/notification-center/internal/services"
)

// SetupUserRoutes configures user routes.
func SetupUserRoutes(router *gin.RouterGroup, deps *Dependencies) {
	// Initialize dependencies
	userService := services.NewUserSyncService(deps.DB, deps.Logger)
	projectRepo := repository.NewProjectRepository(deps.DB)
	handler := handlers.NewUserHandler(userService, projectRepo, deps.Logger)

	// User route group
	users := router.Group("/users")
	{
		users.GET("/me", handler.GetMe)
		users.PATCH("/me", handler.UpdateMe)
		users.GET("/me/projects", handler.GetMyProjects)
	}
}

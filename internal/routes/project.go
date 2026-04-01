package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/your-org/notification-center/internal/handlers"
	"github.com/your-org/notification-center/internal/repository"
	"github.com/your-org/notification-center/internal/services"
)

// SetupProjectRoutes configures project routes.
func SetupProjectRoutes(router *gin.RouterGroup, deps *Dependencies) {
	// Initialize dependencies
	projectRepo := repository.NewProjectRepository(deps.DB)
	apiKeyRepo := repository.NewAPIKeyRepository(deps.DB)
	userService := services.NewUserSyncService(deps.DB, deps.Logger)
	notifService := services.NewNotificationService(deps.DB, deps.RabbitMQ, deps.Logger)

	projectHandler := handlers.NewProjectHandler(projectRepo, apiKeyRepo, userService, deps.Logger)
	notifHandler := handlers.NewNotificationHandler(notifService, userService, deps.Logger)
	apiKeyHandler := handlers.NewAPIKeyHandler(apiKeyRepo, userService, deps.Logger)

	// Project route group
	projects := router.Group("/projects")
	{
		// Project CRUD
		projects.GET("", projectHandler.List)
		projects.POST("", projectHandler.Create)
		projects.GET("/:id", projectHandler.Get)
		projects.PATCH("/:id", projectHandler.Update)
		projects.DELETE("/:id", projectHandler.Delete)

		// Project members
		projects.GET("/:id/members", projectHandler.ListMembers)
		projects.POST("/:id/members", projectHandler.AddMember)
		projects.PATCH("/:id/members/:memberId", projectHandler.UpdateMember)
		projects.DELETE("/:id/members/:memberId", projectHandler.RemoveMember)

		// Project notifications
		projects.POST("/:id/notifications", notifHandler.Send)
		projects.POST("/:id/notifications/batch", notifHandler.SendBatch)
		projects.GET("/:id/notifications", notifHandler.List)
		projects.GET("/:id/notifications/:notificationId", notifHandler.Get)

		// Project API keys
		projects.GET("/:id/api-keys", apiKeyHandler.List)
		projects.POST("/:id/api-keys", apiKeyHandler.Create)
		projects.DELETE("/:id/api-keys/:keyId", apiKeyHandler.Revoke)
	}
}

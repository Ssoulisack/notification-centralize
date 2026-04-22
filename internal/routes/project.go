package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/your-org/notification-center/internal/data/repository"
	"github.com/your-org/notification-center/internal/data/services"
	"github.com/your-org/notification-center/internal/handlers"
)

func SetupProjectRoutes(router *gin.RouterGroup, deps *Dependencies) {
	userRepo := repository.NewUserRepository(deps.GormDB)
	projectRepo := repository.NewProjectRepository(deps.GormDB)
	notifRepo := repository.NewNotificationRepository(deps.GormDB)
	apiKeyRepo := repository.NewAPIKeyRepository(deps.GormDB)

	userService := services.NewUserSyncService(userRepo, deps.Logger)
	projectService := services.NewProjectService(projectRepo, apiKeyRepo, userService)
	apiKeyService := services.NewAPIKeyService(apiKeyRepo, userService)
	notifService := services.NewNotificationService(notifRepo, deps.RabbitMQ, deps.Logger)

	projectHandler := handlers.NewProjectHandler(projectService, deps.Logger)
	notifHandler := handlers.NewNotificationHandler(notifService, userService, deps.Logger)
	apiKeyHandler := handlers.NewAPIKeyHandler(apiKeyService, deps.Logger)

	projects := router.Group("/projects")
	{
		projects.GET("", projectHandler.List)
		projects.POST("", projectHandler.Create)
		projects.GET("/:id", projectHandler.Get)
		projects.PATCH("/:id", projectHandler.Update)
		projects.DELETE("/:id", projectHandler.Delete)

		projects.POST("/:id/notifications", notifHandler.Send)
		projects.POST("/:id/notifications/batch", notifHandler.SendBatch)
		projects.GET("/:id/notifications", notifHandler.List)
		projects.GET("/:id/notifications/:notificationId", notifHandler.Get)

		projects.GET("/:id/api-keys", apiKeyHandler.List)
		projects.POST("/:id/api-keys", apiKeyHandler.Create)
		projects.DELETE("/:id/api-keys/:keyId", apiKeyHandler.Revoke)
	}
}

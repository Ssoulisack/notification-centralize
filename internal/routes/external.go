package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/your-org/notification-center/internal/data/repository"
	"github.com/your-org/notification-center/internal/data/services"
	"github.com/your-org/notification-center/internal/handlers"
	"github.com/your-org/notification-center/pkg/middleware"
)

func SetupExternalRoutes(router *gin.RouterGroup, deps *Dependencies) {
	userRepo := repository.NewUserRepository(deps.GormDB)
	notifRepo := repository.NewNotificationRepository(deps.GormDB)
	userService := services.NewUserSyncService(userRepo, deps.Logger)
	notifService := services.NewNotificationService(notifRepo, deps.RabbitMQ, deps.Logger)
	notifHandler := handlers.NewNotificationHandler(notifService, userService, deps.Logger)

	router.POST("/notifications", func(c *gin.Context) {
		project, err := middleware.GetProjectFromContext(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.Params = append(c.Params, gin.Param{Key: "id", Value: project.ID.String()})
		notifHandler.Send(c)
	})

	router.POST("/notifications/batch", func(c *gin.Context) {
		project, err := middleware.GetProjectFromContext(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.Params = append(c.Params, gin.Param{Key: "id", Value: project.ID.String()})
		notifHandler.SendBatch(c)
	})

	router.GET("/notifications", func(c *gin.Context) {
		project, err := middleware.GetProjectFromContext(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.Params = append(c.Params, gin.Param{Key: "id", Value: project.ID.String()})
		notifHandler.List(c)
	})

	router.GET("/notifications/:notificationId", func(c *gin.Context) {
		project, err := middleware.GetProjectFromContext(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.Params = append(c.Params, gin.Param{Key: "id", Value: project.ID.String()})
		notifHandler.Get(c)
	})
}

package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/your-org/notification-center/internal/handlers"
	"github.com/your-org/notification-center/internal/middleware"
	"github.com/your-org/notification-center/internal/services"
)

// SetupExternalRoutes configures external API routes (API key authenticated).
func SetupExternalRoutes(router *gin.RouterGroup, deps *Dependencies) {
	// Initialize dependencies
	userService := services.NewUserSyncService(deps.DB, deps.Logger)
	notifService := services.NewNotificationService(deps.DB, deps.RabbitMQ, deps.Logger)
	notifHandler := handlers.NewNotificationHandler(notifService, userService, deps.Logger)

	// Notifications via API key
	router.POST("/notifications", func(c *gin.Context) {
		project, err := middleware.GetProjectFromContext(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.Set("project_id", project.ID.String())
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


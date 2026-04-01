package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/your-org/notification-center/internal/handlers"
	"github.com/your-org/notification-center/internal/services"
)

// SetupInboxRoutes configures user inbox routes.
func SetupInboxRoutes(router *gin.RouterGroup, deps *Dependencies) {
	// Initialize dependencies
	userService := services.NewUserSyncService(deps.DB, deps.Logger)
	notifService := services.NewNotificationService(deps.DB, deps.RabbitMQ, deps.Logger)
	handler := handlers.NewInboxHandler(notifService, userService, deps.Logger)

	// Inbox route group
	inbox := router.Group("/inbox")
	{
		inbox.GET("", handler.GetInbox)
		inbox.GET("/unread/count", handler.GetUnreadCount)
		inbox.POST("/:id/read", handler.MarkAsRead)
	}
}

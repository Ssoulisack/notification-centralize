package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/your-org/notification-center/internal/data/repository"
	"github.com/your-org/notification-center/internal/data/services"
	"github.com/your-org/notification-center/internal/handlers"
)

func SetupInboxRoutes(router *gin.RouterGroup, deps *Dependencies) {
	userRepo := repository.NewUserRepository(deps.GormDB)
	notifRepo := repository.NewNotificationRepository(deps.GormDB)
	userService := services.NewUserSyncService(userRepo, deps.Logger)
	notifService := services.NewNotificationService(notifRepo, deps.RabbitMQ, deps.Logger)
	handler := handlers.NewInboxHandler(notifService, userService, deps.Logger)

	inbox := router.Group("/inbox")
	{
		inbox.GET("", handler.GetInbox)
		inbox.GET("/unread/count", handler.GetUnreadCount)
		inbox.POST("/:id/read", handler.MarkAsRead)
	}
}

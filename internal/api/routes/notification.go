package routes

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/your-org/notification-center/internal/api/handler"
	"github.com/your-org/notification-center/internal/service"
)

func NewNotificationRoutes(rg *gin.RouterGroup, svc *service.NotificationService, logger *slog.Logger) {
	h := handler.NewNotificationHandler(svc, logger)

	rg.POST("/notifications", h.Send)
	rg.GET("/notifications/:id", h.Get)
	rg.GET("/notifications", h.List)
}

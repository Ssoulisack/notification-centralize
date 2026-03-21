package routes

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/your-org/notification-center/internal/api/handler"
	"github.com/your-org/notification-center/internal/store"
)

func NewDeviceRoutes(rg *gin.RouterGroup, store store.DeviceStore, logger *slog.Logger) {
	h := handler.NewDeviceHandler(store, logger)

	rg.POST("/devices", h.Register)
	rg.DELETE("/devices/:id", h.Remove)
}

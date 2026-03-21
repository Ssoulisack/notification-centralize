package routes

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/your-org/notification-center/internal/api/handler"
	"github.com/your-org/notification-center/internal/store"
)

func NewPreferenceRoutes(rg *gin.RouterGroup, prefStore store.PreferenceStore, logger *slog.Logger) {
	h := handler.NewPreferenceHandler(prefStore, logger)

	rg.GET("/preferences/:user_id", h.Get)
	rg.PUT("/preferences/:user_id", h.Update)
}

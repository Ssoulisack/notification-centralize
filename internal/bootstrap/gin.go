package bootstrap

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/your-org/notification-center/internal/config"
)

func NewRouter(cfg *config.Config, logger *slog.Logger) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)
	router := gin.New()
	router.Use(gin.Recovery())

	return router
}

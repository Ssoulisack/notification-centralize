package bootstrap

import (
	"log/slog"
	"os"

	"github.com/your-org/notification-center/internal/config"
)

func NewConfig(logger *slog.Logger) *config.Config {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	return cfg
}

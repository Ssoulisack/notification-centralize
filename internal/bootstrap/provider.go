package bootstrap

import (
	"log/slog"

	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/provider"
	emailprov "github.com/your-org/notification-center/internal/provider/email"
	smsprov "github.com/your-org/notification-center/internal/provider/sms"
	pushprov "github.com/your-org/notification-center/internal/provider/push"
	slackprov "github.com/your-org/notification-center/internal/provider/slack"
	telegramprov "github.com/your-org/notification-center/internal/provider/telegram"
	lineprov "github.com/your-org/notification-center/internal/provider/line"
)

func NewProviders(cfg *config.Config, logger *slog.Logger) *provider.Registry {
	registry := provider.NewRegistry()

	if cfg.Provider.Email.Enabled {
		registry.Register(emailprov.New(cfg.Provider.Email))
		logger.Info("email provider enabled")
	}
	if cfg.Provider.SMS.Enabled {
		registry.Register(smsprov.New(cfg.Provider.SMS))
		logger.Info("sms provider enabled")
	}
	if cfg.Provider.Push.Enabled {
		registry.Register(pushprov.New(cfg.Provider.Push))
		logger.Info("push provider enabled")
	}
	if cfg.Provider.Slack.Enabled {
		registry.Register(slackprov.New(cfg.Provider.Slack))
		logger.Info("slack provider enabled")
	}
	if cfg.Provider.Telegram.Enabled {
		registry.Register(telegramprov.New(cfg.Provider.Telegram))
		logger.Info("telegram provider enabled")
	}
	if cfg.Provider.Line.Enabled {
		registry.Register(lineprov.New(cfg.Provider.Line))
		logger.Info("line provider enabled")
	}

	return registry
}

package telegram

import (
	"context"
	"fmt"

	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/model"
)

type Provider struct {
	cfg config.TelegramProviderConfig
}

func New(cfg config.TelegramProviderConfig) *Provider {
	return &Provider{cfg: cfg}
}

func (p *Provider) Channel() model.Channel { return model.ChannelTelegram }
func (p *Provider) Name() string           { return "telegram" }

func (p *Provider) Send(ctx context.Context, n *model.Notification) error {
	// TODO: implement with github.com/go-telegram-bot-api/telegram-bot-api/v5
	fmt.Printf("[telegram] sending to %s: %s\n", n.Recipient, n.Body)
	return nil
}

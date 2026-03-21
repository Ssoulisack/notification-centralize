package line

import (
	"context"
	"fmt"

	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/model"
)

type Provider struct {
	cfg config.LineProviderConfig
}

func New(cfg config.LineProviderConfig) *Provider {
	return &Provider{cfg: cfg}
}

func (p *Provider) Channel() model.Channel { return model.ChannelLine }
func (p *Provider) Name() string           { return "line" }

func (p *Provider) Send(ctx context.Context, n *model.Notification) error {
	// TODO: implement with github.com/line/line-bot-sdk-go
	fmt.Printf("[line] sending to %s: %s\n", n.Recipient, n.Body)
	return nil
}

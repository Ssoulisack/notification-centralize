package slack

import (
	"context"
	"fmt"

	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/model"
)

// Provider sends notifications via Slack.
type Provider struct {
	cfg config.SlackProviderConfig
}

func New(cfg config.SlackProviderConfig) *Provider {
	return &Provider{cfg: cfg}
}

func (p *Provider) Channel() model.Channel { return model.ChannelSlack }
func (p *Provider) Name() string           { return "slack" }

func (p *Provider) Send(ctx context.Context, n *model.Notification) error {
	// TODO: implement with github.com/slack-go/slack
	//
	// client := slack.New(p.cfg.Token)
	// _, _, err := client.PostMessageContext(ctx, n.Recipient,
	//     slack.MsgOptionText(n.Body, false),
	// )
	// return err

	fmt.Printf("[slack] sending to %s: %s\n", n.Recipient, n.Body)
	return nil
}

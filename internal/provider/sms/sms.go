package sms

import (
	"context"
	"fmt"

	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/model"
)

// Provider sends notifications via Twilio SMS.
type Provider struct {
	cfg config.SMSProviderConfig
}

func New(cfg config.SMSProviderConfig) *Provider {
	return &Provider{cfg: cfg}
}

func (p *Provider) Channel() model.Channel { return model.ChannelSMS }
func (p *Provider) Name() string           { return "twilio" }

func (p *Provider) Send(ctx context.Context, n *model.Notification) error {
	// TODO: implement with github.com/twilio/twilio-go
	//
	// client := twilio.NewRestClientWithParams(twilio.ClientParams{
	//     Username: p.cfg.AccountSID,
	//     Password: p.cfg.AuthToken,
	// })
	// params := &api.CreateMessageParams{}
	// params.SetTo(n.Recipient)
	// params.SetFrom(p.cfg.FromNumber)
	// params.SetBody(n.Body)
	// _, err := client.Api.CreateMessage(params)
	// return err

	fmt.Printf("[sms] sending to %s: %s\n", n.Recipient, n.Body)
	return nil
}

package email

import (
	"context"
	"fmt"

	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/model"
)

// Provider sends notifications via SMTP email.
type Provider struct {
	cfg config.EmailProviderConfig
}

// New creates an email provider from config.
func New(cfg config.EmailProviderConfig) *Provider {
	return &Provider{cfg: cfg}
}

func (p *Provider) Channel() model.Channel { return model.ChannelEmail }
func (p *Provider) Name() string           { return "smtp" }

func (p *Provider) Send(ctx context.Context, n *model.Notification) error {
	// TODO: implement with github.com/wneessen/go-mail
	//
	// m := mail.NewMsg()
	// m.From(p.cfg.From)
	// m.To(n.Recipient)
	// m.Subject(n.Subject)
	// m.SetBodyString(mail.TypeTextHTML, n.Body)
	//
	// client, err := mail.NewClient(p.cfg.Host,
	//     mail.WithPort(p.cfg.Port),
	//     mail.WithSMTPAuth(mail.SMTPAuthPlain),
	//     mail.WithUsername(p.cfg.Username),
	//     mail.WithPassword(p.cfg.Password),
	// )
	// return client.DialAndSendWithContext(ctx, m)

	fmt.Printf("[email] sending to %s: %s\n", n.Recipient, n.Subject)
	return nil
}

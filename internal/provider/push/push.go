package push

import (
	"context"
	"fmt"

	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/model"
)

// Provider sends push notifications via FCM and APNs.
type Provider struct {
	cfg config.PushProviderConfig
}

func New(cfg config.PushProviderConfig) *Provider {
	return &Provider{cfg: cfg}
}

func (p *Provider) Channel() model.Channel { return model.ChannelPush }
func (p *Provider) Name() string           { return "fcm+apns" }

func (p *Provider) Send(ctx context.Context, n *model.Notification) error {
	// TODO: detect platform from metadata and route to FCM or APNs
	//
	// For FCM (Android):
	//   app, _ := firebase.NewApp(ctx, nil, option.WithCredentialsFile(p.cfg.FCM.CredentialsFile))
	//   client, _ := app.Messaging(ctx)
	//   _, err := client.Send(ctx, &messaging.Message{
	//       Token: n.Recipient,
	//       Notification: &messaging.Notification{Title: n.Subject, Body: n.Body},
	//   })
	//
	// For APNs (iOS):
	//   cert, _ := certificate.FromPemFile(p.cfg.APNs.CertFile, "")
	//   client := apns2.NewClient(cert)
	//   resp, err := client.Push(&apns2.Notification{
	//       DeviceToken: n.Recipient,
	//       Topic:       "com.yourapp",
	//       Payload:     payload.NewPayload().Alert(n.Body),
	//   })

	fmt.Printf("[push] sending to device %s: %s\n", n.Recipient, n.Subject)
	return nil
}

# Adding New Providers

## Steps to Add a New Channel

### 1. Define Channel Type

In `internal/model/notification.go`, add channel constant:
```go
const (
    ChannelEmail    Channel = "email"
    ChannelWhatsApp Channel = "whatsapp"  // new
)
```

### 2. Add Config

In `internal/config/config.go`:
```go
type ProviderConfig struct {
    // existing...
    WhatsApp WhatsAppProviderConfig `mapstructure:"whatsapp"`
}

type WhatsAppProviderConfig struct {
    Enabled   bool   `mapstructure:"enabled"`
    AccountID string `mapstructure:"account_id"`
    Token     string `mapstructure:"token"`
}
```

### 3. Implement Provider

Create `internal/provider/whatsapp/whatsapp.go`:
```go
package whatsapp

type Provider struct {
    cfg config.WhatsAppProviderConfig
}

func New(cfg config.WhatsAppProviderConfig) *Provider {
    return &Provider{cfg: cfg}
}

func (p *Provider) Send(ctx context.Context, n *model.Notification) error {
    // Implementation using WhatsApp Business API
    return nil
}

func (p *Provider) Channel() model.Channel {
    return model.ChannelWhatsApp
}

func (p *Provider) Name() string {
    return "whatsapp"
}
```

### 4. Register in Bootstrap

In `internal/bootstrap/app.go`, add to `newProviders()`:
```go
if cfg.Provider.WhatsApp.Enabled {
    registry.Register(whatsappprov.New(cfg.Provider.WhatsApp))
    logger.Info("whatsapp provider enabled")
}
```

### 5. Update Config File

In `config.yaml`:
```yaml
providers:
  whatsapp:
    enabled: true
    account_id: ${WHATSAPP_ACCOUNT_ID}
    token: ${WHATSAPP_TOKEN}
```

## Provider Interface

All providers must implement:
```go
type Provider interface {
    Send(ctx context.Context, n *Notification) error
    Channel() Channel
    Name() string
}
```

The registry looks up providers by channel type when workers process notifications.

# Architecture Guide

## System Flow

```
Client → REST API → NotificationService → Queue.Publish()
                                              ↓
                    Worker.consume() ← Queue.Subscribe()
                          ↓
                    Provider.Send()
                          ↓
                    External Service (SMTP, Twilio, FCM, etc.)
```

## Package Structure

### `internal/bootstrap/`
Central DI container. `Application` struct holds all dependencies:
```go
type Application struct {
    Config   *config.Config
    Logger   *slog.Logger
    DB       *pgxpool.Pool
    Cache    *rediscache.Cache
    Queue    queue.Queue
    Registry *provider.Registry
    Router   *gin.Engine
    Stores   *Stores
}
```

### `internal/service/`
Business logic layer. Two send modes:
- **Direct**: Specify channel + recipient explicitly
- **Event-based**: Resolve from user preferences, fan-out to multiple channels

### `internal/provider/`
Pluggable channel adapters implementing:
```go
type Provider interface {
    Send(ctx context.Context, n *Notification) error
    Channel() Channel
    Name() string
}
```

### `internal/queue/`
Message broker abstraction:
```go
type Queue interface {
    Publish(n *Notification) error
    Subscribe() <-chan *Notification
    Close() error
}
```
Implementations: `memory/` (dev), `rabbitmq/` (prod)

### `internal/worker/`
Consumes from queue with:
- Idempotency check (Redis SetNX, 24h TTL)
- Rate limiting (sliding window per provider)
- Retry with backoff (configurable max retries)

### `internal/store/`
Repository pattern with interfaces:
- `NotificationStore` - Notification CRUD
- `DeviceStore` - Push token management
- `PreferenceStore` - User notification settings
- `TemplateStore` - Message templates

## Database Schema

4 core tables in PostgreSQL:
- `notifications` - Delivery log with JSONB metadata
- `device_tokens` - Push notification registrations
- `user_preferences` - Channel opt-ins, quiet hours
- `notification_templates` - Go template subject/body

## Configuration Hierarchy

```
config.yaml → Environment Variables → Defaults
```

Viper loads config with `mapstructure` tags. Key sections:
- `server` - Host, port, gin mode
- `queue` - Engine selection, RabbitMQ/NATS config
- `worker` - Concurrency, retry settings
- `providers` - Per-channel enable flags and credentials

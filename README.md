# Notification Center

A multi-channel notification microservice built in Go. Other services send notification requests via REST API, and this service handles routing, queuing, retrying, and delivering through email, SMS, push, Slack, Telegram, and Line.

## Features

- **Multi-Project Support**: Isolate data between projects/tenants using API keys
- **Multi-Channel Delivery**: Email, SMS, Push, Slack, Telegram, Line
- **Async Processing**: Queue-based workers with retry logic
- **User Preferences**: Per-user channel settings and opt-outs
- **Templating**: Reusable notification templates with variable substitution

## Architecture

```
Services ──▶ REST API ──▶ Queue ──▶ Workers ──▶ Providers (Email/SMS/Push/Slack/...)
                │                      │
                ▼                      ▼
           PostgreSQL              Redis (cache + rate limit)
```

## Quick Start

```bash
# 1. Clone and configure
cp .env.example .env
# Edit .env with your credentials

# 2. Start infrastructure
make docker-up

# 3. Run migrations
migrate -path migrations -database "$DATABASE_URL" up

# 4. Run the server
make run

# 5. Send a test notification
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "X-API-Key: default-api-key" \
  -H "Content-Type: application/json" \
  -d '{"channel":"email","recipient":"test@example.com","body":"Hello"}'
```

## Multi-Project Support

Each project has its own API key and isolated data. All API requests require the `X-API-Key` header.

### Create a new project

```sql
INSERT INTO projects (id, name, api_key)
VALUES ('proj-123', 'My App', 'my-secret-api-key');
```

### Use the API key in requests

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "X-API-Key: my-secret-api-key" \
  -H "Content-Type: application/json" \
  -d '{"channel":"email","recipient":"user@example.com","body":"Hello"}'
```

Data isolation is enforced at the database level - notifications, devices, preferences, and templates are scoped to each project.

## API Endpoints

All endpoints require the `X-API-Key` header.

### Send notification (direct)
```bash
POST /api/v1/notifications
X-API-Key: your-api-key

{
  "channel": "email",
  "recipient": "user@example.com",
  "subject": "Hello",
  "body": "World",
  "priority": "normal"
}
```

### Send notification (event-based)
```bash
POST /api/v1/notifications
X-API-Key: your-api-key

{
  "user_id": "user_123",
  "event": "order_shipped",
  "template_id": "tmpl_order_shipped_email",
  "data": {"order_id": "ORD-001", "name": "John"}
}
```

### Register device token
```bash
POST /api/v1/devices
X-API-Key: your-api-key

{
  "user_id": "user_123",
  "token": "fcm_token_abc",
  "platform": "android"
}
```

### Check notification status
```bash
GET /api/v1/notifications/:id
X-API-Key: your-api-key
```

### User preferences
```bash
GET /api/v1/preferences/:user_id
PUT /api/v1/preferences/:user_id
X-API-Key: your-api-key

{
  "enabled_channels": ["email", "push"],
  "quiet_start": 22,
  "quiet_end": 8,
  "opted_out_events": ["marketing"]
}
```

## Project Structure

```
cmd/server/main.go              Entry point
internal/
  api/context/                  Request context helpers
  api/handler/                  HTTP handlers
  api/middleware/               Auth, logging, CORS
  api/routes/                   Route registration
  bootstrap/                    Dependency injection
  config/                       YAML + env config
  model/                        Domain models + DTOs
  provider/                     Channel implementations
    email/ sms/ push/ slack/ telegram/ line/
  job/                          Message queue abstraction
    memory/                     In-memory (dev)
    rabbitmq/                   RabbitMQ (production)
  service/                      Business logic
  store/                        Persistence interfaces
    postgres/                   PostgreSQL implementations
    redis/                      Cache + rate limiter
  template/                     Notification template engine
  worker/                       Async queue consumers
migrations/                     SQL schema files
deployments/                    Docker Compose
```

## Tech Stack

| Layer      | Technology                          |
|------------|-------------------------------------|
| Language   | Go 1.22+                            |
| HTTP       | Gin                                 |
| Database   | PostgreSQL 16 + pgx                 |
| Cache      | Redis 7                             |
| Queue      | In-memory / RabbitMQ                |
| Email      | go-mail (SMTP)                      |
| SMS        | Twilio                              |
| Push       | Firebase (FCM) + APNs               |
| Chat       | Slack, Telegram, Line               |
| Config     | Viper                               |
| Logging    | slog (stdlib)                       |

## Database Schema

### Projects
```sql
CREATE TABLE projects (
    id          VARCHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    api_key     VARCHAR(255) NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

All other tables (`notifications`, `device_tokens`, `user_preferences`, `notification_templates`) include a `project_id` foreign key for data isolation.

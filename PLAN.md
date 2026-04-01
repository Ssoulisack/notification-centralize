# Notification Center Service - Implementation Plan

## Overview
A notification center service using Go + Gin, PostgreSQL, Keycloak, and RabbitMQ. Users register via Keycloak, join projects with roles, and receive notifications through multiple channels.

## Tech Stack
- **Backend**: Go + Gin
- **Database**: PostgreSQL
- **Auth**: Keycloak (JWT)
- **Message Queue**: RabbitMQ
- **API Key**: Auto-generated on registration

---

## Implementation Status

### Completed Files

| File | Status | Description |
|------|--------|-------------|
| `migrations/001_initial_schema.sql` | ✅ Done | Complete database schema with all tables |
| `pkg/database/postgres.go` | ✅ Done | PostgreSQL connection pool client |
| `pkg/messaging/rabbitmq.go` | ✅ Done | RabbitMQ client with exchanges/queues |
| `internal/models/models.go` | ✅ Done | All data models and enums |
| `internal/config/config.go` | ✅ Done | Configuration management |
| `internal/middleware/auth.go` | ✅ Done | Keycloak JWT validation |
| `internal/middleware/api_key.go` | ✅ Done | API key authentication |
| `internal/services/user_sync.go` | ✅ Done | Keycloak user synchronization |
| `internal/services/notification_service.go` | ✅ Done | Notification business logic |
| `internal/repository/project_repository.go` | ✅ Done | Project database operations |
| `internal/repository/api_key_repository.go` | ✅ Done | API key database operations |
| `internal/handlers/*.go` | ✅ Done | All HTTP handlers |
| `internal/workers/notification_worker.go` | ✅ Done | Background notification workers |
| `cmd/api/main.go` | ✅ Done | API server entry point |
| `cmd/worker/main.go` | ✅ Done | Worker entry point |

### Remaining Files

| File | Status | Description |
|------|--------|-------------|
| `docker-compose.yml` | ⏳ Pending | Docker services configuration |
| `Dockerfile` | ⏳ Pending | Multi-stage build for API and worker |
| `go.mod` | ⏳ Pending | Go module definition |
| `go.sum` | ⏳ Pending | Go dependencies checksum |
| `.env.example` | ⏳ Pending | Environment variables template |

---

## Database Schema

### Tables
| Table | Purpose |
|-------|---------|
| `users` | User profiles synced from Keycloak |
| `projects` | Projects that contain users |
| `roles` | Role definitions (owner, admin, member, viewer) |
| `project_members` | Links users to projects with roles |
| `api_keys` | API keys per user/project |
| `notification_templates` | Reusable notification templates |
| `notifications` | Notification records |
| `notification_recipients` | Delivery status per user/channel |
| `notification_events` | Audit log for notification lifecycle |
| `device_tokens` | Push notification device registration |
| `webhook_endpoints` | Webhook delivery configuration |
| `webhook_deliveries` | Webhook delivery attempts |
| `user_notification_preferences` | Per-user notification settings |

### Role Permissions
```json
{
  "projects": { "read": true, "update": false, "delete": false, "manage_members": false },
  "notifications": { "read": true, "create": true, "send_bulk": false },
  "api_keys": { "read": true, "create": false, "revoke": false },
  "analytics": { "view": true, "export": false }
}
```

---

## API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/login` | Redirect to Keycloak |
| POST | `/auth/callback` | Handle Keycloak callback |
| POST | `/auth/logout` | Logout user |
| GET | `/auth/me` | Get current user (syncs from Keycloak) |

### Users
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/users/me` | Get current user |
| PATCH | `/users/me` | Update profile |
| GET | `/users/me/projects` | List user's projects |
| GET | `/users/me/notifications` | User inbox |

### Projects
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/projects` | List projects |
| POST | `/projects` | Create project (user becomes owner) |
| GET | `/projects/:id` | Get project |
| PATCH | `/projects/:id` | Update project |
| DELETE | `/projects/:id` | Delete project |

### Project Members
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/projects/:id/members` | List members |
| POST | `/projects/:id/members` | Add member (auto-generates API key) |
| PATCH | `/projects/:id/members/:memberId` | Update role |
| DELETE | `/projects/:id/members/:memberId` | Remove member |

### API Keys
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/projects/:id/api-keys` | List API keys |
| POST | `/projects/:id/api-keys` | Create API key |
| DELETE | `/projects/:id/api-keys/:keyId` | Revoke key |

### Notifications
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/projects/:id/notifications` | Send notification |
| POST | `/projects/:id/notifications/batch` | Batch send |
| GET | `/projects/:id/notifications` | List notifications |
| GET | `/projects/:id/notifications/:notificationId` | Get notification |

### User Inbox
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/inbox` | Get notifications |
| GET | `/inbox/unread/count` | Unread count |
| POST | `/inbox/:id/read` | Mark as read |

---

## RabbitMQ Architecture

### Exchanges
| Exchange | Type | Purpose |
|----------|------|---------|
| `notifications.commands` | direct | Route to channel queues |
| `notifications.events` | topic | Tracking events |
| `notifications.dlx` | fanout | Dead letter exchange |

### Queues
| Queue | Purpose |
|-------|---------|
| `notifications.send.in_app` | In-app notifications |
| `notifications.send.email` | Email delivery |
| `notifications.send.sms` | SMS delivery |
| `notifications.send.push` | Push notifications |
| `notifications.dlq` | Failed messages |

### Message Flow
```
API Request → notifications.commands → Channel Queue → Worker → Delivery
                                                          ↓
                                            notifications.events → Analytics
```

---

## Project Structure
```
notification-service/
├── cmd/
│   ├── api/main.go           # API server
│   └── worker/main.go        # Background workers
├── internal/
│   ├── config/config.go      # Configuration
│   ├── handlers/             # HTTP handlers
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── project_handler.go
│   │   ├── notification_handler.go
│   │   ├── inbox_handler.go
│   │   ├── api_key_handler.go
│   │   └── common.go
│   ├── middleware/           # Auth middleware
│   │   ├── auth.go           # Keycloak JWT
│   │   └── api_key.go        # API key auth
│   ├── models/models.go      # Data models
│   ├── repository/           # Database queries
│   │   ├── project_repository.go
│   │   └── api_key_repository.go
│   ├── services/             # Business logic
│   │   ├── user_sync.go
│   │   └── notification_service.go
│   └── workers/              # Channel workers
│       └── notification_worker.go
├── pkg/
│   ├── messaging/rabbitmq.go # RabbitMQ client
│   └── database/postgres.go  # PostgreSQL client
├── migrations/
│   └── 001_initial_schema.sql
├── docker-compose.yml        # Services config
├── Dockerfile               # Build config
├── go.mod
└── go.sum
```

---

## API Key Format

When a user joins a project:
1. Create `project_members` record
2. Generate 32-byte random key
3. Hash with SHA-256 for storage
4. Return full key once (never stored in plaintext)

Key format: `nc_live_xxxxxxxxxxxxxxxxxxxx`

---

## Keycloak Integration

1. **User Login**: Frontend redirects to Keycloak
2. **JWT Validation**: Backend validates token from Keycloak JWKS
3. **User Sync**: On first login, create user record in PostgreSQL
4. **Claims Used**: `sub`, `email`, `preferred_username`, `given_name`, `family_name`

---

## Environment Variables

```env
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=notification_center
DB_SSL_MODE=disable

# RabbitMQ
RABBITMQ_URL=amqp://user:password@localhost:5672/

# Keycloak
KEYCLOAK_BASE_URL=http://localhost:8180
KEYCLOAK_REALM=notification-center
KEYCLOAK_CLIENT_ID=notification-api
KEYCLOAK_CLIENT_SECRET=your-secret
```

---

## Verification Steps

1. Start services: `docker-compose up -d`
2. Access Keycloak: `http://localhost:8180`
3. Create realm "notification-center"
4. Create client "notification-api"
5. Create user in Keycloak
6. Login via API, verify JWT validation
7. Create project, verify API key generated
8. Send notification, verify RabbitMQ message
9. Check user inbox for notification

---

## Next Steps

1. Create `docker-compose.yml` with PostgreSQL, RabbitMQ, Keycloak
2. Create `Dockerfile` for multi-stage build
3. Create `go.mod` and `go.sum`
4. Create `.env.example`
5. Test the complete flow

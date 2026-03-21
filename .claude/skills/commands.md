# Commands Reference

## Build & Run

```bash
make build              # go build -ldflags="-s -w" -o bin/notification-center ./cmd/server
make run                # go run ./cmd/server
```

## Testing

```bash
make test               # go test -v -race -cover ./...
go test ./internal/service/...    # Test specific package
go test -run TestSend ./...       # Run specific test
```

## Linting

```bash
make lint               # golangci-lint run ./...
```

## Docker

```bash
make docker-up          # docker compose up -d --build
make docker-down        # docker compose down -v
make docker-logs        # docker compose logs -f
```

## Database

```bash
make migrate-up         # Apply migrations
make migrate-down       # Rollback migrations
```

## Manual Testing

```bash
# Send notification
make test-send

# Or via curl:
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{"channel":"email","recipient":"user@example.com","subject":"Test","body":"Hello"}'

# Register device
make test-device
```

## Environment

Required services for full functionality:
- PostgreSQL (required)
- Redis (optional, degrades gracefully)
- RabbitMQ (optional, uses memory queue if not configured)

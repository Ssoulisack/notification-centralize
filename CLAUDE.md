# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Notification Center is a Go microservice for multi-channel notification delivery (Email, SMS, Push, Slack, Telegram, Line) with async processing, user preferences, and templating.

## Quick Commands

```bash
make build          # Build binary to bin/
make run            # Run server
make test           # Run all tests with race detector
make lint           # golangci-lint
make docker-up      # Start PostgreSQL, Redis via docker-compose
```

## Architecture

```
API Request → Service → Queue → Workers → Providers
                ↓                   ↓
            PostgreSQL           Redis (cache/rate-limit)
```

**Key packages:**
- `cmd/server/` - Entry point
- `internal/bootstrap/` - DI container, Application struct
- `internal/service/` - Business logic
- `internal/provider/` - Channel implementations (pluggable)
- `internal/queue/` - Message queue abstraction (memory/RabbitMQ)
- `internal/store/` - Data access layer (PostgreSQL/Redis)
- `internal/worker/` - Async notification processor

## Key Patterns

- **Interface-based design**: Provider, Queue, Store interfaces allow swapping implementations
- **Bootstrap pattern**: `bootstrap.App(ctx)` initializes all dependencies
- **Registry pattern**: Providers registered by channel type

## Configuration

Set `queue.engine` in config.yaml:
- `memory` - Development (default)
- `rabbitmq` - Production

## Skills

See `.claude/skills/` for detailed guides:
- `architecture.md` - Full system architecture
- `commands.md` - All available commands
- `providers.md` - Adding new notification channels

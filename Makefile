.PHONY: build build-server build-worker run run-server run-worker run-all \
        test lint tidy \
        docker-up docker-down docker-logs \
        migrate-up migrate-down \
        encrypt decrypt

APP_NAME   = notification-center
SERVER_BIN = bin/server
WORKER_BIN = bin/worker

SOPS    = sops
ENCRYPT = $(SOPS) --encrypt --output config.enc.yaml config.yaml
DECRYPT = $(SOPS) --decrypt --output config.yaml config.enc.yaml

## ── Build ────────────────────────────────────────────────────────────────────

build: build-server build-worker

build-server:
	go build -ldflags="-s -w" -o $(SERVER_BIN) ./cmd/api

build-worker:
	go build -ldflags="-s -w" -o $(WORKER_BIN) ./cmd/worker

## ── Run ──────────────────────────────────────────────────────────────────────

# API server only
run-server:
	go run ./cmd/api

# Background worker only
run-worker:
	go run ./cmd/worker

# Both server and worker (logs interleaved — use separate terminals for cleaner output)
run-all:
	@trap 'kill 0' INT TERM; go run ./cmd/api & go run ./cmd/worker & wait

# Default: server only
run: run-server

## ── Test / Lint ──────────────────────────────────────────────────────────────

test:
	go test -v -race -cover ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

## ── Docker ───────────────────────────────────────────────────────────────────

docker-up:
	docker compose -f deployments/docker-compose.yaml up -d --build

docker-down:
	docker compose -f deployments/docker-compose.yaml down -v

docker-logs:
	docker compose -f deployments/docker-compose.yaml logs -f

## ── Database migrations (requires golang-migrate CLI) ────────────────────────

# Override with: make migrate-up DB_URL="postgres://user:pass@host:port/db?sslmode=disable"
DB_URL ?= postgres://lotto_admin:ithq@2026@10.150.1.85:30432/notification_center?sslmode=disable

migrate-up:
	migrate -path bootstrap/database/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path bootstrap/database/migrations -database "$(DB_URL)" down

## ── Secrets ──────────────────────────────────────────────────────────────────

encrypt:
	$(ENCRYPT)

decrypt:
	$(DECRYPT)
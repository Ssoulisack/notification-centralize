.PHONY: build run test lint docker-up docker-down migrate-up migrate-down

APP_NAME = notification-center
MAIN     = ./cmd/server

SOPS=sops
ENCRYPT=$(SOPS) --encrypt --output config.enc.yaml config.yaml
DECRYPT=$(SOPS) --decrypt --output config.yaml config.enc.yaml

# Build the binary
build:
	go build -ldflags="-s -w" -o bin/$(APP_NAME) $(MAIN)

# Run locally
run:
	go run $(MAIN)

# Run tests
test:
	go test -v -race -cover ./...

# Lint
lint:
	golangci-lint run ./...

# Docker
docker-up:
	docker compose -f deployments/docker-compose.yaml up -d --build

docker-down:
	docker compose -f deployments/docker-compose.yaml down -v

docker-logs:
	docker compose -f deployments/docker-compose.yaml logs -f app

# Database migrations (requires golang-migrate CLI)
migrate-up:
	migrate -path migrations -database "postgres://notify:notify_secret@localhost:5432/notification_center?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://notify:notify_secret@localhost:5432/notification_center?sslmode=disable" down

# Generate proto (if using gRPC)
proto:
	protoc --go_out=. --go-grpc_out=. proto/*.proto

# Quick test: send a notification
test-send:
	curl -X POST http://localhost:8080/api/v1/notifications \
		-H "Content-Type: application/json" \
		-H "X-API-Key: $${API_KEY}" \
		-d '{ \
			"channel": "email", \
			"recipient": "test@example.com", \
			"subject": "Test Notification", \
			"body": "Hello from notification center!", \
			"priority": "normal" \
		}'

# Quick test: register a device
test-device:
	curl -X POST http://localhost:8080/api/v1/devices \
		-H "Content-Type: application/json" \
		-H "X-API-Key: $${API_KEY}" \
		-d '{ \
			"user_id": "user_123", \
			"token": "fcm_device_token_abc", \
			"platform": "android", \
			"app_version": "1.0.0" \
		}'

tidy:
	$(GO_MOD_TIDY)

swag:
	$(SWAG_DOC)

encrypt:
	$(ENCRYPT)

decrypt:
	$(DECRYPT)
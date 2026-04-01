# Build stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build API
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/api ./cmd/api

# Build Worker
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/worker ./cmd/worker

# API runtime stage
FROM alpine:3.19 AS api

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /bin/api .
COPY config.yaml .
COPY migrations/ ./migrations/

EXPOSE 8080

ENTRYPOINT ["./api"]

# Worker runtime stage
FROM alpine:3.19 AS worker

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /bin/worker .
COPY config.yaml .

ENTRYPOINT ["./worker"]

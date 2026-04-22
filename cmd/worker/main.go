package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/your-org/notification-center/bootstrap/database"
	"github.com/your-org/notification-center/bootstrap/messaging"
	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/data/repository"
	"github.com/your-org/notification-center/internal/workers"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load("")
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gormDB, err := database.New(cfg.Database.DSN())
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	logger.Info("connected to database")

	rabbitmq, err := messaging.NewClient(ctx, messaging.Config{
		URL:           cfg.RabbitMQ.URL,
		PrefetchCount: cfg.RabbitMQ.PrefetchCount,
	}, logger)
	if err != nil {
		logger.Error("failed to connect to RabbitMQ", "error", err)
		os.Exit(1)
	}
	defer rabbitmq.Close()
	logger.Info("connected to RabbitMQ")

	recipientRepo := repository.NewRecipientRepository(gormDB)
	worker := workers.NewNotificationWorker(rabbitmq, recipientRepo, &cfg.Worker, logger)

	if err := worker.Start(ctx); err != nil {
		logger.Error("failed to start worker", "error", err)
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("shutting down worker...")
	worker.Stop()
	cancel()
	logger.Info("worker stopped")
}

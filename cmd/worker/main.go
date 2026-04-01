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
	"github.com/your-org/notification-center/internal/workers"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database
	dbClient, err := database.NewClient(ctx, database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
		MaxConns: cfg.Database.MaxConns,
		MinConns: cfg.Database.MinConns,
	})
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbClient.Close()

	logger.Info("connected to database")

	// Initialize RabbitMQ
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

	// Initialize worker
	worker := workers.NewNotificationWorker(rabbitmq, dbClient.Pool, &cfg.Worker, logger)

	// Start worker
	if err := worker.Start(ctx); err != nil {
		logger.Error("failed to start worker", "error", err)
		os.Exit(1)
	}

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("shutting down worker...")

	// Stop worker gracefully
	worker.Stop()

	cancel()

	logger.Info("worker stopped")
}

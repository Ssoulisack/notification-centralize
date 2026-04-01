package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/your-org/notification-center/bootstrap/database"
	"github.com/your-org/notification-center/bootstrap/messaging"
	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/middleware"
	"github.com/your-org/notification-center/internal/routes"
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

	// Set Gin mode
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
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

	// Initialize auth middleware
	authMiddleware, err := middleware.NewAuthMiddleware(&cfg.Keycloak, logger)
	if err != nil {
		logger.Error("failed to initialize auth middleware", "error", err)
		os.Exit(1)
	}
	defer authMiddleware.Close()

	// Initialize API key middleware
	apiKeyMiddleware := middleware.NewAPIKeyMiddleware(dbClient.Pool, logger)

	// Setup router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-API-Key")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Setup routes
	deps := &routes.Dependencies{
		DB:               dbClient.Pool,
		RabbitMQ:         rabbitmq,
		Config:           cfg,
		Logger:           logger,
		AuthMiddleware:   authMiddleware,
		APIKeyMiddleware: apiKeyMiddleware,
	}
	routes.Setup(router, deps)

	// Start server
	server := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		logger.Info("shutting down server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("server shutdown error", "error", err)
		}

		cancel()
	}()

	logger.Info("starting API server", "address", cfg.Server.Address())

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}

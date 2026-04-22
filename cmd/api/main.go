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
	"github.com/your-org/notification-center/internal/routes"
	"github.com/your-org/notification-center/pkg/middleware"
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

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
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

	authMiddleware, err := middleware.NewAuthMiddleware(&cfg.Keycloak, logger)
	if err != nil {
		logger.Error("failed to initialize auth middleware", "error", err)
		os.Exit(1)
	}
	defer authMiddleware.Close()

	apiKeyMiddleware := middleware.NewAPIKeyMiddleware(gormDB, logger)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

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

	routes.Setup(router, &routes.Dependencies{
		GormDB:           gormDB,
		RabbitMQ:         rabbitmq,
		Config:           cfg,
		Logger:           logger,
		AuthMiddleware:   authMiddleware,
		APIKeyMiddleware: apiKeyMiddleware,
	})

	server := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

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

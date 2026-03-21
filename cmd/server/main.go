package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/your-org/notification-center/internal/api/routes"
	"github.com/your-org/notification-center/internal/bootstrap"
	"github.com/your-org/notification-center/internal/worker"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := bootstrap.App(ctx)
	defer app.Close()

	cfg := app.Config
	logger := app.Logger
	router := app.Router

	routes.New(router, app, cfg)

	// Workers
	w := worker.New(app.Queue, app.Registry, app.Stores.Notification, app.Cache, logger,
		cfg.Worker.MaxRetry, cfg.Worker.RetryDelay)
	wg := w.Start(ctx, cfg.Worker.Concurrency)
	logger.Info("workers started", "concurrency", cfg.Worker.Concurrency)

	// Server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: router}

	go func() {
		logger.Info("notification center running", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", "error", err)
	}

	wg.Wait()
	logger.Info("notification center stopped")
}

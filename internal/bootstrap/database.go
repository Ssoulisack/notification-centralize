package bootstrap

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/store/postgres"
	rediscache "github.com/your-org/notification-center/internal/store/redis"
)

func NewStorage(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*pgxpool.Pool, *rediscache.Cache) {
	dbPool, err := postgres.Connect(ctx, cfg.Database)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	logger.Info("connected to PostgreSQL")

	cache, err := rediscache.Connect(cfg.Redis)
	if err != nil {
		logger.Warn("redis not available, running without cache", "error", err)
		cache = nil
	} else {
		logger.Info("connected to Redis")
	}

	return dbPool, cache
}

func NewStores(dbPool *pgxpool.Pool) *Stores {
	return &Stores{
		Project:      postgres.NewProjectStore(dbPool),
		Notification: postgres.NewNotificationStore(dbPool),
		Device:       postgres.NewDeviceStore(dbPool),
		Preference:   postgres.NewPreferenceStore(dbPool),
		Template:     postgres.NewTemplateStore(dbPool),
	}
}

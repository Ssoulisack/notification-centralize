package bootstrap

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/provider"
	"github.com/your-org/notification-center/internal/job"
	"github.com/your-org/notification-center/internal/store/postgres"
	rediscache "github.com/your-org/notification-center/internal/store/redis"
)

type Application struct {
	Config   *config.Config
	Logger   *slog.Logger
	DB       *pgxpool.Pool
	Cache    *rediscache.Cache
	Queue    job.Queue
	Registry *provider.Registry
	Router   *gin.Engine
	Stores   *Stores
}

type Stores struct {
	Project      *postgres.ProjectStore
	Notification *postgres.NotificationStore
	Device       *postgres.DeviceStore
	Preference   *postgres.PreferenceStore
	Template     *postgres.TemplateStore
}

func App(ctx context.Context) *Application {
	app := &Application{}

	app.Logger = NewLogger()
	slog.SetDefault(app.Logger)

	app.Config = NewConfig(app.Logger)
	app.DB, app.Cache = NewStorage(ctx, app.Config, app.Logger)
	app.Stores = NewStores(app.DB)
	app.Registry = NewProviders(app.Config, app.Logger)
	app.Queue = NewJob(app.Config, app.Logger)
	app.Router = NewRouter(app.Config, app.Logger)

	return app
}

func (a *Application) Close() {
	if a.Queue != nil {
		a.Queue.Close()
	}
	if a.DB != nil {
		a.DB.Close()
	}
	if a.Cache != nil {
		a.Cache.Close()
	}
}
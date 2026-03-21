package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/your-org/notification-center/internal/api/middleware"
	"github.com/your-org/notification-center/internal/bootstrap"
	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/service"
	"github.com/your-org/notification-center/internal/template"
)

func New(r *gin.Engine, app *bootstrap.Application, cfg *config.Config) {
	logger := app.Logger

	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.ProjectAuth(app.Stores.Project, logger))

	// Service
	tmplEngine := template.New()
	svc := service.NewNotificationService(
		app.Stores.Notification,
		app.Stores.Device,
		app.Stores.Preference,
		app.Stores.Template,
		app.Queue,
		tmplEngine,
		logger,
	)

	// API v1
	v1 := r.Group("/api/v1")

	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})
	// Register routes
	NewNotificationRoutes(v1, svc, logger)
	NewDeviceRoutes(v1, app.Stores.Device, logger)
	NewPreferenceRoutes(v1, app.Stores.Preference, logger)
}

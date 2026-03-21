package bootstrap

import (
	"log/slog"
	"os"

	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/job"
	"github.com/your-org/notification-center/internal/job/memory"
	"github.com/your-org/notification-center/internal/job/rabbitmq"
)

func NewJob(cfg *config.Config, logger *slog.Logger) job.Queue {
	switch cfg.Queue.Engine {
	case "rabbitmq":
		mq, err := rabbitmq.New(cfg.Queue.RabbitMQ)
		if err != nil {
			logger.Error("failed to connect to RabbitMQ", "error", err)
			os.Exit(1)
		}
		logger.Info("connected to RabbitMQ", "queue", cfg.Queue.RabbitMQ.Queue)
		return mq
	default:
		logger.Info("using in-memory queue")
		return memory.New(cfg.Queue.BufferSize)
	}
}

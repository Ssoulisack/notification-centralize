package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/your-org/notification-center/internal/model"
	"github.com/your-org/notification-center/internal/provider"
	"github.com/your-org/notification-center/internal/job"
	"github.com/your-org/notification-center/internal/store"
	rediscache "github.com/your-org/notification-center/internal/store/redis"
)

// Worker consumes notifications from the queue and sends them via providers.
type Worker struct {
	queue      job.Queue
	registry   *provider.Registry
	notifStore store.NotificationStore
	cache      *rediscache.Cache
	logger     *slog.Logger
	maxRetry   int
	retryDelay time.Duration
}

// New creates a Worker.
func New(
	q job.Queue,
	reg *provider.Registry,
	notifStore store.NotificationStore,
	cache *rediscache.Cache,
	logger *slog.Logger,
	maxRetry int,
	retryDelay time.Duration,
) *Worker {
	return &Worker{
		queue:      q,
		registry:   reg,
		notifStore: notifStore,
		cache:      cache,
		logger:     logger,
		maxRetry:   maxRetry,
		retryDelay: retryDelay,
	}
}

// Start launches the specified number of worker goroutines.
func (w *Worker) Start(ctx context.Context, concurrency int) *sync.WaitGroup {
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			w.logger.Info("worker started", "worker_id", workerID)
			w.process(ctx, workerID)
		}(i)
	}

	return &wg
}

func (w *Worker) process(ctx context.Context, workerID int) {
	for {
		select {
		case <-ctx.Done():
			w.logger.Info("worker stopping", "worker_id", workerID)
			return

		case n, ok := <-w.queue.Subscribe():
			if !ok {
				return // queue closed
			}
			w.handle(ctx, n, workerID)
		}
	}
}

func (w *Worker) handle(ctx context.Context, n *model.Notification, workerID int) {
	log := w.logger.With("id", n.ID, "project_id", n.ProjectID, "channel", n.Channel, "worker_id", workerID)

	// Idempotency check
	if w.cache != nil {
		dup, err := w.cache.CheckIdempotency(ctx, n.ID)
		if err != nil {
			log.Error("idempotency check failed", "error", err)
		}
		if dup {
			log.Warn("duplicate notification, skipping")
			return
		}
	}

	// Rate limit check
	if w.cache != nil {
		p, _ := w.registry.Get(n.Channel)
		if p != nil {
			allowed, err := w.cache.AllowSend(ctx, p.Name(), 100, time.Minute)
			if err != nil {
				log.Error("rate limit check failed", "error", err)
			}
			if !allowed {
				log.Warn("rate limited, re-queuing")
				time.Sleep(w.retryDelay)
				_ = w.queue.Publish(n)
				return
			}
		}
	}

	// Update status to processing
	_ = w.notifStore.UpdateStatus(ctx, n.ProjectID, n.ID, model.StatusProcessing, "")

	// Find the provider
	p, ok := w.registry.Get(n.Channel)
	if !ok {
		log.Error("no provider registered for channel")
		_ = w.notifStore.UpdateStatus(ctx, n.ProjectID, n.ID, model.StatusFailed,
			"no provider for channel: "+string(n.Channel))
		return
	}

	// Send
	if err := p.Send(ctx, n); err != nil {
		log.Error("send failed", "attempt", n.RetryCount+1, "error", err)

		if n.RetryCount < w.maxRetry {
			n.RetryCount++
			n.Status = model.StatusQueued

			time.Sleep(w.retryDelay * time.Duration(n.RetryCount)) // exponential-ish backoff
			if pubErr := w.queue.Publish(n); pubErr != nil {
				log.Error("re-queue failed", "error", pubErr)
			}
			return
		}

		// Max retries exhausted
		_ = w.notifStore.UpdateStatus(ctx, n.ProjectID, n.ID, model.StatusFailed, err.Error())
		log.Error("notification permanently failed after max retries")
		return
	}

	// Success
	_ = w.notifStore.UpdateStatus(ctx, n.ProjectID, n.ID, model.StatusSent, "")
	log.Info("notification sent successfully")
}

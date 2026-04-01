package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/your-org/notification-center/bootstrap/messaging"
	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/models"
)

// NotificationWorker processes notification messages from RabbitMQ.
type NotificationWorker struct {
	rabbitmq *messaging.Client
	db       *pgxpool.Pool
	config   *config.WorkerConfig
	logger   *slog.Logger

	wg      sync.WaitGroup
	cancel  context.CancelFunc
	stopped bool
	mu      sync.Mutex
}

// NewNotificationWorker creates a new notification worker.
func NewNotificationWorker(
	rabbitmq *messaging.Client,
	db *pgxpool.Pool,
	cfg *config.WorkerConfig,
	logger *slog.Logger,
) *NotificationWorker {
	return &NotificationWorker{
		rabbitmq: rabbitmq,
		db:       db,
		config:   cfg,
		logger:   logger,
	}
}

// Start starts the worker pool.
func (w *NotificationWorker) Start(ctx context.Context) error {
	ctx, w.cancel = context.WithCancel(ctx)

	queues := []string{
		messaging.QueueInApp,
		messaging.QueueEmail,
		messaging.QueueSMS,
		messaging.QueuePush,
	}

	for _, queue := range queues {
		for i := 0; i < w.config.Concurrency; i++ {
			w.wg.Add(1)
			go w.consumeQueue(ctx, queue, fmt.Sprintf("%s-worker-%d", queue, i))
		}
	}

	w.logger.Info("notification worker started",
		"concurrency", w.config.Concurrency,
		"queues", queues,
	)

	return nil
}

// Stop gracefully stops the worker.
func (w *NotificationWorker) Stop() {
	w.mu.Lock()
	if w.stopped {
		w.mu.Unlock()
		return
	}
	w.stopped = true
	w.mu.Unlock()

	if w.cancel != nil {
		w.cancel()
	}

	w.wg.Wait()
	w.logger.Info("notification worker stopped")
}

// consumeQueue consumes messages from a queue.
func (w *NotificationWorker) consumeQueue(ctx context.Context, queue, consumerTag string) {
	defer w.wg.Done()

	deliveries, err := w.rabbitmq.Consume(ctx, queue, consumerTag)
	if err != nil {
		w.logger.Error("failed to start consuming",
			"queue", queue,
			"error", err,
		)
		return
	}

	w.logger.Info("started consuming",
		"queue", queue,
		"consumer", consumerTag,
	)

	for {
		select {
		case <-ctx.Done():
			return

		case delivery, ok := <-deliveries:
			if !ok {
				return
			}

			if err := w.processMessage(ctx, &delivery); err != nil {
				w.logger.Error("failed to process message",
					"queue", queue,
					"error", err,
				)
				delivery.Nack(false, true) // Requeue
			} else {
				delivery.Ack(false)
			}
		}
	}
}

// processMessage processes a single notification message.
func (w *NotificationWorker) processMessage(ctx context.Context, delivery *amqp.Delivery) error {
	var msg models.QueueMessage
	if err := json.Unmarshal(delivery.Body, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	w.logger.Info("processing notification",
		"notification_id", msg.NotificationID,
		"recipient_id", msg.RecipientID,
		"channel", msg.Channel,
	)

	// Update status to processing
	if err := w.updateRecipientStatus(ctx, msg.RecipientID, models.StatusPending, ""); err != nil {
		w.logger.Error("failed to update status",
			"recipient_id", msg.RecipientID,
			"error", err,
		)
	}

	// Process based on channel
	var err error
	switch msg.Channel {
	case models.ChannelInApp:
		err = w.processInApp(ctx, &msg)
	case models.ChannelEmail:
		err = w.processEmail(ctx, &msg)
	case models.ChannelSMS:
		err = w.processSMS(ctx, &msg)
	case models.ChannelPush:
		err = w.processPush(ctx, &msg)
	default:
		err = fmt.Errorf("unsupported channel: %s", msg.Channel)
	}

	if err != nil {
		// Check retry count
		if msg.RetryCount < w.config.MaxRetry {
			w.logger.Warn("retrying notification",
				"notification_id", msg.NotificationID,
				"retry_count", msg.RetryCount+1,
				"error", err,
			)

			// Increment retry count
			if updateErr := w.incrementRetryCount(ctx, msg.RecipientID); updateErr != nil {
				w.logger.Error("failed to update retry count", "error", updateErr)
			}

			// Requeue with delay
			time.Sleep(w.config.RetryDelay)
			return err // Return error to trigger nack and requeue
		}

		// Max retries exceeded
		w.logger.Error("notification failed after max retries",
			"notification_id", msg.NotificationID,
			"channel", msg.Channel,
			"error", err,
		)

		if updateErr := w.updateRecipientStatus(ctx, msg.RecipientID, models.StatusFailed, err.Error()); updateErr != nil {
			w.logger.Error("failed to update failed status", "error", updateErr)
		}

		// Publish failure event
		w.publishEvent(ctx, "notification.failed", &msg, err.Error())

		return nil // Don't return error - message is handled
	}

	// Success
	if updateErr := w.updateRecipientStatus(ctx, msg.RecipientID, models.StatusSent, ""); updateErr != nil {
		w.logger.Error("failed to update sent status", "error", updateErr)
	}

	w.publishEvent(ctx, "notification.sent", &msg, "")

	w.logger.Info("notification sent",
		"notification_id", msg.NotificationID,
		"channel", msg.Channel,
	)

	return nil
}

// processInApp handles in-app notification delivery.
func (w *NotificationWorker) processInApp(ctx context.Context, msg *models.QueueMessage) error {
	// In-app notifications are just stored in DB and marked as delivered
	if err := w.updateRecipientStatus(ctx, msg.RecipientID, models.StatusDelivered, ""); err != nil {
		return err
	}
	return nil
}

// processEmail handles email notification delivery.
func (w *NotificationWorker) processEmail(ctx context.Context, msg *models.QueueMessage) error {
	// TODO: Implement actual email sending via SMTP
	// For now, just log and mark as sent
	w.logger.Info("sending email",
		"to", msg.Recipient,
		"subject", msg.Title,
	)

	// Simulate email sending delay
	time.Sleep(100 * time.Millisecond)

	return nil
}

// processSMS handles SMS notification delivery.
func (w *NotificationWorker) processSMS(ctx context.Context, msg *models.QueueMessage) error {
	// TODO: Implement actual SMS sending via Twilio/etc
	w.logger.Info("sending SMS",
		"to", msg.Recipient,
		"body", msg.Body,
	)

	// Simulate SMS sending delay
	time.Sleep(100 * time.Millisecond)

	return nil
}

// processPush handles push notification delivery.
func (w *NotificationWorker) processPush(ctx context.Context, msg *models.QueueMessage) error {
	// TODO: Implement actual push notification via FCM/APNs
	w.logger.Info("sending push notification",
		"token", msg.Recipient[:20]+"...",
		"title", msg.Title,
	)

	// Simulate push sending delay
	time.Sleep(100 * time.Millisecond)

	return nil
}

// updateRecipientStatus updates the status of a notification recipient.
func (w *NotificationWorker) updateRecipientStatus(ctx context.Context, recipientID uuid.UUID, status models.Status, errorMsg string) error {
	var query string
	var args []interface{}

	switch status {
	case models.StatusSent:
		query = `UPDATE notification_recipients SET status = $1, sent_at = NOW() WHERE id = $2`
		args = []interface{}{status, recipientID}
	case models.StatusDelivered:
		query = `UPDATE notification_recipients SET status = $1, delivered_at = NOW() WHERE id = $2`
		args = []interface{}{status, recipientID}
	case models.StatusFailed:
		query = `UPDATE notification_recipients SET status = $1, failed_at = NOW(), error_message = $3 WHERE id = $2`
		args = []interface{}{status, recipientID, errorMsg}
	default:
		query = `UPDATE notification_recipients SET status = $1 WHERE id = $2`
		args = []interface{}{status, recipientID}
	}

	_, err := w.db.Exec(ctx, query, args...)
	return err
}

// incrementRetryCount increments the retry count for a recipient.
func (w *NotificationWorker) incrementRetryCount(ctx context.Context, recipientID uuid.UUID) error {
	query := `UPDATE notification_recipients SET retry_count = retry_count + 1 WHERE id = $1`
	_, err := w.db.Exec(ctx, query, recipientID)
	return err
}

// publishEvent publishes a notification event to RabbitMQ.
func (w *NotificationWorker) publishEvent(ctx context.Context, eventType string, msg *models.QueueMessage, errorMsg string) {
	event := map[string]interface{}{
		"event_type":      eventType,
		"notification_id": msg.NotificationID,
		"recipient_id":    msg.RecipientID,
		"channel":         msg.Channel,
		"timestamp":       time.Now().UTC(),
	}

	if errorMsg != "" {
		event["error"] = errorMsg
	}

	if err := w.rabbitmq.PublishEvent(ctx, eventType, event); err != nil {
		w.logger.Error("failed to publish event",
			"event_type", eventType,
			"error", err,
		)
	}
}

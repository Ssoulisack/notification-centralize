package job

import "github.com/your-org/notification-center/internal/model"

// Queue abstracts the message broker.
// Swap implementations (memory, NATS, RabbitMQ) without changing business logic.
type Queue interface {
	// Publish enqueues a notification for async processing.
	Publish(n *model.Notification) error

	// Subscribe returns a channel that receives notifications.
	Subscribe() <-chan *model.Notification

	// Close shuts down the queue gracefully.
	Close() error
}

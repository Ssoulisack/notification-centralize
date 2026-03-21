package memory

import (
	"github.com/your-org/notification-center/internal/model"
)

// Queue is a simple in-memory queue backed by a Go channel.
// Good for development and testing. Use NATS for production.
type Queue struct {
	ch chan *model.Notification
}

func New(bufferSize int) *Queue {
	return &Queue{
		ch: make(chan *model.Notification, bufferSize),
	}
}

func (q *Queue) Publish(n *model.Notification) error {
	q.ch <- n
	return nil
}

func (q *Queue) Subscribe() <-chan *model.Notification {
	return q.ch
}

func (q *Queue) Close() error {
	close(q.ch)
	return nil
}

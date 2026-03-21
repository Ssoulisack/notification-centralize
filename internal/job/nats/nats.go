package nats

import (
	"encoding/json"
	"fmt"

	"github.com/your-org/notification-center/internal/model"
)

// Queue uses NATS as the message broker for production use.
type Queue struct {
	// conn    *nats.Conn
	subject string
	ch      chan *model.Notification
}

func New(url, subject string) (*Queue, error) {
	// TODO: implement with github.com/nats-io/nats.go
	//
	// conn, err := nats.Connect(url)
	// if err != nil {
	//     return nil, fmt.Errorf("nats connect: %w", err)
	// }

	q := &Queue{
		subject: subject,
		ch:      make(chan *model.Notification, 10000),
	}

	// Subscribe in background
	// conn.Subscribe(subject, func(msg *nats.Msg) {
	//     var n model.Notification
	//     if err := json.Unmarshal(msg.Data, &n); err != nil {
	//         return
	//     }
	//     q.ch <- &n
	// })

	return q, nil
}

func (q *Queue) Publish(n *model.Notification) error {
	data, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("marshal notification: %w", err)
	}

	// TODO: q.conn.Publish(q.subject, data)
	_ = data
	return nil
}

func (q *Queue) Subscribe() <-chan *model.Notification {
	return q.ch
}

func (q *Queue) Close() error {
	// q.conn.Close()
	close(q.ch)
	return nil
}

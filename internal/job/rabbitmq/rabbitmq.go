package rabbitmq

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/model"
)

type Queue struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
	msgs    <-chan amqp.Delivery
	out     chan *model.Notification
	done    chan struct{}
}

func New(cfg config.RabbitMQConfig) (*Queue, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		cfg.Queue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to register consumer: %w", err)
	}

	queue := &Queue{
		conn:    conn,
		channel: ch,
		queue:   q,
		msgs:    msgs,
		out:     make(chan *model.Notification),
		done:    make(chan struct{}),
	}

	go queue.consume()

	return queue, nil
}

func (q *Queue) consume() {
	for {
		select {
		case <-q.done:
			return
		case msg, ok := <-q.msgs:
			if !ok {
				return
			}
			var n model.Notification
			if err := json.Unmarshal(msg.Body, &n); err != nil {
				msg.Nack(false, false)
				continue
			}
			q.out <- &n
			msg.Ack(false)
		}
	}
}

func (q *Queue) Publish(n *model.Notification) error {
	body, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	err = q.channel.Publish(
		"",           // exchange
		q.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (q *Queue) Subscribe() <-chan *model.Notification {
	return q.out
}

func (q *Queue) Close() error {
	close(q.done)
	if err := q.channel.Close(); err != nil {
		return err
	}
	return q.conn.Close()
}

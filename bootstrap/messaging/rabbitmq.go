package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Exchange names
const (
	ExchangeCommands = "notifications.commands"
	ExchangeEvents   = "notifications.events"
	ExchangeDLX      = "notifications.dlx"
)

// Queue names
const (
	QueueInApp = "notifications.send.in_app"
	QueueEmail = "notifications.send.email"
	QueueSMS   = "notifications.send.sms"
	QueuePush  = "notifications.send.push"
	QueueDLQ   = "notifications.dlq"
)

// Routing keys
const (
	RoutingKeyInApp = "in_app"
	RoutingKeyEmail = "email"
	RoutingKeySMS   = "sms"
	RoutingKeyPush  = "push"
)

// Config holds RabbitMQ connection settings.
type Config struct {
	URL            string
	PrefetchCount  int
	ReconnectDelay time.Duration
}

// Client wraps a RabbitMQ connection.
type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  Config
	logger  *slog.Logger
	mu      sync.RWMutex
	closed  bool
}

// NewClient creates a new RabbitMQ client.
func NewClient(ctx context.Context, cfg Config, logger *slog.Logger) (*Client, error) {
	if cfg.PrefetchCount == 0 {
		cfg.PrefetchCount = 10
	}
	if cfg.ReconnectDelay == 0 {
		cfg.ReconnectDelay = 5 * time.Second
	}

	client := &Client{
		config: cfg,
		logger: logger,
	}

	if err := client.connect(); err != nil {
		return nil, err
	}

	if err := client.setupTopology(); err != nil {
		client.Close()
		return nil, err
	}

	return client, nil
}

// connect establishes connection and channel.
func (c *Client) connect() error {
	conn, err := amqp.Dial(c.config.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	if err := ch.Qos(c.config.PrefetchCount, 0, false); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.channel = ch
	c.mu.Unlock()

	return nil
}

// setupTopology declares exchanges and queues.
func (c *Client) setupTopology() error {
	// Declare exchanges
	exchanges := []struct {
		name string
		kind string
	}{
		{ExchangeCommands, "direct"},
		{ExchangeEvents, "topic"},
		{ExchangeDLX, "fanout"},
	}

	for _, ex := range exchanges {
		if err := c.channel.ExchangeDeclare(
			ex.name,
			ex.kind,
			true,  // durable
			false, // auto-delete
			false, // internal
			false, // no-wait
			nil,
		); err != nil {
			return fmt.Errorf("failed to declare exchange %s: %w", ex.name, err)
		}
	}

	// Declare queues with DLX
	queues := []struct {
		name       string
		routingKey string
	}{
		{QueueInApp, RoutingKeyInApp},
		{QueueEmail, RoutingKeyEmail},
		{QueueSMS, RoutingKeySMS},
		{QueuePush, RoutingKeyPush},
	}

	dlxArgs := amqp.Table{
		"x-dead-letter-exchange": ExchangeDLX,
	}

	for _, q := range queues {
		if _, err := c.channel.QueueDeclare(
			q.name,
			true,    // durable
			false,   // auto-delete
			false,   // exclusive
			false,   // no-wait
			dlxArgs, // arguments
		); err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", q.name, err)
		}

		if err := c.channel.QueueBind(
			q.name,
			q.routingKey,
			ExchangeCommands,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed to bind queue %s: %w", q.name, err)
		}
	}

	// Declare DLQ
	if _, err := c.channel.QueueDeclare(
		QueueDLQ,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to declare DLQ: %w", err)
	}

	if err := c.channel.QueueBind(
		QueueDLQ,
		"",
		ExchangeDLX,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind DLQ: %w", err)
	}

	return nil
}

// Publish sends a message to the specified routing key.
func (c *Client) Publish(ctx context.Context, routingKey string, message any) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return fmt.Errorf("client is closed")
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return c.channel.PublishWithContext(
		ctx,
		ExchangeCommands,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	)
}

// PublishEvent sends an event to the events exchange.
func (c *Client) PublishEvent(ctx context.Context, routingKey string, event any) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return fmt.Errorf("client is closed")
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return c.channel.PublishWithContext(
		ctx,
		ExchangeEvents,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	)
}

// Consume returns a channel of deliveries from the specified queue.
func (c *Client) Consume(ctx context.Context, queue string, consumerTag string) (<-chan amqp.Delivery, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, fmt.Errorf("client is closed")
	}

	return c.channel.ConsumeWithContext(
		ctx,
		queue,
		consumerTag,
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
}

// Close closes the connection.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true

	var errs []error

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close channel: %w", err))
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close connection: %w", err))
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

// IsConnected returns whether the client is connected.
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn != nil && !c.conn.IsClosed()
}

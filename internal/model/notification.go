package model

import "time"

// Channel represents a notification delivery channel.
type Channel string

const (
	ChannelEmail    Channel = "email"
	ChannelSMS      Channel = "sms"
	ChannelPush     Channel = "push"
	ChannelSlack    Channel = "slack"
	ChannelTelegram Channel = "telegram"
	ChannelLine     Channel = "line"
)

// Priority controls processing order.
type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityNormal   Priority = "normal"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

// Status tracks the notification lifecycle.
type Status string

const (
	StatusPending    Status = "pending"
	StatusQueued     Status = "queued"
	StatusProcessing Status = "processing"
	StatusSent       Status = "sent"
	StatusFailed     Status = "failed"
	StatusCancelled  Status = "cancelled"
)

// Notification is the core domain object.
type Notification struct {
	ID           string            `json:"id"            db:"id"`
	ProjectID    string            `json:"project_id"    db:"project_id"`
	UserID       string            `json:"user_id"       db:"user_id"`
	Channel      Channel           `json:"channel"       db:"channel"`
	Recipient    string            `json:"recipient"     db:"recipient"`
	Subject      string            `json:"subject"       db:"subject"`
	Body         string            `json:"body"          db:"body"`
	TemplateID   string            `json:"template_id"   db:"template_id"`
	Priority     Priority          `json:"priority"      db:"priority"`
	Status       Status            `json:"status"        db:"status"`
	RetryCount   int               `json:"retry_count"   db:"retry_count"`
	Metadata     map[string]string `json:"metadata"      db:"metadata"`
	ErrorMessage string            `json:"error_message" db:"error_message"`
	CreatedAt    time.Time         `json:"created_at"    db:"created_at"`
	SentAt       *time.Time        `json:"sent_at"       db:"sent_at"`
}

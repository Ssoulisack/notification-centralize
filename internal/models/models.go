package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Channel represents notification delivery channels.
type Channel string

const (
	ChannelInApp Channel = "in_app"
	ChannelEmail Channel = "email"
	ChannelSMS   Channel = "sms"
	ChannelPush  Channel = "push"
)

// Status represents notification delivery status.
type Status string

const (
	StatusPending   Status = "pending"
	StatusSent      Status = "sent"
	StatusDelivered Status = "delivered"
	StatusFailed    Status = "failed"
	StatusRead      Status = "read"
)

// Priority represents notification priority levels.
type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityNormal   Priority = "normal"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

// MemberRole represents project member roles.
type MemberRole string

const (
	RoleOwner  MemberRole = "owner"
	RoleAdmin  MemberRole = "admin"
	RoleMember MemberRole = "member"
	RoleViewer MemberRole = "viewer"
)

// User represents a user synced from Keycloak.
type User struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	KeycloakID    string     `json:"keycloak_id" db:"keycloak_id"`
	Email         string     `json:"email" db:"email"`
	Username      string     `json:"username" db:"username"`
	FirstName     string     `json:"first_name,omitempty" db:"first_name"`
	LastName      string     `json:"last_name,omitempty" db:"last_name"`
	AvatarURL     string     `json:"avatar_url,omitempty" db:"avatar_url"`
	EmailVerified bool       `json:"email_verified" db:"email_verified"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// Project represents a notification project.
type Project struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	Name        string          `json:"name" db:"name"`
	Description string          `json:"description,omitempty" db:"description"`
	Slug        string          `json:"slug" db:"slug"`
	Settings    json.RawMessage `json:"settings,omitempty" db:"settings"`
	IsActive    bool            `json:"is_active" db:"is_active"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

// Role represents a role with permissions.
type Role struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	Name        MemberRole  `json:"name" db:"name"`
	Permissions Permissions `json:"permissions" db:"permissions"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
}

// Permissions defines role capabilities.
type Permissions struct {
	Projects      PermissionSet `json:"projects"`
	Notifications PermissionSet `json:"notifications"`
	APIKeys       PermissionSet `json:"api_keys"`
	Analytics     PermissionSet `json:"analytics"`
	Settings      PermissionSet `json:"settings,omitempty"`
}

// PermissionSet defines individual permissions.
type PermissionSet map[string]bool

// ProjectMember links users to projects with roles.
type ProjectMember struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	ProjectID uuid.UUID  `json:"project_id" db:"project_id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	RoleID    uuid.UUID  `json:"role_id" db:"role_id"`
	JoinedAt  time.Time  `json:"joined_at" db:"joined_at"`
	InvitedBy *uuid.UUID `json:"invited_by,omitempty" db:"invited_by"`

	// Joined fields
	User *User `json:"user,omitempty"`
	Role *Role `json:"role,omitempty"`
}

// APIKey represents an API key for project access.
type APIKey struct {
	ID         uuid.UUID       `json:"id" db:"id"`
	ProjectID  uuid.UUID       `json:"project_id" db:"project_id"`
	UserID     uuid.UUID       `json:"user_id" db:"user_id"`
	Name       string          `json:"name" db:"name"`
	KeyPrefix  string          `json:"key_prefix" db:"key_prefix"`
	KeyHash    string          `json:"-" db:"key_hash"`
	Scopes     json.RawMessage `json:"scopes" db:"scopes"`
	LastUsedAt *time.Time      `json:"last_used_at,omitempty" db:"last_used_at"`
	ExpiresAt  *time.Time      `json:"expires_at,omitempty" db:"expires_at"`
	IsActive   bool            `json:"is_active" db:"is_active"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
	RevokedAt  *time.Time      `json:"revoked_at,omitempty" db:"revoked_at"`
}

// NotificationTemplate represents a reusable notification template.
type NotificationTemplate struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	ProjectID       uuid.UUID       `json:"project_id" db:"project_id"`
	Name            string          `json:"name" db:"name"`
	Slug            string          `json:"slug" db:"slug"`
	Channel         Channel         `json:"channel" db:"channel"`
	SubjectTemplate string          `json:"subject_template,omitempty" db:"subject_template"`
	BodyTemplate    string          `json:"body_template" db:"body_template"`
	Variables       json.RawMessage `json:"variables,omitempty" db:"variables"`
	IsActive        bool            `json:"is_active" db:"is_active"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

// Notification represents a notification message.
type Notification struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	ProjectID   uuid.UUID       `json:"project_id" db:"project_id"`
	TemplateID  *uuid.UUID      `json:"template_id,omitempty" db:"template_id"`
	SenderID    *uuid.UUID      `json:"sender_id,omitempty" db:"sender_id"`
	Title       string          `json:"title" db:"title"`
	Body        string          `json:"body" db:"body"`
	Data        json.RawMessage `json:"data,omitempty" db:"data"`
	Priority    Priority        `json:"priority" db:"priority"`
	ScheduledAt *time.Time      `json:"scheduled_at,omitempty" db:"scheduled_at"`
	ExpiresAt   *time.Time      `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`

	// For API responses
	Recipients []NotificationRecipient `json:"recipients,omitempty"`
}

// NotificationRecipient tracks delivery per user/channel.
type NotificationRecipient struct {
	ID               uuid.UUID       `json:"id" db:"id"`
	NotificationID   uuid.UUID       `json:"notification_id" db:"notification_id"`
	UserID           uuid.UUID       `json:"user_id" db:"user_id"`
	Channel          Channel         `json:"channel" db:"channel"`
	RecipientAddress string          `json:"recipient_address,omitempty" db:"recipient_address"`
	Status           Status          `json:"status" db:"status"`
	SentAt           *time.Time      `json:"sent_at,omitempty" db:"sent_at"`
	DeliveredAt      *time.Time      `json:"delivered_at,omitempty" db:"delivered_at"`
	ReadAt           *time.Time      `json:"read_at,omitempty" db:"read_at"`
	FailedAt         *time.Time      `json:"failed_at,omitempty" db:"failed_at"`
	ErrorMessage     string          `json:"error_message,omitempty" db:"error_message"`
	RetryCount       int             `json:"retry_count" db:"retry_count"`
	Metadata         json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
}

// NotificationEvent represents an audit log entry.
type NotificationEvent struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	NotificationID uuid.UUID       `json:"notification_id" db:"notification_id"`
	RecipientID    *uuid.UUID      `json:"recipient_id,omitempty" db:"recipient_id"`
	EventType      string          `json:"event_type" db:"event_type"`
	EventData      json.RawMessage `json:"event_data,omitempty" db:"event_data"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
}

// DeviceToken represents a push notification device.
type DeviceToken struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	ProjectID  uuid.UUID  `json:"project_id" db:"project_id"`
	Token      string     `json:"token" db:"token"`
	Platform   string     `json:"platform" db:"platform"`
	DeviceName string     `json:"device_name,omitempty" db:"device_name"`
	AppVersion string     `json:"app_version,omitempty" db:"app_version"`
	IsActive   bool       `json:"is_active" db:"is_active"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

// WebhookEndpoint represents a webhook delivery target.
type WebhookEndpoint struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	ProjectID uuid.UUID       `json:"project_id" db:"project_id"`
	URL       string          `json:"url" db:"url"`
	Secret    string          `json:"-" db:"secret"`
	Events    json.RawMessage `json:"events" db:"events"`
	IsActive  bool            `json:"is_active" db:"is_active"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

// WebhookDelivery represents a webhook delivery attempt.
type WebhookDelivery struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	EndpointID     uuid.UUID       `json:"endpoint_id" db:"endpoint_id"`
	EventType      string          `json:"event_type" db:"event_type"`
	Payload        json.RawMessage `json:"payload" db:"payload"`
	ResponseStatus int             `json:"response_status,omitempty" db:"response_status"`
	ResponseBody   string          `json:"response_body,omitempty" db:"response_body"`
	Status         string          `json:"status" db:"status"`
	Attempts       int             `json:"attempts" db:"attempts"`
	NextRetryAt    *time.Time      `json:"next_retry_at,omitempty" db:"next_retry_at"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	CompletedAt    *time.Time      `json:"completed_at,omitempty" db:"completed_at"`
}

// UserNotificationPreference stores per-user notification settings.
type UserNotificationPreference struct {
	ID              uuid.UUID `json:"id" db:"id"`
	UserID          uuid.UUID `json:"user_id" db:"user_id"`
	ProjectID       uuid.UUID `json:"project_id" db:"project_id"`
	Channel         Channel   `json:"channel" db:"channel"`
	Enabled         bool      `json:"enabled" db:"enabled"`
	QuietHoursStart string    `json:"quiet_hours_start,omitempty" db:"quiet_hours_start"`
	QuietHoursEnd   string    `json:"quiet_hours_end,omitempty" db:"quiet_hours_end"`
	Frequency       string    `json:"frequency" db:"frequency"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// QueueMessage represents a message in the notification queue.
type QueueMessage struct {
	NotificationID uuid.UUID `json:"notification_id"`
	RecipientID    uuid.UUID `json:"recipient_id"`
	Channel        Channel   `json:"channel"`
	Recipient      string    `json:"recipient"`
	Title          string    `json:"title"`
	Body           string    `json:"body"`
	Data           any       `json:"data,omitempty"`
	Priority       Priority  `json:"priority"`
	RetryCount     int       `json:"retry_count"`
}

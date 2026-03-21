package model

// SendRequest is the API payload for sending a notification.
// Supports two modes:
//   - Direct: provide channel + recipient
//   - Event-based: provide user_id + event (resolver figures out the rest)
type SendRequest struct {
	// Direct mode
	Channel   Channel `json:"channel"   validate:"omitempty,oneof=email sms push slack telegram line"`
	Recipient string  `json:"recipient" validate:"omitempty"`

	// Event mode
	UserID string `json:"user_id" validate:"omitempty"`
	Event  string `json:"event"   validate:"omitempty"` // e.g. "order_shipped"

	// Content (used in direct mode, or overrides template)
	Subject string `json:"subject"  validate:"omitempty,max=500"`
	Body    string `json:"body"     validate:"omitempty,max=10000"`

	// Template mode (used with event mode)
	TemplateID string            `json:"template_id" validate:"omitempty"`
	Data       map[string]string `json:"data"`     // template variables

	Priority Priority          `json:"priority" validate:"omitempty,oneof=low normal high critical"`
	Metadata map[string]string `json:"metadata"`
}

// SendResponse is returned after a notification is accepted.
type SendResponse struct {
	ID      string `json:"id"`
	Status  Status `json:"status"`
	Message string `json:"message,omitempty"`
}

// DeviceRegisterRequest is the payload for registering a device token.
type DeviceRegisterRequest struct {
	UserID     string   `json:"user_id"     validate:"required"`
	Token      string   `json:"token"       validate:"required"`
	Platform   Platform `json:"platform"    validate:"required,oneof=ios android huawei web"`
	AppVersion string   `json:"app_version" validate:"omitempty"`
}

// PreferenceUpdateRequest is the payload for updating user preferences.
type PreferenceUpdateRequest struct {
	EnabledChannels []Channel `json:"enabled_channels" validate:"omitempty,dive,oneof=email sms push slack telegram line"`
	QuietHoursStart *int      `json:"quiet_start"      validate:"omitempty,min=0,max=23"`
	QuietHoursEnd   *int      `json:"quiet_end"        validate:"omitempty,min=0,max=23"`
	OptedOutEvents  []string  `json:"opted_out_events"`
}

// StatusResponse is returned when checking notification status.
type StatusResponse struct {
	ID           string  `json:"id"`
	Status       Status  `json:"status"`
	Channel      Channel `json:"channel"`
	Recipient    string  `json:"recipient"`
	RetryCount   int     `json:"retry_count"`
	ErrorMessage string  `json:"error_message,omitempty"`
	CreatedAt    string  `json:"created_at"`
	SentAt       string  `json:"sent_at,omitempty"`
}

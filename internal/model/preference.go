package model

import "time"

// UserPreference controls how a user receives notifications.
type UserPreference struct {
	ProjectID       string    `json:"project_id"       db:"project_id"`
	UserID          string    `json:"user_id"          db:"user_id"`
	EnabledChannels []Channel `json:"enabled_channels" db:"enabled_channels"`
	QuietHoursStart int       `json:"quiet_start"      db:"quiet_start"` // 0-23
	QuietHoursEnd   int       `json:"quiet_end"        db:"quiet_end"`   // 0-23
	OptedOutEvents  []string  `json:"opted_out_events" db:"opted_out_events"`
	UpdatedAt       time.Time `json:"updated_at"       db:"updated_at"`
}

// IsChannelEnabled checks if a channel is in the user's enabled list.
func (p *UserPreference) IsChannelEnabled(ch Channel) bool {
	for _, c := range p.EnabledChannels {
		if c == ch {
			return true
		}
	}
	return false
}

// IsEventOptedOut checks if the user opted out of a specific event.
func (p *UserPreference) IsEventOptedOut(event string) bool {
	for _, e := range p.OptedOutEvents {
		if e == event {
			return true
		}
	}
	return false
}

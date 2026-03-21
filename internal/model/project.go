package model

import "time"

// Project represents a tenant in the multi-project system.
type Project struct {
	ID        string    `json:"id"         db:"id"`
	Name      string    `json:"name"       db:"name"`
	APIKey    string    `json:"-"          db:"api_key"` // Never expose in JSON
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

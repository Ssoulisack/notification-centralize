package model

import "time"

// Template is a reusable notification template.
type Template struct {
	ID              string    `json:"id"               db:"id"`
	ProjectID       string    `json:"project_id"       db:"project_id"`
	Name            string    `json:"name"             db:"name"` // e.g. "order_shipped"
	Channel         Channel   `json:"channel"          db:"channel"`
	SubjectTemplate string    `json:"subject_template" db:"subject_template"`
	BodyTemplate    string    `json:"body_template"    db:"body_template"`
	CreatedAt       time.Time `json:"created_at"       db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"       db:"updated_at"`
}

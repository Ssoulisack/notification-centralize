package model

import "time"

// Platform represents a mobile platform.
type Platform string

const (
	PlatformIOS     Platform = "ios"
	PlatformAndroid Platform = "android"
	PlatformHuawei  Platform = "huawei"
	PlatformWeb     Platform = "web"
)

// DeviceToken stores a push notification device registration.
type DeviceToken struct {
	ID         string    `json:"id"          db:"id"`
	ProjectID  string    `json:"project_id"  db:"project_id"`
	UserID     string    `json:"user_id"     db:"user_id"`
	Token      string    `json:"token"       db:"token"`
	Platform   Platform  `json:"platform"    db:"platform"`
	AppVersion string    `json:"app_version" db:"app_version"`
	CreatedAt  time.Time `json:"created_at"  db:"created_at"`
	LastUsedAt time.Time `json:"last_used_at" db:"last_used_at"`
}

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type APIKey struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID  uuid.UUID      `gorm:"type:uuid;not null;index"`
	Project    Project        `gorm:"foreignKey:ProjectID"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null"`
	User       User           `gorm:"foreignKey:UserID"`
	Name       string         `gorm:"not null;size:255"`
	KeyPrefix  string         `gorm:"not null;size:16"`
	KeyHash    string         `gorm:"not null;size:64"`
	Scopes     datatypes.JSON `gorm:"type:jsonb;default:'[\"notifications:write\",\"notifications:read\"]'"`
	LastUsedAt *time.Time
	ExpiresAt  *time.Time
	IsActive   bool       `gorm:"default:true"`
	CreatedAt  time.Time
	RevokedAt  *time.Time
}

func (APIKey) TableName() string { return "api_keys" }

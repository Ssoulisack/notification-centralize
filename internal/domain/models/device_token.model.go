package models

import (
	"time"

	"github.com/google/uuid"
)

type DeviceToken struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	User       User      `gorm:"foreignKey:UserID"`
	ProjectID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Project    Project   `gorm:"foreignKey:ProjectID"`
	Token      string    `gorm:"not null"`
	Platform   string    `gorm:"not null;size:50"`
	DeviceName string    `gorm:"size:255"`
	AppVersion string    `gorm:"size:50"`
	IsActive   bool      `gorm:"default:true"`
	LastUsedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (DeviceToken) TableName() string { return "device_tokens" }

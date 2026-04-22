package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	KeycloakID  string     `gorm:"uniqueIndex;not null;size:255"`
	Email       string     `gorm:"uniqueIndex;not null;size:255"`
	Username    string     `gorm:"uniqueIndex;not null;size:255"`
	FirstName   string     `gorm:"size:255"`
	LastName    string     `gorm:"size:255"`
	IsActive    bool       `gorm:"default:true"`
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (User) TableName() string { return "users" }

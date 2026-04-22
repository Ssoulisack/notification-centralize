package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Project struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OwnerID     uuid.UUID      `gorm:"type:uuid;not null;index"`
	Owner       User           `gorm:"foreignKey:OwnerID"`
	Name        string         `gorm:"not null;size:255"`
	Description string
	Slug        string         `gorm:"uniqueIndex;not null;size:255"`
	Settings    datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	IsActive    bool           `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Project) TableName() string { return "projects" }

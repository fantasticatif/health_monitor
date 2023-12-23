package data

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HitPoint struct {
	gorm.Model
	UUID      string `gorm:"type:varchar(100);primaryKey;<-:create"`
	Name      string `gorm:"not null"`
	ProjectID int
	Project   Project
}

func (proj *HitPoint) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New()
	proj.UUID = uuid.String()
	return nil
}

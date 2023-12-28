package data

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HitPoint struct {
	gorm.Model
	UUID        string `gorm:"type:varchar(100);primaryKey;<-:create"`
	Name        string `gorm:"not null"` // add composite constraint
	Description string `gorm:"type:text"`
	ProjectID   uint
	Project     Project
}

func (h *HitPoint) BeforeCreate(db *gorm.DB) error {
	h.UUID = uuid.New().String()
	return nil
}

func (h *HitPoint) Create(db *gorm.DB) error {
	h.BeforeCreate(db)
	tx := db.Create(h)
	return tx.Error
}

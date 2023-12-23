package data

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name string `gorm:"not null"`
	UUID string `gorm:"type:varchar(100);primaryKey;"`
}

func (proj *Project) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New()
	proj.UUID = uuid.String()
	return nil
}

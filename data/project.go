package data

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name      string `gorm:"not null;unique;"`
	UUID      string `gorm:"type:varchar(100);primaryKey;"`
	AccountID uint   `gorm:"not null;"`
	Account   Account
}

func (proj *Project) BeforeCreate(tx *gorm.DB) error {
	proj.UUID = uuid.NewString()
	return nil
}

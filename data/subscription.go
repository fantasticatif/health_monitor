package data

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	UUID      string `gorm:"type:varchar(100);primaryKey;"`
	PlanID    uint   `gorm:"not null;"`
	Plan      Plan
	AccountID uint `gorm:"not null;"`
	Account   Account
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	s.UUID = uuid.NewString()
	return nil
}

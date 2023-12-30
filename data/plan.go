package data

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlanDuration uint8

const (
	Lifetime PlanDuration = 1
	Monthly  PlanDuration = 2
)

type Plan struct {
	gorm.Model
	Name         string `gorm:"not null;"`
	UUID         string `gorm:"type:varchar(100);primaryKey;"`
	DurationType PlanDuration
}

func (p *Plan) BeforeCreate(tx *gorm.DB) error {
	p.UUID = uuid.NewString()
	return nil
}

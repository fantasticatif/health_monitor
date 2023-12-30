package data

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Name     string
	UUID     string `gorm:"type:varchar(100);primaryKey;"`
	Users    []User
	Projects []Project
}

func (a *Account) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New()
	a.UUID = uuid.String()
	return nil
}

package data

import (
	"errors"

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

func (p *Project) CreateProject(user User, db *gorm.DB) error {
	if user.ID == 0 {
		return errors.New("user id not found")
	}
	if p.UUID == "" {
		p.BeforeCreate(db)
	}
	tx := db.Begin()

	if err := tx.Create(p).Error; err != nil {
		tx.Rollback() // Rollback if an error occurs
		return err
	}

	tx.Commit()
	return nil
}

package data

import (
	"errors"

	"gorm.io/gorm"
)

type ProjectUser struct {
	ProjectId uint `gorm:"not null;"`
	UserId    uint `gorm:"not null;"`
	AccountID uint
	Account   Account
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

	pu := ProjectUser{ProjectId: p.ID, UserId: user.ID}
	if err := tx.Create(&pu).Error; err != nil {
		tx.Rollback() // Rollback if an error occurs
		return err
	}

	// Commit if everything went well
	tx.Commit()
	return nil
}

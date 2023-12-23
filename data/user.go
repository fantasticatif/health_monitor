package data

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `gorm:"type:varchar(255);not null"`
	Email        string `gorm:"type:varchar(500);not null;unique"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
}

func (user *User) SetPassword(tx *gorm.DB, password string) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)
	return nil
}

func (user *User) VerifyPasswordMatch(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}

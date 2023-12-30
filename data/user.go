package data

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID         string `gorm:"type:varchar(100);primaryKey;"`
	Name         string `gorm:"type:varchar(255);not null"`
	Email        string `gorm:"type:varchar(500);not null;unique"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	IsVerified   bool   `gorm:"not null;"`
	AccountID    uint   `gorm:"not null;"`
	Account      Account
}

func (user *User) SetPassword(password string) (err error) {
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

func (user *User) Create(db *gorm.DB, password string) error {
	pass_len := len(password)
	if pass_len < 6 {
		return errors.New("Password must be at least 6 character long")
	}
	user.SetPassword(password)
	user.UUID = uuid.New().String()
	tx := db.Create(&user)

	if tx.Error != nil {
		return tx.Error
	} else {
		return nil
	}
}

func UserById(db *gorm.DB, id int) (*User, error) {
	var user = User{}
	tx := db.First(&user, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

func UserByEmail(db *gorm.DB, email string) (*User, error) {
	var user = User{}
	tx := db.Where("email", email).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

/*
It will get user by email and verify its password matches. if authentication fail, error will be returned.
User will be returned if found even if password dont match.
*/

func AuthenticateUserByEmailPassword(db *gorm.DB, email string, password string) (*User, error) {
	if password == "" {
		return nil, errors.New("password is not provided")
	}
	user, err := UserByEmail(db, email)
	if err != nil {
		return nil, err
	}
	err = user.VerifyPasswordMatch(password)
	return user, err
}

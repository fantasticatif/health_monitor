package data

type AccountUser struct {
	AccountID   uint   `gorm:"not null;"`
	UserID      uint   `gorm:"not null;"`
	AccountUUID string `gorm:"not null;"`
	UserUUID    string `gorm:"not null;"`

	Account Account
	User    User
}

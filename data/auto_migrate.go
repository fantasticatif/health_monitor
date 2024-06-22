package data

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&Account{}, &User{}, &Project{}, &HitPoint{}, &HitEvent{}, &AccountUser{})

}

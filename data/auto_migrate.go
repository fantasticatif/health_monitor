package data

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Project{}, &HitPoint{}, &HitEvent{})
}

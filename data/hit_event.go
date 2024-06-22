package data

import "gorm.io/gorm"

type HitEvent struct {
	gorm.Model
	Body       string    `gorm:"type:text;nullable"`
	StatusRaw  string    `gorm:"type:varchar(25);nullable"`
	Status     HP_STATUS `gorm:"type:varchar(25)"`
	Query      string    `gorm:"type:varchar(200);nullable;"`
	FromIp     string    `gorm:"type:varchar(15)"`
	UserAgent  string    `gorm:"type:varchar(200)"`
	HttpMethod string    `gorm:"type:varchar(10)"`
	HitPointID uint      `gorm:"not null;"`
	HitPoint   HitPoint
}

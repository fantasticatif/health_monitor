package data

import "gorm.io/gorm"

type HitEvent struct {
	gorm.Model
	Body       string `gorm:"type:text"`
	Query      string `gorm:"type:varchar(200)"`
	FromIp     string `gorm:"type:varchar(15)"`
	UserAgent  string `gorm:"type:varchar(200)"`
	HttpMethod string `gorm:"type:varchar(10)"`
}

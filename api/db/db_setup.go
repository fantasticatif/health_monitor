package db

import (
	"log"
	"os"

	"github.com/fantasticatif/health_monitor/data"
	"gorm.io/gorm"
)

var SharedDB *gorm.DB

func SetupDb() {
	dbConfig := data.GormTCPConnectionConfig{
		UserName: os.Getenv("HM_DB_USERNAME"),
		Password: os.Getenv("HM_DB_PASSWORD"),
		DBName:   os.Getenv("HM_DB_NAME"),
		Host:     os.Getenv("HM_DB_HOST"),
	}
	db, err := dbConfig.OpenMySql()
	if err != nil {
		log.Fatal(err)
	}
	data.AutoMigrate(db)
	SharedDB = db
}

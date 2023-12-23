package data

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormTCPConnectionConfig struct {
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	UserName string
	Password string
	DBName   string
	Host     string
	Charset  string
	Port     string
}

func (c *GormTCPConnectionConfig) dsn(defaultPort string) string {
	cs := c.Charset
	if cs == "" {
		cs = "utf8mb4"
	}
	port := c.Port
	if port == "" {
		port = defaultPort
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True",
		c.UserName,
		c.Password,
		c.Host,
		port,
		c.DBName,
		cs)
}

func (c *GormTCPConnectionConfig) OpenMySql() (*gorm.DB, error) {
	return gorm.Open(mysql.Open(c.dsn("3306")), &gorm.Config{})
}

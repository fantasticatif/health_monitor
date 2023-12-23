package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/fantasticatif/health_monitor/data"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var sharedDB *gorm.DB

func ping(ctx *gin.Context) {
	// platform := ClientPlatform(ctx.Param("platform"))
	// if platform == IosApp || platform == AndroidApp || platform == Web {
	// 	ctx.JSON(200, "Hello world:"+platform)
	// } else {
	// 	ctx.JSON(400, "invalid platform:"+platform)
	// }
	target_id := ctx.Param("uuid")

	body := ""
	if strings.ToLower(ctx.Request.Method) == "post" {
		data, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			fmt.Println("error reading body " + err.Error())
		}
		if data != nil {
			body = string(body)
		}
	}
	fmt.Printf("target_id: %s\n", target_id)
	ctx.JSON(200, "Hello world:")
}

func main() {

	godotenv.Load()

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
	sharedDB = db
	data.AutoMigrate(db)

	router := gin.Default()
	v1 := router.Group("/api/v1")
	pingRoute := v1.Group("/ping")
	pingRoute.GET("/:uuid", ping)
	pingRoute.POST("/:uuid", ping)
	pingRoute.GET("/:uuid/:flag", ping)
	pingRoute.POST("/:uuid/:flag", ping)
	router.Run()
}

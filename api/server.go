package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/fantasticatif/health_monitor/api/db"
	"github.com/fantasticatif/health_monitor/api/projectroute"
	"github.com/fantasticatif/health_monitor/api/shareddata"
	"github.com/fantasticatif/health_monitor/api/userroute"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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
	shareddata.ResetEnvVariables()
	db.SetupDb()

	router := gin.Default()
	v1 := router.Group("/api/v1")
	pingRoute := v1.Group("/ping")
	pingRoute.GET("/:uuid", ping)
	pingRoute.POST("/:uuid", ping)
	pingRoute.GET("/:uuid/:flag", ping)
	pingRoute.POST("/:uuid/:flag", ping)

	userroute.SetupUserRoute(router)
	projectroute.SetupProjectRoute(router)
	router.Run()
}

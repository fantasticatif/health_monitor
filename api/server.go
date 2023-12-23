package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func ping(ctx *gin.Context) {
	// platform := ClientPlatform(ctx.Param("platform"))
	// if platform == IosApp || platform == AndroidApp || platform == Web {
	// 	ctx.JSON(200, "Hello world:"+platform)
	// } else {
	// 	ctx.JSON(400, "invalid platform:"+platform)
	// }
	target_id := ctx.Param("target_id")
	fmt.Printf("target_id: %s\n", target_id)
	ctx.JSON(200, "Hello world:")
}

func main() {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	pingRoute := v1.Group("/ping")
	pingRoute.GET("/:target_id", ping)
	pingRoute.GET("/:target_id/:flag", ping)
	router.Run()
}

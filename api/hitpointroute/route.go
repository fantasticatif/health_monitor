package hitpointroute

import (
	"net/http"

	"github.com/fantasticatif/health_monitor/api/db"
	"github.com/fantasticatif/health_monitor/api/middleware"
	"github.com/fantasticatif/health_monitor/api/shareddata"
	"github.com/fantasticatif/health_monitor/data"
	"github.com/gin-gonic/gin"
)

func createHitPoint(c *gin.Context) {
	/* expected post param field
	project_uuid, name, description
	*/
	user := c.MustGet("user").(data.User)

	projUUId := c.PostForm("project_uuid")
	if projUUId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "missing project uuid",
		})
		return
	}

	proj := data.Project{}
	projTx := db.SharedDB.Where("uuid", projUUId).First(&proj)
	if projTx.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	projUser := data.ProjectUser{}
	projUserTx := db.SharedDB.Where("user_id", user.ID).Where("project_id", proj.ID).First(&projUser)
	if projUserTx.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":          "you dont have access to project",
			"internal_error": projUserTx.Error.Error(),
		})
		return
	}

	hp := data.HitPoint{
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
		ProjectID:   projUser.ProjectId,
	}
	err := hp.Create(db.SharedDB)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":          "some error occured",
			"internal_error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"hitpoint": hp,
		"project":  proj,
	})
}

func getHitPointByUUID(c *gin.Context) {

}

func SetupProjectRoute(router *gin.Engine) {
	auth_ur := router.Group(shareddata.Authenticated_api_route, middleware.AuthenticateMiddleware)
	// auth_ur.GET("/v1/hitpoint/all", )
	auth_ur.POST("/v1/hitpoint/create", createHitPoint)
	auth_ur.GET("/v1/hitpoint/uuid/:uuid", getHitPointByUUID)
}

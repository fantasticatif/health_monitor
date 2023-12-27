package projectroute

import (
	"net/http"

	"github.com/fantasticatif/health_monitor/api/db"
	"github.com/fantasticatif/health_monitor/api/middleware"
	"github.com/fantasticatif/health_monitor/api/shareddata"
	"github.com/fantasticatif/health_monitor/data"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func getProjects(c *gin.Context) {
	user := c.MustGet("user").(data.User)
	projects := []data.Project{}
	tx := db.SharedDB.Model(&data.Project{}).Joins(
		"inner join project_users pu on pu.project_id = projects.id and pu.user_id = ?",
		user.ID,
	).Find(&projects)
	if tx.Error != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error": "some error occured",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
	})
}

func createProject(c *gin.Context) {
	user := c.MustGet("user").(data.User)
	projName := c.PostForm("name")
	if projName == "" {
		c.AbortWithStatusJSON(500, gin.H{
			"error": "some error occured",
		})
		return
	}
	p := data.Project{Name: projName}
	err := p.CreateProject(user, db.SharedDB)
	if err != nil {
		if sqlErr := err.(*mysql.MySQLError); sqlErr != nil {
			if sqlErr.Number == 1062 {
				c.AbortWithStatusJSON(400, gin.H{
					"error": "Project name already exist",
				})
				return
			}
		}
	}

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error": "some error occured",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"project": p,
	})
}

func SetupProjectRoute(router *gin.Engine) {
	auth_ur := router.Group(shareddata.Authenticated_api_route, middleware.AuthenticateMiddleware)
	auth_ur.GET("/v1/project/all", getProjects)
	auth_ur.POST("/v1/project/create", createProject)
}

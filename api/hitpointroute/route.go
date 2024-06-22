package hitpointroute

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/fantasticatif/health_monitor/api/db"
	"github.com/fantasticatif/health_monitor/api/middleware"
	"github.com/fantasticatif/health_monitor/api/shareddata"
	"github.com/fantasticatif/health_monitor/data"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		if projTx.Error == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid project id",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":        "some error occured",
				"internal_err": projTx.Error.Error(),
			})
		}
		return
	}

	var accountUser data.AccountUser
	accUserTx := db.SharedDB.Where("account_id", proj.AccountID).Where("user_id", user.ID).First(&accountUser)
	if accUserTx.Error == gorm.ErrRecordNotFound {
		if accUserTx.Error == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "user dont have access to account",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":        "some error occured",
				"internal_err": accUserTx.Error.Error(),
			})
		}
		return
	}

	s_duration, _ := strconv.Atoi(c.PostForm("duration"))
	g_duration, _ := strconv.Atoi(c.PostForm("grace_duration"))

	sched, schedErr := data.CreatePeriodicSchedule(data.Duration{
		DType:    data.DurationType(c.PostForm("duration_type")),
		Duration: uint(s_duration),
	}, data.Duration{
		DType:    data.DurationType(c.PostForm("duration_type")),
		Duration: uint(g_duration),
	})
	if schedErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":        "some error schedule",
			"internal_err": schedErr.Error(),
		})
		return
	}
	sched_json, sched_json_err := sched.Json()
	if sched_json_err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":        "some error schedule",
			"internal_err": sched_json_err.Error(),
		})
		return
	}
	hp := data.HitPoint{
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
		ProjectID:   proj.ID,
		Schedule:    sched_json,
		AccountID:   accountUser.AccountID,
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
	hp_uuid := c.Param("uuid")
	var hp data.HitPoint
	dbtx := db.SharedDB.Where("uuid", hp_uuid).First(&hp)

	if dbtx.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "not found",
		})
		return
	}

	// user := c.MustGet("user").(data.User)

	c.JSON(http.StatusOK, gin.H{
		"hitpoint": hp,
	})
}

func getProjectHitPoints(c *gin.Context) {
	// user := c.MustGet("user").(data.User)
	projUUId := c.Param("project_uuid")

	proj := data.Project{}
	projTx := db.SharedDB.Where("uuid", projUUId).First(&proj)
	if projTx.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	// var pu data.ProjectUser
	// dbtx := db.SharedDB.Where("project_id", proj.ID).Where("user_id", user.ID).First(&pu)
	// if dbtx.Error != nil {
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
	// 		"error": "you dont have access to the resource",
	// 	})
	// 	return
	// }

	var hps []data.HitPoint
	dbtx := db.SharedDB.Where("project_id", proj.ID).Find(&hps)
	if dbtx.Error != nil {
		if dbtx.Error != gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":        "some error occured",
				"internal_err": dbtx.Error.Error(),
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"hitpoints": hps,
	})
}

func pingHitPoint(c *gin.Context) {
	hp_uuid := c.Param("uuid")
	var hp data.HitPoint
	tx := db.SharedDB.Where("uuid", hp_uuid).First(&hp)
	if tx.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":        "some error occured",
			"internal_err": tx.Error.Error(),
		})
		return
	}
	// bodyStream, err := c.Request.Body()
	var body string = ""
	// if err != nil {
	// 	bodyBytes, err := io.ReadAll(bodyStream)
	// 	if err == nil {
	// 		body = string(bodyBytes)
	// 	}
	// }
	q := c.Request.URL.RawQuery
	status_raw := c.Param("status")
	hp_status := data.HP_STATUS_OK
	if status_raw == "error" || status_raw == "failed" {
		hp_status = data.HP_STATUS_DOWN
	}
	qRaw, _ := url.QueryUnescape(q)

	ev := data.HitEvent{Body: body, Query: qRaw,
		FromIp:     c.Request.RemoteAddr,
		UserAgent:  c.Request.UserAgent(),
		HttpMethod: c.Request.Method,
		HitPointID: hp.ID,
		StatusRaw:  status_raw,
		Status:     hp_status,
	}

	createTx := db.SharedDB.Create(&ev)
	if createTx.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":        "some error occured",
			"internal_err": createTx.Error.Error(),
		})
		return
	}
	old_status := hp.Status
	hp.Status = hp_status
	tx = db.SharedDB.Model(&hp).Where("id = ?", hp.ID).Update("status", hp_status).Update("last_hit_event_id", ev.ID).Update("last_hit_event_at", ev.CreatedAt)

	if tx.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":        "some error occured",
			"internal_err": tx.Error.Error(),
		})
		return
	}
	if hp_status != old_status {
		fmt.Printf("status changed from %s to %s\n", old_status, hp_status)
		fmt.Println("handle alerting for added channel")
	} else {
		fmt.Printf("status is unchanged: %s\n", hp_status)
	}
}

func SetupProjectRoute(router *gin.Engine) {
	auth_ur := router.Group(shareddata.Authenticated_api_route, middleware.AuthenticateMiddleware)
	auth_ur.GET("/v1/project/:project_uuid/hitpoint/all", getProjectHitPoints)
	auth_ur.POST("/v1/hitpoint/create", createHitPoint)
	auth_ur.GET("/v1/hitpoint/uuid/:uuid", getHitPointByUUID)
	router.GET("/ping/:uuid", pingHitPoint)
	router.GET("/ping/:uuid/:status", pingHitPoint)
}

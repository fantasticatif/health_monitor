package account

import (
	"net/http"

	"github.com/fantasticatif/health_monitor/api/db"
	"github.com/fantasticatif/health_monitor/api/middleware"
	"github.com/fantasticatif/health_monitor/api/shareddata"
	"github.com/fantasticatif/health_monitor/api/transformer"
	"github.com/fantasticatif/health_monitor/data"
	"github.com/gin-gonic/gin"
)

func createAccount(c *gin.Context) {
	accountName := c.PostForm("account_name")
	userName := c.PostForm("user_name")
	userEmail := c.PostForm("user_email")
	password := c.PostForm("password")

	field_validations := make(map[string]string)

	if accountName == "" {
		field_validations["name"] = "missing account name"
	}
	if password == "" {
		field_validations["password"] = "missing password"
	}

	if userName == "" {
		field_validations["password"] = "missing user name"
	}

	if userEmail == "" {
		field_validations["password"] = "missing user email"
	}
	if len(field_validations) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"validation_errors": field_validations,
		})
	}

	tx := db.SharedDB.Begin()

	acc := data.Account{}
	acc.Name = accountName
	acc.BeforeCreate(db.SharedDB)
	accTx := db.SharedDB.Create(&acc)
	if accTx.Error != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":        "some error occured",
			"internal_err": accTx.Error.Error(),
		})
		return
	}

	user := data.User{Name: userName, Email: userEmail, AccountID: acc.ID}
	err := user.Create(db.SharedDB, password)
	if err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":        "some error occured",
			"internal_err": err.Error(),
		})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"account": transformer.AccountAPIResponse(acc),
		"user":    transformer.UserAPIResponse(user),
	})
}

func account(c *gin.Context) {
	user := c.MustGet("user").(data.User)
	c.JSON(http.StatusOK, gin.H{
		"account": transformer.AccountAPIResponse(user.Account),
		"user":    transformer.UserAPIResponse(user),
	})
}

func SetupUserRoute(router *gin.Engine) {
	router.POST("/api/v1/account/create", createAccount)
	auth_ur := router.Group(shareddata.Authenticated_api_route+"/v1/account", middleware.AuthenticateMiddleware)
	auth_ur.GET("/", account)
}

package userroute

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fantasticatif/health_monitor/api/auth"
	"github.com/fantasticatif/health_monitor/api/db"
	"github.com/fantasticatif/health_monitor/api/middleware"
	"github.com/fantasticatif/health_monitor/api/shareddata"
	"github.com/fantasticatif/health_monitor/data"
	"github.com/gin-gonic/gin"
)

func createUserRequestHandler(ctx *gin.Context) {
	bodyData, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(400, "Missing user creation request body")
		return
	}
	var body = make(map[string]interface{})
	err = json.Unmarshal(bodyData, &body)
	if err != nil {
		ctx.JSON(400, "invalid request body")
		return
	}
	fmt.Println("received create user request")
	fmt.Println(body)
	response := map[string]interface{}{}
	status_code := 200
	name, ok := body["name"].(string)
	if !ok {
		status_code = 400
		response["message"] = "missing data"
	}
	email, ok := body["email"].(string)
	if !ok || len(email) == 0 {
		status_code = 400
		response["message"] = "missing data"
	}
	password, ok := body["password"].(string)
	pass_len := len(password)
	if !ok || pass_len == 0 {
		status_code = 400
		response["message"] = "missing data"
	}

	if status_code != 200 {
		ctx.JSON(400, response)
		return
	}

	user := data.User{Name: name, Email: email}
	user.SetPassword(password)
	tx := db.SharedDB.Create(&user)

	if tx.Error != nil {
		fmt.Println(tx.Error.Error())
		response["error"] = tx.Error.Error()
		response["message"] = "failed to create user"
		ctx.JSON(400, response)
	} else {
		user_data := make(map[string]interface{})
		user_data["id"] = user.ID
		user_data["name"] = user.Name
		user_data["email"] = user.Email
		user_data["created_at"] = user.CreatedAt
		response["user"] = user_data
		response["message"] = "sucess"
		ctx.JSON(200, response)
	}

}

func loginRequestHandler(c *gin.Context) {

	// Retrieve username and password from request body
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Validate credentials against your user database or service
	user, err := data.AuthenticateUserByEmailPassword(db.SharedDB, email, password)
	if err != nil || user == nil {
		// Handle authentication failure, e.g., return 401 Unauthorized
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// Upon successful authentication:
	// - Generate an authentication token (e.g., JWT)
	token, err := auth.GenerateAuthToken(*user)
	if err != nil {
		// Handle token generation error
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	// // - Set the token in a cookie or response header
	c.SetCookie("token", token, 0, "/", "", false, true)
	// c.SetCookie("token", token, ...)

	// - Return a success response with user information (optional)
	c.JSON(200, gin.H{"message": "Login successful", "user": user})
}

func getUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists || user == nil {
		c.JSON(401, gin.H{"error": "user not found"})
	}
	user_data := user.(data.User)

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"ID":    user_data.ID,
			"Name":  user_data.Name,
			"Email": user_data.Email,
		},
	})
}

func logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)

	// Send a success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

func SetupUserRoute(router *gin.Engine) {
	ur := router.Group("/api/v1/user")
	ur.POST("/create", createUserRequestHandler)
	ur.POST("/login", loginRequestHandler)
	ur.GET("/logout", logout)

	auth_ur := router.Group(shareddata.Authenticated_api_route, middleware.AuthenticateMiddleware)
	auth_ur.GET("/v1/user", getUser)
}

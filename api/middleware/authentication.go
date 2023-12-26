package middleware

import (
	"github.com/fantasticatif/health_monitor/api/auth"
	"github.com/fantasticatif/health_monitor/api/db"
	"github.com/fantasticatif/health_monitor/api/util"
	"github.com/fantasticatif/health_monitor/data"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthenticateMiddleware(c *gin.Context) {
	token, err := c.Cookie("token")
	if err == nil {
		decode, err := auth.DecodeJwtToken(token)
		if err == nil && decode != nil && decode.Valid {
			user_id, err := util.IntValue((decode.Claims.(jwt.MapClaims))["user_id"])
			if err == nil {
				user, err := data.UserById(db.SharedDB, user_id)
				if err == nil {
					c.Set("user", *user)
					return
				}
			}
		}
	}
	c.AbortWithStatusJSON(401, gin.H{
		"error": "user is not authenticated",
	})

}

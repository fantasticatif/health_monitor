package auth

import (
	"crypto/rand"
	"errors"
	"os"

	"github.com/fantasticatif/health_monitor/data"
	"github.com/golang-jwt/jwt/v5"
)

func generateRandomKey() []byte {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err) // Handle error appropriately in production
	}

	return key
}

func jwtKey() []byte {
	key := os.Getenv("AUTH_KEY")
	return []byte(key)
}

func GenerateAuthToken(user data.User) (string, error) {
	// key := os.Getenv("AUTH_KEY")
	// key := generateRandomKey()
	key := jwtKey()
	println(key)
	if key == nil {
		return "", errors.New("missing AUTH_KEY enviroment variable")
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
	})

	return t.SignedString(key)
}

func DecodeJwtToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return jwtKey(), nil
	})
}

package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     jwt.TimeFunc().Add(24 * time.Hour).Unix(), // Token expires in 1 day
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//signing method hs256 -> HMAC with SHA-256 hashing algorithm
	return token.SignedString([]byte(os.Getenv("JWT_SECRET"))) // byte slice secret -> better than string
}

package utils

import (
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	cost, err := strconv.ParseInt(os.Getenv("SALT_ROUNDS"), 10, 64)
	//receive salt rounds from env convert it into int base 10 and 64 bit size
	if err != nil {
		cost = 10 //default cost if env variable is not set or invalid
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), int(cost))
	//bcrypt work with byte slices and not with strings
	//were hashing the password using bcrypt algorithm
	// it receives the password as a byte slice and the salt rounds

	return string(bytes), err
}

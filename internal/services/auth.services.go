package services

import (
	"context"

	"github.com/Deepsayan-Das/chatter_GO/internal/db"
)

func CreateUser(username, email, password string) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`

	_, err := db.DB.Exec(context.Background(), query, username, email, password)

	return err
}

func GetUserByEmail(email string) (int, string, error) {
	query := `SELECT id, password_hash FROM users WHERE email = $1`
	var id int
	var passwordHash string

	err := db.DB.QueryRow(context.Background(), query, email).Scan(&id, &passwordHash)

	return id, passwordHash, err
}

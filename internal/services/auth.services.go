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

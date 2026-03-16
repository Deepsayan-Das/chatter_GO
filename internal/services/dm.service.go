package services

import (
	"context"

	"github.com/Deepsayan-Das/chatter_GO/internal/db"
)

func CreateDM(senderID int, receiverID int, content string) (int, error) {

	query := `
	INSERT INTO direct_messages (sender_id, receiver_id, content)
	VALUES ($1,$2,$3)
	RETURNING id
	`

	var id int

	err := db.DB.QueryRow(
		context.Background(),
		query,
		senderID,
		receiverID,
		content,
	).Scan(&id)

	return id, err
}

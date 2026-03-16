package services

import (
	"context"

	"github.com/Deepsayan-Das/chatter_GO/internal/db"
	"github.com/Deepsayan-Das/chatter_GO/internal/models"
)

// CREATE TABLE rooms (
//     id SERIAL PRIMARY KEY,
//     name TEXT NOT NULL,
//     created_by INTEGER REFERENCES users(id),
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );

// CREATE TABLE room_members (
//     user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
//     room_id INTEGER REFERENCES rooms(id) ON DELETE CASCADE,
//     joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     PRIMARY KEY (user_id, room_id)
// );

func CreateRoom(name string, userId int) (int, error) {
	query := `INSERT INTO rooms (name , created_by) VALUES ($1, $2) RETURNING id`
	var roomId int
	err := db.DB.QueryRow(context.Background(), query, name, userId).Scan(&roomId)
	return roomId, err
}
func JoinRoom(roomId int, userId int) error {
	query := `INSERT INTO room_members (user_id, room_id) VALUES ($1, $2)`
	_, err := db.DB.Exec(context.Background(), query, userId, roomId)
	return err
}
func LeaveRoom(userID int, roomID int) error {

	query := `
	DELETE FROM room_members
	WHERE user_id = $1 AND room_id = $2
	`

	_, err := db.DB.Exec(
		context.Background(),
		query,
		userID,
		roomID,
	)

	return err
}
func FindRoomsByName(name string) ([]models.Room, error) {

	query := `
	SELECT id, name
	FROM rooms
	WHERE name ILIKE '%' || $1 || '%'
	`

	rows, err := db.DB.Query(context.Background(), query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room

	for rows.Next() {

		var r models.Room

		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, err
		}

		rooms = append(rooms, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}
func FindRoomsByUserId(userId int) ([]models.Room, error) {

	query := `
	SELECT r.id, r.name
	FROM rooms r
	JOIN room_members rm ON r.id = rm.room_id
	WHERE rm.user_id = $1
	`

	rows, err := db.DB.Query(context.Background(), query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room

	for rows.Next() {

		var r models.Room

		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, err
		}

		rooms = append(rooms, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

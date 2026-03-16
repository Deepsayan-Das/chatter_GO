package models

import "time"

type Message struct {
	ID        int        `json:"id"`
	RoomID    int        `json:"room_id"`
	UserID    int        `json:"user_id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	EditedAt  *time.Time `json:"edited_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

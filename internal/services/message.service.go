package services

import (
	"context"
	"errors"

	"github.com/Deepsayan-Das/chatter_GO/internal/db"
	"github.com/Deepsayan-Das/chatter_GO/internal/models"
)

var (
	ErrUserNotInRoom   = errors.New("user is not a member of this room")
	ErrMessageNotFound = errors.New("message not found")
	ErrUnauthorized    = errors.New("unauthorized: you do not own this message")
	ErrRoomNotFound    = errors.New("room not found")
)

// CreateMessage inserts a new message into the given room on behalf of the user.
// Returns the new message ID, or an error if the user is not in the room.
func CreateMessage(userID int, roomID int, content string) (int, error) {
	member, err := IsUserInRoom(userID, roomID)
	if err != nil {
		return 0, err
	}
	if !member {
		return 0, ErrUserNotInRoom
	}

	query := `
		INSERT INTO messages (room_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var messageID int
	err = db.DB.QueryRow(
		context.Background(),
		query,
		roomID,
		userID,
		content,
	).Scan(&messageID)

	return messageID, err
}

// GetMessagesByRoomID returns up to `limit` non-deleted messages for a room,
// ordered newest-first, starting at `offset`. Returns ErrRoomNotFound if the
// room does not exist.
func GetMessagesByRoomID(roomID int, limit int, offset int) ([]models.Message, error) {
	// Validate that the room exists before querying messages.
	roomExists, err := roomExists(roomID)
	if err != nil {
		return nil, err
	}
	if !roomExists {
		return nil, ErrRoomNotFound
	}

	// Apply safe defaults for pagination parameters.
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT id, room_id, user_id, content, created_at, edited_at, deleted_at
		FROM messages
		WHERE room_id = $1
		  AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.DB.Query(context.Background(), query, roomID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(
			&m.ID,
			&m.RoomID,
			&m.UserID,
			&m.Content,
			&m.CreatedAt,
			&m.EditedAt,
			&m.DeletedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	// Check for any error that occurred during row iteration.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// EditMessage updates the content of a message. Returns ErrMessageNotFound if
// no such message exists, or ErrUnauthorized if userID does not own it.
func EditMessage(userID int, messageID int, content string) error {
	query := `
		UPDATE messages
		SET content  = $1,
		    edited_at = NOW()
		WHERE id      = $2
		  AND user_id  = $3
		  AND deleted_at IS NULL
	`

	result, err := db.DB.Exec(
		context.Background(),
		query,
		content,
		messageID,
		userID,
	)
	if err != nil {
		return err
	}

	return checkAffected(result.RowsAffected(), messageID, userID)
}

// DeleteMessage soft-deletes a message by setting deleted_at. Returns
// ErrMessageNotFound if the message is already deleted or does not exist,
// or ErrUnauthorized if userID does not own it.
func DeleteMessage(userID int, messageID int) error {
	query := `
		UPDATE messages
		SET deleted_at = NOW()
		WHERE id        = $1
		  AND user_id   = $2
		  AND deleted_at IS NULL
	`

	result, err := db.DB.Exec(
		context.Background(),
		query,
		messageID,
		userID,
	)
	if err != nil {
		return err
	}

	return checkAffected(result.RowsAffected(), messageID, userID)
}

// ── helpers ──────────────────────────────────────────────────────────────────

// checkAffected translates a RowsAffected count into the appropriate sentinel
// error. It cannot distinguish "wrong user" from "not found" without an extra
// query, so it checks ownership first.
func checkAffected(rowsAffected int64, messageID int, userID int) error {
	if rowsAffected == 0 {
		// Determine whether the message exists at all so we can return a more
		// precise error to the caller.
		exists, err := messageExists(messageID)
		if err != nil {
			return err
		}
		if !exists {
			return ErrMessageNotFound
		}
		return ErrUnauthorized
	}
	return nil
}

// messageExists returns true when a non-deleted message with the given ID exists.
func messageExists(messageID int) (bool, error) {
	var exists bool
	err := db.DB.QueryRow(
		context.Background(),
		`SELECT EXISTS(SELECT 1 FROM messages WHERE id = $1 AND deleted_at IS NULL)`,
		messageID,
	).Scan(&exists)
	return exists, err
}

// roomExists returns true when a room with the given ID exists.
func roomExists(roomID int) (bool, error) {
	var exists bool
	err := db.DB.QueryRow(
		context.Background(),
		`SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1)`,
		roomID,
	).Scan(&exists)
	return exists, err
}

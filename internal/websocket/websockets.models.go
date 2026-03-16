package websocket

type Type int

const (
	join_room = iota
	leave_room
	send_message
	send_dm
)

type WSMessage struct {
	Type       Type   `json:"type"`
	RoomID     int    `json:"room_id"`
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
}

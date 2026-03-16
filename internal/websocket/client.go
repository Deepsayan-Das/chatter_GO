package websocket

import (
	"encoding/json"
	"time"

	"github.com/Deepsayan-Das/chatter_GO/internal/services"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID     int             // user id
	Conn   *websocket.Conn // websocket connection
	Send   chan []byte     // outgoing messages
	RoomID int             // current room
	Hub    *Hub            // hub reference
}

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

/*
READ LOOP
Handles incoming websocket messages
*/
func (c *Client) ReadPump() {

	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(4096)

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))

	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {

		var msg WSMessage

		err := c.Conn.ReadJSON(&msg)

		if err != nil {
			break
		}

		switch msg.Type {

		case join_room:

			allowed, err := services.IsUserInRoom(c.ID, msg.RoomID)

			if err != nil || !allowed {
				continue
			}

			c.Hub.Manager.JoinRoom(msg.RoomID, c)

		case leave_room:

			c.Hub.Manager.LeaveRoom(msg.RoomID, c)

		case send_message:

			messageID, err := services.CreateMessage(c.ID, msg.RoomID, msg.Content)

			if err != nil {
				continue
			}

			payload := map[string]interface{}{
				"type":       "new_message",
				"message_id": messageID,
				"room_id":    msg.RoomID,
				"user_id":    c.ID,
				"content":    msg.Content,
			}

			data, _ := json.Marshal(payload)

			c.Hub.Manager.Broadcast(msg.RoomID, data)

		case send_dm:

			messageID, err := services.CreateDM(c.ID, msg.ReceiverID, msg.Content)

			if err != nil {
				continue
			}

			payload := map[string]interface{}{
				"type":        "new_dm",
				"message_id":  messageID,
				"sender_id":   c.ID,
				"receiver_id": msg.ReceiverID,
				"content":     msg.Content,
			}

			data, _ := json.Marshal(payload)

			c.Hub.SendToUser(msg.ReceiverID, data)
		}
	}
}

/*
WRITE LOOP
Sends outgoing messages to websocket clients
*/
func (c *Client) WritePump() {

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {

		select {

		case message, ok := <-c.Send:

			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.Conn.WriteMessage(websocket.TextMessage, message)

			if err != nil {
				return
			}

		case <-ticker.C:

			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

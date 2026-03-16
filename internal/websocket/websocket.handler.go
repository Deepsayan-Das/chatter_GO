package websocket

import (
	"net/http"

	"github.com/Deepsayan-Das/chatter_GO/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}



func ServeWS(hub *Hub, c *gin.Context) {

	token := c.Query("token")

	userID, err := utils.ValidateToken(token)

	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &Client{
		ID:   userID,
		Conn: conn,
		Send: make(chan []byte),
		Hub:  hub,
	}

	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

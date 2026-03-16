package handlers

import (
	"fmt"
	"strconv"

	"github.com/Deepsayan-Das/chatter_GO/internal/services"
	"github.com/gin-gonic/gin"
)

type sendMessageReq struct {
	RoomID  int    `json:"room_id"`
	Content string `json:"content"`
}

type EditMessageReq struct {
	Content string `json:"content"`
}

func SendMessage(ctx *gin.Context) {
	var req sendMessageReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request "})
		return
	}
	userID := ctx.GetInt("user_id")

	id, err := services.CreateMessage(userID, req.RoomID, req.Content)

	if err != nil {
		fmt.Printf("CreateMessage Error: %v\n", err)

		ctx.JSON(500, gin.H{"error": "Error sending message"})
		return
	}
	ctx.JSON(200, gin.H{"message": "Message sent successfully", "message_id": id})

}
func GetMessages(c *gin.Context) {

	roomID := c.Param("id")

	roomIDInt, err := strconv.Atoi(roomID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid room id"})
		return
	}

	var limit int = 50
	var offset int = 0

	messages, err := services.GetMessagesByRoomID(roomIDInt, limit, offset)

	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch messages"})
		return
	}

	c.JSON(200, messages)
}

func EditMessage(c *gin.Context) {

	userID := c.GetInt("user_id")

	messageID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{"error": "invalid message id"})
		return
	}

	var req EditMessageReq

	c.ShouldBindJSON(&req)

	err = services.EditMessage(userID, messageID, req.Content)

	if err != nil {
		c.JSON(500, gin.H{"error": "edit failed"})
		return
	}

	c.JSON(200, gin.H{"message": "edited"})
}

func DeleteMessage(c *gin.Context) {

	userID := c.GetInt("user_id")

	messageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid message id"})
		return
	}

	err = services.DeleteMessage(userID, messageID)

	if err != nil {
		c.JSON(500, gin.H{"error": "delete failed"})
		return
	}

	c.JSON(200, gin.H{"message": "deleted"})
}

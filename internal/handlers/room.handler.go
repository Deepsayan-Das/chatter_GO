package handlers

import (
	"fmt"

	"github.com/Deepsayan-Das/chatter_GO/internal/services"
	"github.com/gin-gonic/gin"
)

type createRoomReq struct {
	Name string `json:"name"`
}

type joinRoomReq struct {
	RoomID int `json:"room_id"`
}

type LeaveRoomRequest struct {
	RoomID int `json:"room_id"`
}

func CreateRoom(ctx *gin.Context) {
	var req createRoomReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	userID := ctx.GetInt("user_id")

	roomId, err := services.CreateRoom(req.Name, userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error creating room"})
		return
	}
	ctx.JSON(200, gin.H{"message": "Room created successfully", "room_id": roomId})

}

func JoinRoom(ctx *gin.Context) {
	var req joinRoomReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	userID := ctx.GetInt("user_id")

	err = services.JoinRoom(req.RoomID, userID)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(500, gin.H{"error": "Error joining room"})
		return
	}
	ctx.JSON(200, gin.H{"message": "Joined room successfully"})
}
func LeaveRoom(ctx *gin.Context) {

	var req LeaveRoomRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	userID := ctx.GetInt("user_id")

	err := services.LeaveRoom(userID, req.RoomID)

	if err != nil {
		ctx.JSON(500, gin.H{"error": "failed to leave room"})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "left room",
	})
}

func SearchRooms(c *gin.Context) {

	name := c.Query("name")

	rooms, err := services.FindRoomsByName(name)

	if err != nil {
		c.JSON(500, gin.H{"error": "failed to search rooms"})
		return
	}

	c.JSON(200, rooms)
}

func MyRooms(c *gin.Context) {

	userID := c.GetInt("user_id")

	rooms, err := services.FindRoomsByUserId(userID)

	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch rooms"})
		return
	}

	c.JSON(200, rooms)
}

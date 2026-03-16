package handlers

import (
	"fmt"
	"net/http"

	"github.com/Deepsayan-Das/chatter_GO/internal/services"
	"github.com/Deepsayan-Das/chatter_GO/internal/utils"
	"github.com/gin-gonic/gin"
)

type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(ctx *gin.Context) {
	var req CreateUserReq
	err := ctx.ShouldBindJSON(&req) // incoming json -> our created struct
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	hashed, err := utils.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error hashing password"})
		return
	}

	err = services.CreateUser(req.Username, req.Email, hashed)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error creating user"})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"message": "User registered successfully"})
}

func Login(ctx *gin.Context) {
	var req LoginReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	id, hashedPwd, err := services.GetUserByEmail(req.Email)
	if err != nil {
		ctx.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}
	if utils.ComparePassword(hashedPwd, req.Password) {
		token, err := utils.GenerateJWT(id)
		if err != nil {
			fmt.Println("Error generating token: ", err.Error())
			ctx.JSON(500, gin.H{"error": "Error generating token"})
			return
		}
		ctx.JSON(200, gin.H{"token": token})
	} else {
		ctx.JSON(401, gin.H{"error": "Invalid email or password"})
	}

}

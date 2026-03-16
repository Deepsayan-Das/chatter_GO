package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Deepsayan-Das/chatter_GO/internal/db"
	"github.com/Deepsayan-Das/chatter_GO/internal/handlers"
	"github.com/Deepsayan-Das/chatter_GO/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting Chatter GO server...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	err = db.ConnectDB()
	if err != nil {
		log.Fatal("Error connecting to the database: ", err.Error())
	}
	fmt.Println("DB connection established")
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Chat backend running"})
	})

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	authRoutes := r.Group("/")
	authRoutes.Use(middleware.AuthMiddleware())

	authRoutes.POST("/rooms/create", handlers.CreateRoom)
	authRoutes.POST("/rooms/join", handlers.JoinRoom)
	authRoutes.POST("/rooms/leave", handlers.LeaveRoom)
	authRoutes.GET("/rooms/search", handlers.SearchRooms)
	authRoutes.GET("/rooms/my", handlers.MyRooms)

	err = r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal("Error starting the server: ", err.Error())
	}
}

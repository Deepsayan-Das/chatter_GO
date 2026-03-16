package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Deepsayan-Das/chatter_GO/internal/db"
	"github.com/Deepsayan-Das/chatter_GO/internal/handlers"
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
	r := gin.Default() //central router for our server

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Welcome to Chatter GO",
		})
	})
	r.POST("/register", handlers.Register)

	err = r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal("Error starting the server: ", err.Error())
	}
}

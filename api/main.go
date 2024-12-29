package main

import (
	"fmt"
	"log"
	"os"
	"whiteboarder/handler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	websocketURL := os.Getenv("WEBSOCKET_URL")
	fmt.Println(websocketURL)
	if websocketURL == "" {
		websocketURL = "ws://localhost:8080/ws"
	}

	r := gin.Default()

	r.Static("/static", "./public")
	r.LoadHTMLFiles("templates/index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"websocketURL": websocketURL,
		})
	})

	r.GET("/ws", handler.HandleConnection)

	go handler.HandleBroadcast()

	fmt.Println("Server is listening on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}

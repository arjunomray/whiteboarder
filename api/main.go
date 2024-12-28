package main

import (
	"fmt"
	"log"
	"whiteboarder/handler"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.Static("/static", "./public")
	r.LoadHTMLFiles("templates/index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/ws", handler.HandleConnection)

	go handler.HandleBroadcast()

	fmt.Println("Server is listening on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}

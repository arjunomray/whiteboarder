package handler

import (
	"fmt"
	"log"
	"net/http"
	"whiteboarder/model"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {

		return true
	},
}

var clients = make(map[*model.Client]bool)
var broadcast = make(chan []byte)

func HandleConnection(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	client := &model.Client{Conn: conn, Send: make(chan []byte)}
	clients[client] = true

	go handleClientMessages(client)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			delete(clients, client)
			return
		}

		if len(msg) > 0 {
			broadcast <- msg
		}
	}
}

func handleClientMessages(client *model.Client) {
	for msg := range client.Send {

		err := client.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Error writing message:", err)
			delete(clients, client)
			client.Conn.Close()
			broadcastClientCount()
			return
		}
	}
}

func HandleBroadcast() {
	for {
		msg := <-broadcast

		for c := range clients {
			c.Send <- msg
		}
	}
}

func broadcastClientCount() {
	clientCount := fmt.Sprintf(`{"type": "client_count", "count": %d}`, len(clients))
	for c := range clients {
		c.Send <- []byte(clientCount)
	}
}

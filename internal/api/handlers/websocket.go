package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by default.
		// In a production environment, you might want to restrict this.
		return true
	},
}

// WebSocketHandler handles the WebSocket connection.
func WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket connection established.")

	// This is a simple loop to demonstrate sending messages.
	// In a real application, you would have a more sophisticated message handling system.
	for {
		// Send a test message
		message := map[string]string{"message": "This is a test message from the WebSocket server."}
		if err := conn.WriteJSON(message); err != nil {
			log.Printf("Error sending message: %v", err)
			break
		}

		// In a real application, you'd likely have a select statement here
		// to handle incoming messages, tickers for periodic messages, and a done channel.
		// For this task, we'll just send one message and then read to keep the connection open.

		// Wait for a message from the client (or for the connection to close)
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket connection closed: %v", err)
			break
		}
	}
}
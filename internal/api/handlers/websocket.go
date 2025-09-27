package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).Msg("Failed to set websocket upgrade")
		return
	}
	defer conn.Close()

	log.Info().Str("remote_addr", conn.RemoteAddr().String()).Msg("WebSocket connection established")

	// This is a simple loop to demonstrate sending messages.
	// In a real application, you would have a more sophisticated message handling system.
	for {
		// Send a test message
		message := map[string]string{"message": "This is a test message from the WebSocket server."}
		if err := conn.WriteJSON(message); err != nil {
			log.Error().Err(err).Msg("Error sending message over WebSocket")
			break
		}

		// In a real application, you'd likely have a select statement here
		// to handle incoming messages, tickers for periodic messages, and a done channel.
		// For this task, we'll just send one message and then read to keep the connection open.

		// Wait for a message from the client (or for the connection to close)
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().Err(err).Msg("Unexpected WebSocket close error")
			} else {
				log.Info().Str("remote_addr", conn.RemoteAddr().String()).Msg("WebSocket connection closed")
			}
			break
		}
	}
}
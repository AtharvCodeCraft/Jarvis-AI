package main

import (
	"log"
	"net/http"

	"jarvis-ai/modules/commands"
	"jarvis-ai/modules/system"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// API Endpoints
	r.POST("/command", handleCommand)
	r.GET("/system/status", handleSystemStatus)

	// WebSocket Endpoint
	r.GET("/ws", handleWebSocket)

	return r
}

func handleCommand(c *gin.Context) {
	var req struct {
		Text string `json:"text"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result := commands.ParseAndExecute(req.Text)
	c.JSON(http.StatusOK, gin.H{"result": result})
}

func handleSystemStatus(c *gin.Context) {
	status := system.GetStatus()
	c.JSON(http.StatusOK, status)
}

func handleWebSocket(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}
	defer ws.Close()

	log.Println("Client connected to WS")

	for {
		var msg map[string]interface{}
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("WS Read Error:", err)
			break
		}

		// Handle incoming WS event (e.g., voice_command, manual_command)
		eventType, _ := msg["type"].(string)

		if eventType == "voice_command" || eventType == "manual_command" {
			text, _ := msg["text"].(string)
			// Echo back the language the client sent so TTS uses the right voice
			language, _ := msg["language"].(string)
			if language == "" {
				language = "en-IN"
			}
			log.Printf("Received command via WS [%s]: %s", language, text)

			// Process command
			response := commands.ParseAndExecute(text)

			// Send response back with language so frontend TTS speaks in correct language
			ws.WriteJSON(map[string]interface{}{
				"type":     "command_result",
				"text":     response,
				"language": language,
			})
		}
	}
}

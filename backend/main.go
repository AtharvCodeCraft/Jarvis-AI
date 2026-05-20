package main

import (
	"log"

	"jarvis-ai/modules/database"
	"jarvis-ai/modules/plugins"
)

func main() {
	log.Println("Initializing JARVIS Core Engine...")

	// Initialize SQLite Database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Load Plugins
	plugins.LoadAll()

	// Setup Server and Router
	r := SetupRouter()

	log.Println("JARVIS system online. Listening on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

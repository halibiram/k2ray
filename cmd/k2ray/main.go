package main

import (
	"github.com/gin-gonic/gin"
	"k2ray/internal/api"
	"k2ray/internal/config"
	"k2ray/internal/db"
	"log"
)

func main() {
	// Load application configuration
	config.LoadConfig("") // Load from default path "configs/system.env"
	log.Println("Configuration loaded successfully.")

	// Initialize database connection
	db.InitDB()


	router := gin.Default()

	// Setup routes from the internal/api package
	api.SetupRouter(router, true) // Enable rate limiting in production

	// In a real application, the port should be configurable.
	// For now, we'll hardcode it to 8080.
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

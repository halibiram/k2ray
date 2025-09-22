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

	// Run database migrations
	db.RunMigrations()

	// Seed database with initial data if necessary
	db.SeedDatabase()

	router := gin.Default()

	// Setup routes from the internal/api package
	api.SetupRouter(router)

	// In a real application, the port should be configurable.
	// For now, we'll hardcode it to 8080.
	if err := router.Run(":8080"); err != nil {
		// A structured logger should be used in a real application.
		panic("failed to start server: " + err.Error())
	}
}

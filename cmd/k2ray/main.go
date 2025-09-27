package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"k2ray/internal/api"
	"k2ray/internal/config"
	"k2ray/internal/db"
	"k2ray/internal/logger"
)

func main() {
	// Initialize the structured logger as the first step.
	logger.InitLogger()

	// Load application configuration
	config.LoadConfig("") // Load from default path "configs/system.env"
	log.Info().Msg("Configuration loaded successfully.")

	// Initialize database connection
	db.InitDB()

	router := gin.Default()

	// Setup routes from the internal/api package
	api.SetupRouter(router, true) // Enable rate limiting in production

	// In a real application, the port should be configurable.
	// For now, we'll hardcode it to 8080.
	log.Info().Msg("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
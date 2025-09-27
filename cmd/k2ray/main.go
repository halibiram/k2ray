package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"k2ray/internal/api"
	"k2ray/internal/config"
	"k2ray/internal/db"
	"k2ray/internal/logger"
	"k2ray/internal/redis"
)

// @title K2Ray API
// @version 1.0
// @description This is the API for K2Ray, a modern V2Ray management panel.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Initialize the structured logger as the first step.
	logger.InitLogger()

	// Load application configuration
	config.LoadConfig("") // Load from default path "configs/system.env"
	log.Info().Msg("Configuration loaded successfully.")

	// Initialize database connection
	db.InitDB()

	// Initialize Redis connection
	redis.InitRedis()

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
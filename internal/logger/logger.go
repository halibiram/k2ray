package logger

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger initializes the global zerolog logger.
// It configures a human-friendly console logger for development
// and a structured JSON logger for production.
func InitLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Use console writer for development for better readability
	if gin.IsDebugging() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		// Use JSON logger in production
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	log.Info().Msg("Logger initialized")
}
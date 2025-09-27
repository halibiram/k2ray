package logger

import (
	"io"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger initializes the global zerolog logger.
// It supports structured JSON logging, file rotation, and remote TCP logging.
// Configuration is managed via environment variables:
// - LOG_LEVEL: (trace, debug, info, warn, error, fatal, panic), defaults to "info"
// - LOG_PATH: Path to the log file. If empty, logs to stderr.
// - REMOTE_LOG_URL: TCP address for remote logging (e.g., "localhost:4000").
func InitLogger() {
	// 1. Set Log Level from Environment Variable
	logLevelStr := os.Getenv("LOG_LEVEL")
	level, err := zerolog.ParseLevel(logLevelStr)
	if err != nil || logLevelStr == "" {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// 2. Create a list of writers
	var writers []io.Writer

	// Always add a JSON writer to stderr. In production, this is the primary console log.
	// In development, it provides a machine-readable log alongside the human-friendly one.
	writers = append(writers, zerolog.New(os.Stderr).With().Timestamp().Logger())

	// Add a human-friendly console writer only when in debug mode
	if gin.IsDebugging() {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	// 3. Configure File-based Logging with Rotation
	logPath := os.Getenv("LOG_PATH")
	if logPath != "" {
		fileLogger := &lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    10, // megabytes
			MaxBackups: 5,
			MaxAge:     30, // days
			Compress:   true,
		}
		writers = append(writers, fileLogger)
		log.Info().Str("log_path", logPath).Msg("File logging enabled")
	}

	// 4. Configure Remote Logging
	remoteLogURL := os.Getenv("REMOTE_LOG_URL")
	if remoteLogURL != "" {
		conn, err := net.Dial("tcp", remoteLogURL)
		if err != nil {
			log.Error().Err(err).Msg("Failed to connect to remote log server")
		} else {
			writers = append(writers, conn)
			log.Info().Str("remote_url", remoteLogURL).Msg("Remote logging enabled")
		}
	}

	// 5. Create a multi-writer to combine all configured outputs
	multiWriter := io.MultiWriter(writers...)

	// 6. Build the final logger instance
	log.Logger = zerolog.New(multiWriter).
		With().
		Timestamp().
		Str("service", "k2ray"). // Add a static field for the service name
		Logger()

	// Add a hook to include caller information (file and line number)
	// This is useful for debugging but can add overhead.
	if zerolog.GlobalLevel() <= zerolog.DebugLevel {
		log.Logger = log.Logger.With().Caller().Logger()
	}

	log.Info().
		Str("log_level", level.String()).
		Bool("debug_mode", gin.IsDebugging()).
		Msg("Logger initialized successfully")
}

// GinLogger is a middleware that logs HTTP requests using zerolog.
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Don't log favicon requests for cleaner logs
		if path == "/favicon.ico" {
			return
		}

		end := time.Now()
		latency := end.Sub(start)

		// Get additional details
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		event := log.Info()
		if statusCode >= 500 {
			event = log.Error()
		} else if statusCode >= 400 {
			event = log.Warn()
		}

		event.
			Str("method", method).
			Str("path", path).
			Int("status_code", statusCode).
			Str("client_ip", clientIP).
			Dur("latency", latency).
			Int("body_size", bodySize).
			Str("error", errorMessage).
			Msg("Request completed")
	}
}

// Helper function to convert string to int, with a fallback.
func getEnvAsInt(name string, fallback int) int {
	if value, exists := os.LookupEnv(name); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}
package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// Config holds the application configuration.
// Using a struct provides type safety and a single source of truth for config values.
type Config struct {
	DatabaseURL string
	JWTSecret   string
	AppName     string
}

// AppConfig is a singleton instance of the Config struct.
var (
	AppConfig *Config
	once      sync.Once
)

// ResetForTesting is a function to reset the singleton for testing purposes.
func ResetForTesting() {
	once = sync.Once{}
	AppConfig = nil
}

// LoadConfig loads configuration from a .env file and the environment.
// It is safe to call from multiple goroutines.
func LoadConfig(path string) {
	once.Do(func() {
		// If no path is provided, default to "configs/system.env"
		if path == "" {
			path = "configs/system.env"
		}

		// godotenv.Load will not override existing environment variables.
		err := godotenv.Load(path)
		if err != nil {
			log.Warn().Err(err).Str("path", path).Msg("Could not load .env file, falling back to environment variables")
		}

		AppConfig = &Config{
			DatabaseURL: getEnv("DATABASE_URL", "./k2ray.db"),
			JWTSecret:   getEnv("JWT_SECRET", "default-secret-please-change"),
			AppName:     getEnv("APP_NAME", "k2ray"),
		}
	})
}

// getEnv retrieves an environment variable by key, returning a fallback if not found.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

// Config holds the application configuration.
// Using a struct provides type safety and a single source of truth for config values.
type Config struct {
	DatabaseURL string
	JWTSecret   string
}

// AppConfig is a singleton instance of the Config struct.
var (
	AppConfig *Config
	once      sync.Once
)

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
			log.Printf("Warning: could not load %s file: %v. Falling back to environment variables.", path, err)
		}

		AppConfig = &Config{
			DatabaseURL: getEnv("DATABASE_URL", "./k2ray.db"),
			JWTSecret:   getEnv("JWT_SECRET", "default-secret-please-change"),
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

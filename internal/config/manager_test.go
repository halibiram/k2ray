package config_test

import (
	"k2ray/internal/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Reset singleton for test isolation
	config.ResetForTesting()

	// Create a temporary .env file for testing
	content := []byte("DATABASE_URL=test.db\nJWT_SECRET=test-secret")
	tmpfile, err := os.CreateTemp("", "*.env")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name()) // clean up

	_, err = tmpfile.Write(content)
	assert.NoError(t, err)
	err = tmpfile.Close()
	assert.NoError(t, err)

	// Load the config from the temporary file
	config.LoadConfig(tmpfile.Name())

	// Assert that the config values are loaded correctly
	assert.NotNil(t, config.AppConfig)
	assert.Equal(t, "test.db", config.AppConfig.DatabaseURL)
	assert.Equal(t, "test-secret", config.AppConfig.JWTSecret)
}

func TestLoadConfig_Fallback(t *testing.T) {
	// Reset singleton for test isolation
	config.ResetForTesting()

	// Ensure env vars are not set
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("JWT_SECRET")

	// Load config from a non-existent path to trigger fallback
	config.LoadConfig("non-existent-file.env")

	// Assert that the fallback values are used
	assert.NotNil(t, config.AppConfig)
	assert.Equal(t, "./k2ray.db", config.AppConfig.DatabaseURL)
	assert.Equal(t, "default-secret-please-change", config.AppConfig.JWTSecret)
}

package db_test

import (
	"k2ray/internal/config"
	"k2ray/internal/db"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestMain sets up an in-memory SQLite database for the tests in this package.
func TestMain(m *testing.M) {
	// Use an in-memory database for testing
	config.AppConfig = &config.Config{
		DatabaseURL: ":memory:",
	}

	db.InitDB()

	// Run the tests
	code := m.Run()

	db.DB.Close()
	os.Exit(code)
}

func TestTokenBlocklisting(t *testing.T) {
	// Clean up table before test
	_, err := db.DB.Exec("DELETE FROM revoked_tokens")
	assert.NoError(t, err)

	jti1 := "test-jti-1"
	jti2 := "test-jti-2"
	expiresAt := time.Now().Add(1 * time.Hour)

	// 1. Initially, token should not be blocklisted
	isBlocklisted, err := db.IsTokenBlocklisted(jti1)
	assert.NoError(t, err)
	assert.False(t, isBlocklisted)

	// 2. Blocklist a token
	err = db.BlocklistToken(jti1, expiresAt)
	assert.NoError(t, err)

	// 3. Verify it is now blocklisted
	isBlocklisted, err = db.IsTokenBlocklisted(jti1)
	assert.NoError(t, err)
	assert.True(t, isBlocklisted)

	// 4. Verify another token is still not blocklisted
	isBlocklisted, err = db.IsTokenBlocklisted(jti2)
	assert.NoError(t, err)
	assert.False(t, isBlocklisted)
}

// Note: A test for a cleanup function for expired tokens would go here
// if such a function were implemented.

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

func TestCleanupExpiredTokens(t *testing.T) {
	// Clean up table before test
	_, err := db.DB.Exec("DELETE FROM revoked_tokens")
	assert.NoError(t, err)

	// Add some tokens: 1 expired, 1 not
	expiredJTI := "expired-jti"
	validJTI := "valid-jti"
	expiredTime := time.Now().Add(-1 * time.Hour) // Expired 1 hour ago
	validTime := time.Now().Add(1 * time.Hour)   // Expires in 1 hour

	err = db.BlocklistToken(expiredJTI, expiredTime)
	assert.NoError(t, err)
	err = db.BlocklistToken(validJTI, validTime)
	assert.NoError(t, err)

	// Verify both tokens are initially in the blocklist
	isBlocklisted, err := db.IsTokenBlocklisted(expiredJTI)
	assert.NoError(t, err)
	assert.True(t, isBlocklisted, "Expired token should be blocklisted initially")

	isBlocklisted, err = db.IsTokenBlocklisted(validJTI)
	assert.NoError(t, err)
	assert.True(t, isBlocklisted, "Valid token should be blocklisted initially")

	// Run the cleanup function
	rowsAffected, err := db.CleanupExpiredTokens()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected, "Should have cleaned up exactly one expired token")

	// Verify the expired token is gone
	isBlocklisted, err = db.IsTokenBlocklisted(expiredJTI)
	assert.NoError(t, err)
	assert.False(t, isBlocklisted, "Expired token should have been removed")

	// Verify the valid token remains
	isBlocklisted, err = db.IsTokenBlocklisted(validJTI)
	assert.NoError(t, err)
	assert.True(t, isBlocklisted, "Valid token should not have been removed")
}

package auth_test

import (
	"k2ray/internal/auth"
	"k2ray/internal/config"
	"k2ray/internal/db"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setup() {
	// Mock the config for testing
	config.AppConfig = &config.Config{
		JWTSecret: "test-jwt-secret-for-auth-package",
	}
}

func TestTokenGenerationAndValidation(t *testing.T) {
	setup()

	user := db.User{
		ID:       123,
		Username: "testuser",
		Role:     db.RoleUser,
	}

	// 1. Generate tokens
	accessToken, refreshToken, err := auth.GenerateTokens(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	// 2. Validate Access Token
	accessClaims, err := auth.ValidateToken(accessToken)
	assert.NoError(t, err)
	assert.NotNil(t, accessClaims)
	assert.Equal(t, user.ID, accessClaims.UserID)
	assert.Equal(t, user.Username, accessClaims.Username)
	assert.NotEmpty(t, accessClaims.ID) // JTI should exist
	assert.Equal(t, "k2ray", accessClaims.Issuer)
	// Check that expiration is roughly 15 minutes from now
	assert.WithinDuration(t, time.Now().Add(15*time.Minute), accessClaims.ExpiresAt.Time, 5*time.Second)

	// 3. Validate Refresh Token
	refreshClaims, err := auth.ValidateToken(refreshToken)
	assert.NoError(t, err)
	assert.NotNil(t, refreshClaims)
	assert.Equal(t, user.ID, refreshClaims.UserID)
	// Check that expiration is roughly 7 days from now
	assert.WithinDuration(t, time.Now().Add(7*24*time.Hour), refreshClaims.ExpiresAt.Time, 5*time.Second)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	setup()

	// Test with a malformed token string
	_, err := auth.ValidateToken("this.is.not.a.valid.token")
	assert.Error(t, err)

	// Test with a token signed with a different key
	otherSecretConfig := &config.Config{JWTSecret: "other-secret"}
	config.AppConfig = otherSecretConfig
	accessToken, _, _ := auth.GenerateTokens(db.User{ID: 1, Username: "user", Role: db.RoleUser})

	// Switch back to the original secret to validate
	setup()
	_, err = auth.ValidateToken(accessToken)
	assert.Error(t, err, "Should fail validation with a different secret")
}

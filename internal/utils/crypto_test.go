package utils_test

import (
	"k2ray/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHashing(t *testing.T) {
	password := "my-secret-password"

	// Test hashing
	hash, err := utils.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	// Test correct password check
	isCorrect := utils.CheckPasswordHash(password, hash)
	assert.True(t, isCorrect, "Should be true for correct password")

	// Test incorrect password check
	isIncorrect := utils.CheckPasswordHash("wrong-password", hash)
	assert.False(t, isIncorrect, "Should be false for incorrect password")
}

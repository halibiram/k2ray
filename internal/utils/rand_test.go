package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecureIntn(t *testing.T) {
	// Test with a positive number
	n := int64(100)
	for i := 0; i < 100; i++ {
		val := SecureIntn(n)
		assert.True(t, val >= 0 && val < n, "Value should be in [0, %d), but got %d", n, val)
	}

	// Test with a small positive number
	n = 2
	for i := 0; i < 10; i++ {
		val := SecureIntn(n)
		assert.True(t, val >= 0 && val < n, "Value should be in [0, %d), but got %d", n, val)
	}

	// Test with 1
	assert.Equal(t, int64(0), SecureIntn(1), "SecureIntn(1) should always be 0")

	// Test with 0 or negative
	assert.Equal(t, int64(0), SecureIntn(0), "SecureIntn(0) should be 0")
	assert.Equal(t, int64(0), SecureIntn(-10), "SecureIntn with negative input should be 0")
}

func TestSecureUint64n(t *testing.T) {
	// Test with a positive number
	n := uint64(100)
	for i := 0; i < 100; i++ {
		val := SecureUint64n(n)
		assert.True(t, val < n, "Value should be in [0, %d), but got %d", n, val)
	}

	// Test with a small positive number
	n = 2
	for i := 0; i < 10; i++ {
		val := SecureUint64n(n)
		assert.True(t, val < n, "Value should be in [0, %d), but got %d", n, val)
	}

	// Test with 1
	assert.Equal(t, uint64(0), SecureUint64n(1), "SecureUint64n(1) should always be 0")

	// Test with 0
	assert.Equal(t, uint64(0), SecureUint64n(0), "SecureUint64n(0) should be 0")
}

func TestSecureFloat64(t *testing.T) {
	for i := 0; i < 100; i++ {
		val := SecureFloat64()
		assert.True(t, val >= 0.0 && val < 1.0, "Value should be in [0.0, 1.0), but got %f", val)
	}
}
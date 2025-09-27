package twofactor

import (
	"k2ray/internal/config"
	"strings"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Load a dummy config for testing purposes
	config.LoadConfig("../../configs/system.env.example")
	m.Run()
}

func TestGenerateSecret(t *testing.T) {
	key, err := GenerateSecret()
	require.NoError(t, err)
	assert.NotNil(t, key)
	assert.Equal(t, "k2ray", key.Issuer())
	assert.NotEmpty(t, key.Secret())
}

func TestGenerateQRCode(t *testing.T) {
	key, err := GenerateSecret()
	require.NoError(t, err)

	qrCodeBytes, err := GenerateQRCode(key)
	require.NoError(t, err)
	assert.NotEmpty(t, qrCodeBytes)
	// A simple check to see if it looks like a PNG header
	assert.True(t, strings.HasPrefix(string(qrCodeBytes), "\x89PNG"))
}

func TestValidateCode(t *testing.T) {
	// This test is time-sensitive and can be flaky.
	// For a robust test suite, this would require mocking the time.
	// For now, we'll test the basic validation logic.

	key, err := GenerateSecret()
	require.NoError(t, err)

	// Generate a valid passcode
	passcode, err := totp.GenerateCode(key.Secret(), time.Now())
	require.NoError(t, err)

	// Test with a valid code
	assert.True(t, ValidateCode(key.Secret(), passcode))

	// Test with an invalid code
	assert.False(t, ValidateCode(key.Secret(), "000000"))
}

func TestGenerateRecoveryCodes(t *testing.T) {
	count := 10
	length := 12

	codes, err := GenerateRecoveryCodes(count, length)
	require.NoError(t, err)
	assert.Len(t, codes, count)

	// Check for uniqueness and length
	codeMap := make(map[string]bool)
	for _, code := range codes {
		assert.Len(t, code, length)
		assert.False(t, codeMap[code], "Found duplicate recovery code")
		codeMap[code] = true
	}
}

func TestRecoveryCodesToString(t *testing.T) {
	codes := []string{"abc", "def", "ghi"}
	str := RecoveryCodesToString(codes)
	assert.Equal(t, "abc,def,ghi", str)

	emptyCodes := []string{}
	str = RecoveryCodesToString(emptyCodes)
	assert.Equal(t, "", str)
}

func TestStringToRecoveryCodes(t *testing.T) {
	str := "abc,def,ghi"
	codes := StringToRecoveryCodes(str)
	assert.Equal(t, []string{"abc", "def", "ghi"}, codes)

	emptyStr := ""
	codes = StringToRecoveryCodes(emptyStr)
	assert.Equal(t, []string{}, codes)
}
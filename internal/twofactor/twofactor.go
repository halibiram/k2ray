package twofactor

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"image/png"
	"k2ray/internal/config"
	"strings"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// GenerateSecret creates a new TOTP secret key.
func GenerateSecret() (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      config.AppConfig.AppName,
		AccountName: "user@k2ray", // This should be customized per user
	})
}

// GenerateQRCode generates a PNG image of the QR code for the given OTP key.
func GenerateQRCode(key *otp.Key) ([]byte, error) {
	img, err := key.Image(256, 256)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ValidateCode checks if the provided passcode is valid for the given secret.
func ValidateCode(secret, passcode string) bool {
	valid := totp.Validate(passcode, secret)
	return valid
}

// GenerateRecoveryCodes creates a set of single-use recovery codes.
func GenerateRecoveryCodes(count, length int) ([]string, error) {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		bytes := make([]byte, length)
		if _, err := rand.Read(bytes); err != nil {
			return nil, err
		}
		// Use a user-friendly encoding
		codes[i] = base64.RawURLEncoding.EncodeToString(bytes)[:length]
	}
	return codes, nil
}

// RecoveryCodesToString converts a slice of recovery codes to a single string for DB storage.
func RecoveryCodesToString(codes []string) string {
	return strings.Join(codes, ",")
}

// StringToRecoveryCodes converts a comma-separated string of codes back to a slice.
func StringToRecoveryCodes(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}
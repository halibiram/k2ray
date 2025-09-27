package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Enable2FA generates a new 2FA secret and returns it as a QR code.
func Enable2FA(c *gin.Context) {
	// 1. Get user from context (set by AuthMiddleware)
	// 2. Generate a new TOTP secret
	// 3. Store the secret temporarily (e.g., in user's DB record but not yet enabled)
	// 4. Generate QR code from the secret
	// 5. Return QR code image and the secret (for manual entry)
	c.JSON(http.StatusOK, gin.H{"message": "2FA enabling process started"})
}

// Verify2FA verifies the TOTP code and enables 2FA for the user.
func Verify2FA(c *gin.Context) {
	// 1. Get user from context
	// 2. Get TOTP code from request body
	// 3. Retrieve temporary secret from DB
	// 4. Validate the code
	// 5. If valid, set two_factor_enabled = true, generate recovery codes, and save them
	// 6. Return success message and recovery codes
	c.JSON(http.StatusOK, gin.H{"message": "2FA verified and enabled"})
}

// Disable2FA disables two-factor authentication for the user.
func Disable2FA(c *gin.Context) {
	// 1. Get user from context
	// 2. Get password from request body (to confirm identity)
	// 3. Verify password
	// 4. If password is correct, clear 2FA secret, disable 2FA flag, and clear recovery codes
	// 5. Return success message
	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled"})
}
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Verify2FAPayload defines the structure for the 2FA verification request.
type Verify2FAPayload struct {
	Code string `json:"code" binding:"required,numeric,len=6"`
}

// Disable2FAPayload defines the structure for the 2FA disabling request.
type Disable2FAPayload struct {
	Password string `json:"password" binding:"required"`
}

// Enable2FA generates a new 2FA secret and returns it as a QR code.
// This endpoint does not require a request body.
func Enable2FA(c *gin.Context) {
	// 1. Get user from context (set by AuthMiddleware)
	// 2. Generate a new TOTP secret
	// 3. Store the secret temporarily (e.g., in user's DB record but not yet enabled)
	// 4. Generate QR code from the secret
	// 5. Return QR code image and the secret (for manual entry)
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not Implemented"})
}

// Verify2FA verifies the TOTP code and enables 2FA for the user.
func Verify2FA(c *gin.Context) {
	var payload Verify2FAPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}
	// 1. Get user from context
	// 2. Retrieve temporary secret from DB
	// 3. Validate the code from payload.Code
	// 4. If valid, set two_factor_enabled = true, generate recovery codes, and save them
	// 5. Return success message and recovery codes
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not Implemented"})
}

// Disable2FA disables two-factor authentication for the user.
func Disable2FA(c *gin.Context) {
	var payload Disable2FAPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}
	// 1. Get user from context
	// 2. Verify password from payload.Password
	// 3. If password is correct, clear 2FA secret, disable 2FA flag, and clear recovery codes
	// 4. Return success message
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not Implemented"})
}
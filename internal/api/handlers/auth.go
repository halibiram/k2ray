package handlers

import (
	"database/sql"
	"k2ray/internal/api/middleware"
	"k2ray/internal/auth"
	"k2ray/internal/db"
	"k2ray/internal/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// LoginPayload defines the expected JSON structure for a login request.
// Gin's `binding:"required"` tag ensures these fields are present in the request.
type LoginPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshPayload defines the expected JSON structure for a token refresh request.
type RefreshPayload struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Login is the handler for the user authentication endpoint.
func Login(c *gin.Context) {
	var payload LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// 1. Fetch user from the database by username
	user := &db.User{}
	err := db.DB.QueryRow("SELECT id, username, password_hash FROM users WHERE username = ?", payload.Username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		// If no user is found, return a generic "invalid credentials" error to prevent username enumeration.
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
		// For other database errors, log the error and return a generic server error.
		log.Printf("Database error on login for user %s: %v", payload.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// 2. Check if the provided password matches the stored hash
	if !utils.CheckPasswordHash(payload.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 3. Generate JWT access and refresh tokens
	accessToken, refreshToken, err := auth.GenerateTokens(user.Username)
	if err != nil {
		log.Printf("Token generation error for user %s: %v", payload.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate authentication tokens"})
		return
	}

	// 4. Return tokens to the client
	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Refresh is the handler for refreshing JWTs. It implements token rotation.
func Refresh(c *gin.Context) {
	var payload RefreshPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// 1. Validate the provided refresh token
	claims, err := auth.ValidateToken(payload.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// 2. Check if the refresh token has already been used (is blocklisted)
	isBlocklisted, err := db.IsTokenBlocklisted(claims.ID)
	if err != nil {
		log.Printf("Error checking token blocklist for JTI %s: %v", claims.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify token status"})
		return
	}
	if isBlocklisted {
		// This could indicate a stolen token is being re-used.
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Refresh token has already been used"})
		return
	}

	// 3. Invalidate the old refresh token by adding it to the blocklist.
	expiresAt := claims.ExpiresAt.Time
	if err := db.BlocklistToken(claims.ID, expiresAt); err != nil {
		log.Printf("Error blocklisting token for JTI %s: %v", claims.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not process refresh token"})
		return
	}

	// 4. Issue a new pair of tokens
	newAccessToken, newRefreshToken, err := auth.GenerateTokens(claims.Username)
	if err != nil {
		log.Printf("Token generation error for user %s: %v", claims.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate new tokens"})
		return
	}

	// 5. Return the new tokens
	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

// Logout invalidates the user's current access token by adding it to the blocklist.
func Logout(c *gin.Context) {
	// The AuthMiddleware has already validated the token and placed its details in the context.
	jtiVal, ok := c.Get(middleware.ContextTokenJTIKey)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JTI not found in context"})
		return
	}
	jti := jtiVal.(string)

	expiresAtVal, ok := c.Get(middleware.ContextTokenExpiresAtKey)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Expiration not found in context"})
		return
	}
	expiresAt := expiresAtVal.(time.Time)

	// Add the token to the blocklist.
	if err := db.BlocklistToken(jti, expiresAt); err != nil {
		log.Printf("Error blocklisting token for JTI %s: %v", jti, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not process logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

package handlers

import (
	"database/sql"
	"k2ray/internal/auth"
	"k2ray/internal/db"
	"k2ray/internal/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginPayload defines the expected JSON structure for a login request.
// Gin's `binding:"required"` tag ensures these fields are present in the request.
type LoginPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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

package handlers

import (
	"k2ray/internal/api/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMe retrieves information about the currently authenticated user
// from the context (which is populated by the AuthMiddleware).
func GetMe(c *gin.Context) {
	username, exists := c.Get(middleware.ContextUsernameKey)
	if !exists {
		// This should not happen if the middleware is correctly applied.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User information not found in context"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"username": username})
}

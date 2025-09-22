package middleware

import (
	"k2ray/internal/auth"
	"k2ray/internal/db"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// ContextUsernameKey is the key for storing the username in the Gin context.
	ContextUsernameKey = "username"
	// ContextTokenJTIKey is the key for storing the token's JTI in the Gin context.
	ContextTokenJTIKey = "jti"
	// ContextTokenExpiresAtKey is the key for storing the token's expiration time.
	ContextTokenExpiresAtKey = "expires_at"
)

// AuthMiddleware creates a Gin middleware for authenticating requests via JWT.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not provided"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Check if the token has been revoked (is in the blocklist).
		isBlocklisted, err := db.IsTokenBlocklisted(claims.ID)
		if err != nil {
			log.Printf("Error checking token blocklist for JTI %s: %v", claims.ID, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not verify token status"})
			return
		}
		if isBlocklisted {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			return
		}

		// Store user information in the context for downstream handlers to use.
		c.Set(ContextUsernameKey, claims.Username)
		c.Set(ContextTokenJTIKey, claims.ID)
		c.Set(ContextTokenExpiresAtKey, claims.ExpiresAt.Time)

		// Continue to the next handler.
		c.Next()
	}
}

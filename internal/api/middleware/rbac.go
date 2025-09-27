package middleware

import (
	"github.com/gin-gonic/gin"
	"k2ray/internal/auth"
	"k2ray/internal/db"
	"net/http"
)

// AdminRequired is a middleware that ensures the user has the 'admin' role.
// It should be used after the AuthMiddleware.
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user claims from the context, which should have been set by AuthMiddleware
		claims, exists := c.Get("user_claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User claims not found"})
			return
		}

		// Type assert the claims to the correct struct
		userClaims, ok := claims.(*auth.Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not parse user claims"})
			return
		}

		// Check if the user's role is 'admin'
		if userClaims.Role != db.AdminRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}

		// If the user is an admin, continue to the next handler
		c.Next()
	}
}
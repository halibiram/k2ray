package middleware

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimiterMiddleware creates a new rate-limiting middleware.
// It uses an in-memory store, which is suitable for single-instance deployments.
func RateLimiterMiddleware(rateFormat string) gin.HandlerFunc {
	// 1. Parse the rate format string (e.g., "10-M" for 10 requests per minute).
	rate, err := limiter.NewRateFromFormatted(rateFormat)
	if err != nil {
		log.Fatalf("Failed to parse rate limit format: %v", err)
		// This will cause the application to fail on startup if the format is invalid.
	}

	// 2. Create a new in-memory store.
	store := memory.NewStore()

	// 3. Create a new limiter instance.
	instance := limiter.New(store, rate)

	// 4. Return the Gin middleware handler.
	return func(c *gin.Context) {
		// Get the limiter context for the current request.
		context, err := instance.Get(c.Request.Context(), c.ClientIP())
		if err != nil {
			log.Printf("Rate limiter error: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Add rate limit headers to the response.
		c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

		// Abort the request if the limit has been reached.
		if context.Reached {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}

		c.Next()
	}
}
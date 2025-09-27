package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/redis"
	redis_internal "k2ray/internal/redis"
)

// RateLimiterMiddleware creates a new rate-limiting middleware with a dynamic key.
// It limits requests based on the user ID if available, otherwise it falls back to the client's IP address.
func RateLimiterMiddleware(rateFormat string) gin.HandlerFunc {
	// 1. Parse the rate format string.
	rate, err := limiter.NewRateFromFormatted(rateFormat)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to parse rate limit format: %s", rateFormat)
	}

	// 2. Create a new Redis store, using the global Redis client.
	store, err := redis.NewStore(redis_internal.RedisClient)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create redis store for rate limiter")
	}

	// 3. Create a new limiter instance.
	// The key is determined dynamically per request.
	instance := limiter.New(store, rate, limiter.WithTrustForwardHeader(true))

	// 4. Return the Gin middleware handler.
	return func(c *gin.Context) {
		// Determine the key for rate limiting.
		// Prioritize user ID if the user is authenticated.
		key := c.ClientIP()
		if userID, exists := c.Get("user_id"); exists {
			if idStr, ok := userID.(string); ok {
				key = idStr
			}
		}

		// Get the limiter context for the current request.
		context, err := instance.Get(c.Request.Context(), key)
		if err != nil {
			log.Error().Err(err).Str("key", key).Msg("Rate limiter error")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Add rate limit headers to the response.
		c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

		// Abort the request if the limit has been reached.
		if context.Reached {
			log.Warn().Str("key", key).Msg("Rate limit exceeded")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}

		c.Next()
	}
}
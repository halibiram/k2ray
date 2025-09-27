package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"k2ray/internal/metrics"
)

// PrometheusMiddleware creates a Gin middleware that collects Prometheus metrics for each request.
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// After the request has been processed, collect metrics
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// Record request duration
		metrics.HTTPRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration.Seconds())

		// Record total requests
		metrics.HTTPRequestsTotal.WithLabelValues(c.Request.Method, path, strconv.Itoa(statusCode)).Inc()
	}
}
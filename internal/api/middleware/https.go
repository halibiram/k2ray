package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HTTPSRedirectMiddleware redirects HTTP requests to HTTPS.
// It checks the 'X-Forwarded-Proto' header, which is set by many reverse proxies.
func HTTPSRedirectMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only redirect in "production" mode.
		// In development, we often don't use HTTPS.
		if gin.Mode() == gin.ReleaseMode {
			if c.GetHeader("X-Forwarded-Proto") != "https" {
				// Construct the new URL.
				target := "https://" + c.Request.Host + c.Request.URL.Path
				if c.Request.URL.RawQuery != "" {
					target += "?" + c.Request.URL.RawQuery
				}

				// Redirect with a 308 status code (Permanent Redirect).
				c.Redirect(http.StatusPermanentRedirect, target)
				c.Abort()
				return
			}
		}

		// If the request is already HTTPS or not in production, continue.
		c.Next()
	}
}
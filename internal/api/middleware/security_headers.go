package middleware

import (
	"github.com/gin-gonic/gin"
)

// CSPMiddleware sets a strict Content-Security-Policy header.
// This policy helps to prevent Cross-Site Scripting (XSS) and other code injection attacks.
func CSPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This is a very strict policy. You may need to adjust it based on your frontend's needs.
		// For example, if you use inline scripts or styles, you'll need to add 'unsafe-inline' or use hashes/nonces.
		// If you load assets from a CDN, you'll need to add that CDN's domain to the appropriate directives.
		csp := "default-src 'self'; " +
			"script-src 'self'; " + // Only allow scripts from the same origin.
			"style-src 'self' 'unsafe-inline'; " + // Allow inline styles for convenience, but consider removing 'unsafe-inline'.
			"img-src 'self' data:; " + // Allow images from the same origin and data URIs.
			"font-src 'self'; " +
			"object-src 'none'; " + // Disallow plugins like Flash.
			"frame-ancestors 'none'; " + // Prevent clickjacking.
			"form-action 'self'; " +
			"base-uri 'self';"
		c.Header("Content-Security-Policy", csp)
		c.Next()
	}
}

// SecurityHeadersMiddleware is a convenient wrapper for all security-related headers.
func SecurityHeadersMiddleware() gin.HandlerFunc {
	csp := CSPMiddleware()
	return func(c *gin.Context) {
		// Apply CSP
		csp(c)

		// Add other security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Next()
	}
}
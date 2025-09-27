package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

// CORSMiddleware configures and returns a new CORS middleware.
func CORSMiddleware() gin.HandlerFunc {
	// In a production environment, you should be more restrictive.
	// For example, allow only your frontend's origin.
	// cors.Config{
	//     AllowOrigins:     []string{"https://your-frontend.com"},
	//     AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	//     AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	//     ExposeHeaders:    []string{"Content-Length"},
	//     AllowCredentials: true,
	//     MaxAge:           12 * time.Hour,
	// }
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
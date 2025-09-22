package api

import (
	"github.com/gin-gonic/gin"
	"k2ray/internal/api/handlers"
)

// SetupRouter configures the routes for the application.
func SetupRouter(router *gin.Engine) {
	// All API routes will be prefixed with /api/v1
	apiV1 := router.Group("/api/v1")
	{
		systemRoutes := apiV1.Group("/system")
		{
			systemRoutes.GET("/status", handlers.SystemStatus)
		}

		authRoutes := apiV1.Group("/auth")
		{
			authRoutes.POST("/login", handlers.Login)
		}
	}
}

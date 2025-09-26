package api

import (
	"github.com/gin-gonic/gin"
	"k2ray/internal/api/handlers"
	"k2ray/internal/api/middleware"
)

// SetupRouter configures the routes for the application.
func SetupRouter(router *gin.Engine) {
	// All API routes will be prefixed with /api/v1
	apiV1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		systemRoutes := apiV1.Group("/system")
		{
			systemRoutes.GET("/status", handlers.SystemStatus)
		}

		authRoutes := apiV1.Group("/auth")
		{
			authRoutes.POST("/login", handlers.Login)
			authRoutes.POST("/refresh", handlers.Refresh)
		}

		// Protected routes (authentication required)
		protected := apiV1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/auth/logout", handlers.Logout)

			userRoutes := protected.Group("/users")
			{
				userRoutes.GET("/me", handlers.GetMe)
			}

			configRoutes := protected.Group("/configs")
			{
				configRoutes.POST("", handlers.CreateConfig)
				configRoutes.GET("", handlers.ListConfigs)
				configRoutes.GET("/:id", handlers.GetConfig)
				configRoutes.PUT("/:id", handlers.UpdateConfig)
				configRoutes.DELETE("/:id", handlers.DeleteConfig)
			}

			// Protected system routes
			protectedSystemRoutes := protected.Group("/system")
			{
				protectedSystemRoutes.GET("/active-config", handlers.GetActiveConfig)
				protectedSystemRoutes.POST("/active-config", handlers.SetActiveConfig)
				protectedSystemRoutes.GET("/info", handlers.GetSystemInfo)
				protectedSystemRoutes.GET("/logs", handlers.GetSystemLogs)
			}

			// V2Ray process management routes
			v2rayRoutes := protected.Group("/v2ray")
			{
				v2rayRoutes.POST("/start", handlers.StartV2Ray)
				v2rayRoutes.POST("/stop", handlers.StopV2Ray)
				v2rayRoutes.GET("/status", handlers.GetV2RayStatus)
			}

			// Metrics routes
			metricsRoutes := protected.Group("/metrics")
			{
				metricsRoutes.GET("/traffic", handlers.GetTrafficMetrics)
				metricsRoutes.GET("/connections", handlers.GetConnectionMetrics)
				metricsRoutes.GET("/performance", handlers.GetPerformanceMetrics)
			}

			// WebSocket route
			protected.GET("/ws", handlers.WebSocketHandler)
		}
	}
}

package handlers

import (
	"github.com/gin-gonic/gin"
	"k2ray/internal/system"
	"net/http"
)

// GetTrafficMetrics handles the request for traffic metrics.
func GetTrafficMetrics(c *gin.Context) {
	metrics := system.GetTrafficMetrics()
	c.JSON(http.StatusOK, metrics)
}

// GetConnectionMetrics handles the request for connection metrics.
func GetConnectionMetrics(c *gin.Context) {
	metrics := system.GetConnectionMetrics()
	c.JSON(http.StatusOK, metrics)
}

// GetPerformanceMetrics handles the request for performance metrics.
func GetPerformanceMetrics(c *gin.Context) {
	metrics := system.GetPerformanceMetrics()
	c.JSON(http.StatusOK, metrics)
}
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SystemStatus is a handler for the system status endpoint.
// It returns a simple JSON response indicating the service is running.
func SystemStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

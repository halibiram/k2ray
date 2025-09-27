package handlers

import (
	"database/sql"
	"k2ray/internal/api/middleware"
	"k2ray/internal/db"
	"k2ray/internal/system"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const ActiveConfigKey = "active_config_id"

// SystemStatus is a handler for the system status endpoint.
func SystemStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// SetActiveConfigPayload defines the structure for setting the active config.
type SetActiveConfigPayload struct {
	ConfigID int64 `json:"config_id" binding:"required"`
}

// SetActiveConfig sets the system-wide active V2Ray configuration.
func SetActiveConfig(c *gin.Context) {
	var payload SetActiveConfigPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	userID, _ := c.Get(middleware.ContextUserIDKey)

	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM configurations WHERE id = ? AND user_id = ?)", payload.ConfigID, userID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found or access denied"})
		return
	}

	upsertSQL := `INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value;`
	_, err = db.DB.Exec(upsertSQL, ActiveConfigKey, payload.ConfigID)
	if err != nil {
		log.Printf("Error setting active config: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set active configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Active configuration set successfully"})
}

// GetActiveConfig retrieves the currently active V2Ray configuration ID.
func GetActiveConfig(c *gin.Context) {
	var configID int64
	err := db.DB.QueryRow("SELECT value FROM settings WHERE key = ?", ActiveConfigKey).Scan(&configID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No active configuration is set"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get active configuration"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"active_config_id": configID})
}

// GetSystemInfo is the handler for the /system/info endpoint.
func GetSystemInfo(c *gin.Context) {
	info, err := system.GetSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve system information"})
		return
	}
	c.JSON(http.StatusOK, info)
}

// GetSystemLogs is the handler for the /system/logs endpoint.
func GetSystemLogs(c *gin.Context) {
	logs, err := system.GetSystemLogs()
	if err != nil {
		log.Printf("Error reading system logs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve system logs"})
		return
	}
	c.String(http.StatusOK, logs)
}

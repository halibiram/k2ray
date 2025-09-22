package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

import (
	"database/sql"
	"k2ray/internal/api/middleware"
	"k2ray/internal/db"
	"log"
)

const ActiveConfigKey = "active_config_id"

// SystemStatus is a handler for the system status endpoint.
// It returns a simple JSON response indicating the service is running.
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

	// Verify the config exists and belongs to the user.
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM v2ray_configs WHERE id = ? AND user_id = ?)", payload.ConfigID, userID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found or access denied"})
		return
	}

	// Upsert the setting into the system_settings table.
	upsertSQL := `INSERT INTO system_settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value;`
	_, err = db.DB.Exec(upsertSQL, ActiveConfigKey, payload.ConfigID)
	if err != nil {
		log.Printf("Error setting active config: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set active configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Active configuration set successfully"})
}

// GetActiveConfig retrieves the currently active V2Ray configuration.
func GetActiveConfig(c *gin.Context) {
	var configID int64
	err := db.DB.QueryRow("SELECT value FROM system_settings WHERE key = ?", ActiveConfigKey).Scan(&configID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No active configuration is set"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get active configuration"})
		return
	}

	// Optional: Fetch the full config details. For now, just return the ID.
	c.JSON(http.StatusOK, gin.H{"active_config_id": configID})
}

package handlers

import (
	"database/sql"
	"encoding/json"
	"k2ray/internal/api/middleware"
	"k2ray/internal/db"
	"k2ray/internal/v2ray"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateConfigPayload defines the structure for creating a new V2Ray config.
type CreateConfigPayload struct {
	Name       string          `json:"name" binding:"required"`
	Protocol   string          `json:"protocol" binding:"required"`
	ConfigData json.RawMessage `json:"config_data" binding:"required"`
}

// TransportSettings defines common transport settings for V2Ray protocols.
type TransportSettings struct {
	Network     string        `json:"net"`      // "tcp", "kcp", "ws", "h2", "quic", "grpc"
	Security    string        `json:"tls"`      // "none", "tls"
	WsSettings  WsSettings    `json:"wsSettings"`
	GrpcSettings GrpcSettings `json:"grpcSettings"`
}

// WsSettings defines WebSocket-specific transport settings.
type WsSettings struct {
	Path    string `json:"path"`
	Headers map[string]string `json:"headers"`
}

// GrpcSettings defines gRPC-specific transport settings.
type GrpcSettings struct {
	ServiceName string `json:"serviceName"`
}

// VmessConfigData defines the structure for a VMess config.
type VmessConfigData struct {
	V          string   `json:"v"`
	Add        string   `json:"add"`
	Port       any      `json:"port"`
	ID         string   `json:"id"`
	Aid        int      `json:"aid"`
	Type       string   `json:"type"` // Header type
	Host       string   `json:"host"`
	Path       string   `json:"path"`
	TransportSettings
}

// VlessConfigData defines the structure for a VLESS config.
type VlessConfigData struct {
	ID         string   `json:"id"`
	Address    string   `json:"add"`
	Port       any      `json:"port"`
	Encryption string   `json:"encryption"`
	Flow       string   `json:"flow"`
	TransportSettings
}

// ShadowsocksConfigData defines the structure for a Shadowsocks config.
type ShadowsocksConfigData struct {
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
	Password   string `json:"password"`
	Method     string `json:"method"`
	// Shadowsocks doesn't typically have complex transport settings in the same way,
	// but we can include basic ones if needed in the future.
}

// TrojanConfigData defines the structure for a Trojan config.
type TrojanConfigData struct {
	Server     string   `json:"server"`
	ServerPort int      `json:"server_port"`
	Password   string   `json:"password"`
	SNI        string   `json:"sni"`
	TransportSettings
}

// isValidatable defines an interface for config data structs.
type isValidatable interface {
	validate() error
}

func (c VmessConfigData) validate() error {
	if c.Add == "" || c.Port == nil || c.ID == "" {
		return &ValidationError{Msg: "VMess config must include 'add', 'port', and 'id'"}
	}
	return nil
}

func (c VlessConfigData) validate() error {
	if c.ID == "" || c.Address == "" || c.Port == nil {
		return &ValidationError{Msg: "VLESS config must include 'id', 'add', and 'port'"}
	}
	return nil
}

func (c ShadowsocksConfigData) validate() error {
	if c.Server == "" || c.ServerPort == 0 || c.Password == "" || c.Method == "" {
		return &ValidationError{Msg: "Shadowsocks config must include 'server', 'server_port', 'password', and 'method'"}
	}
	return nil
}

func (c TrojanConfigData) validate() error {
	if c.Server == "" || c.ServerPort == 0 || c.Password == "" {
		return &ValidationError{Msg: "Trojan config must include 'server', 'server_port', and 'password'"}
	}
	return nil
}


// ValidationError is a custom error type for validation failures.
type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string {
	return e.Msg
}

// validateAndDecode performs JSON unmarshaling and validation for a given protocol.
func validateAndDecode(protocol string, data json.RawMessage) (isValidatable, error) {
	var v isValidatable
	switch protocol {
	case "vmess":
		v = &VmessConfigData{}
	case "vless":
		v = &VlessConfigData{}
	case "shadowsocks":
		v = &ShadowsocksConfigData{}
	case "trojan":
		v = &TrojanConfigData{}
	default:
		return nil, &ValidationError{Msg: "Protocol not supported"}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return nil, &ValidationError{Msg: "Invalid config_data format: " + err.Error()}
	}

	if err := v.validate(); err != nil {
		return nil, err
	}

	return v, nil
}


// CreateConfig is the handler for creating a new V2Ray configuration.
func CreateConfig(c *gin.Context) {
	var payload CreateConfigPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Validate config_data based on the protocol
	if _, err := validateAndDecode(payload.Protocol, payload.ConfigData); err != nil {
		if verr, ok := err.(*ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": verr.Msg})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred during validation"})
		}
		return
	}

	// Get user ID from the context (set by AuthMiddleware)
	userIDVal, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}
	userID := userIDVal.(int64)

	// Insert into the database
	insertSQL := `INSERT INTO configurations (user_id, name, protocol, config_data) VALUES (?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(insertSQL)
	if err != nil {
		log.Printf("Error preparing SQL for CreateConfig: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create configuration"})
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(userID, payload.Name, payload.Protocol, string(payload.ConfigData))
	if err != nil {
		log.Printf("Error executing SQL for CreateConfig: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create configuration"})
		return
	}

	// Get the ID of the newly created config
	newID, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID for CreateConfig: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve new configuration ID"})
		return
	}

	// Return the newly created resource
	newConfig := db.Configuration{
		ID:         newID,
		UserID:     userID,
		Name:       payload.Name,
		Protocol:   payload.Protocol,
		ConfigData: string(payload.ConfigData),
		CreatedAt:  time.Now(), // Approximate, DB value is more accurate
		UpdatedAt:  time.Now(), // Approximate
	}

	c.JSON(http.StatusCreated, newConfig)
}

// ListConfigs retrieves all V2Ray configurations for the authenticated user.
func ListConfigs(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserIDKey)

	querySQL := `SELECT id, user_id, name, protocol, config_data, created_at, updated_at FROM configurations WHERE user_id = ?`
	rows, err := db.DB.Query(querySQL, userID)
	if err != nil {
		log.Printf("Error querying configs for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve configurations"})
		return
	}
	defer rows.Close()

	configs := []db.Configuration{}
	for rows.Next() {
		var config db.Configuration
		if err := rows.Scan(&config.ID, &config.UserID, &config.Name, &config.Protocol, &config.ConfigData, &config.CreatedAt, &config.UpdatedAt); err != nil {
			log.Printf("Error scanning config row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process configurations"})
			return
		}
		configs = append(configs, config)
	}

	c.JSON(http.StatusOK, configs)
}

// GetConfig retrieves a single V2Ray configuration by its ID.
func GetConfig(c *gin.Context) {
	configID := c.Param("id")
	userID, _ := c.Get(middleware.ContextUserIDKey)

	querySQL := `SELECT id, user_id, name, protocol, config_data, created_at, updated_at FROM configurations WHERE id = ? AND user_id = ?`
	var config db.Configuration
	err := db.DB.QueryRow(querySQL, configID, userID).Scan(&config.ID, &config.UserID, &config.Name, &config.Protocol, &config.ConfigData, &config.CreatedAt, &config.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found or access denied"})
			return
		}
		log.Printf("Error getting config %s for user %d: %v", configID, userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve configuration"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// UpdateConfigPayload defines the structure for updating a V2Ray config.
// Using pointers allows for partial updates (PATCH-like behavior).
type UpdateConfigPayload struct {
	Name       *string         `json:"name"`
	ConfigData *json.RawMessage `json:"config_data"`
}

// UpdateConfig updates a specific V2Ray configuration.
func UpdateConfig(c *gin.Context) {
	configID := c.Param("id")
	userID, _ := c.Get(middleware.ContextUserIDKey)

	// 1. Fetch the existing config to ensure it exists and the user owns it.
	var existingConfig db.Configuration
	querySQL := `SELECT id, user_id, name, protocol, config_data, created_at, updated_at FROM configurations WHERE id = ? AND user_id = ?`
	err := db.DB.QueryRow(querySQL, configID, userID).Scan(&existingConfig.ID, &existingConfig.UserID, &existingConfig.Name, &existingConfig.Protocol, &existingConfig.ConfigData, &existingConfig.CreatedAt, &existingConfig.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found or access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve configuration for update"})
		return
	}

	// 2. Bind the payload for the update.
	var payload UpdateConfigPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// 3. Apply updates from the payload to the existing config.
	if payload.Name != nil {
		existingConfig.Name = *payload.Name
	}
	if payload.ConfigData != nil {
		// Validate new config data before applying, based on the existing config's protocol
		if _, err := validateAndDecode(existingConfig.Protocol, *payload.ConfigData); err != nil {
			if verr, ok := err.(*ValidationError); ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": verr.Msg})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred during validation"})
			}
			return
		}
		existingConfig.ConfigData = string(*payload.ConfigData)
	}

	// 4. Save the updated record back to the database.
	updateSQL := `UPDATE configurations SET name = ?, config_data = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND user_id = ?`
	_, err = db.DB.Exec(updateSQL, existingConfig.Name, existingConfig.ConfigData, configID, userID)
	if err != nil {
		log.Printf("Error executing update for config %s: %v", configID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update configuration"})
		return
	}

	// 5. Return the full updated configuration.
	c.JSON(http.StatusOK, existingConfig)
}

// DeleteConfig deletes a specific V2Ray configuration.
func DeleteConfig(c *gin.Context) {
	configID := c.Param("id")
	userID, _ := c.Get(middleware.ContextUserIDKey)

	deleteSQL := `DELETE FROM configurations WHERE id = ? AND user_id = ?`
	res, err := db.DB.Exec(deleteSQL, configID, userID)
	if err != nil {
		log.Printf("Error deleting config %s: %v", configID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete configuration"})
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check deletion status"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found or access denied"})
		return
	}

	c.Status(http.StatusNoContent)
}

// --- V2Ray Process Management Handlers ---

// StartV2Ray starts the V2Ray service.
func StartV2Ray(c *gin.Context) {
	if err := v2ray.Start(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "V2Ray service started successfully (mocked)."})
}

// StopV2Ray stops the V2Ray service.
func StopV2Ray(c *gin.Context) {
	if err := v2ray.Stop(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "V2Ray service stopped successfully (mocked)."})
}

// GetV2RayStatus gets the current status of the V2Ray service.
func GetV2RayStatus(c *gin.Context) {
	isRunning, pid := v2ray.Status()
	status := "stopped"
	if isRunning {
		status = "running"
	}
	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"pid":    pid, // Will be 0 if not running
	})
}

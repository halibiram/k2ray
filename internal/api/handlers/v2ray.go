package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"k2ray/internal/api/middleware"
	"k2ray/internal/db"
	"k2ray/internal/security"
	"k2ray/internal/v2ray"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// CreateConfigPayload defines the structure for creating a new V2Ray config.
type CreateConfigPayload struct {
	Name       string          `json:"name" binding:"required,min=3,max=50"`
	Protocol   string          `json:"protocol" binding:"required,oneof=vmess vless shadowsocks trojan"`
	ConfigData json.RawMessage `json:"config_data" binding:"required"`
}

// TransportSettings defines common transport settings for V2Ray protocols.
type TransportSettings struct {
	Network      string       `json:"net"` // "tcp", "kcp", "ws", "h2", "quic", "grpc"
	Security     string       `json:"tls"` // "none", "tls"
	WsSettings   WsSettings   `json:"wsSettings"`
	GrpcSettings GrpcSettings `json:"grpcSettings"`
}

// WsSettings defines WebSocket-specific transport settings.
type WsSettings struct {
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
}

// GrpcSettings defines gRPC-specific transport settings.
type GrpcSettings struct {
	ServiceName string `json:"serviceName"`
}

// VmessConfigData defines the structure for a VMess config.
type VmessConfigData struct {
	V                 string   `json:"v"`
	Add               string   `json:"add"`
	Port              any      `json:"port"`
	ID                string   `json:"id"`
	Aid               int      `json:"aid"`
	Type              string   `json:"type"` // Header type
	Host              string   `json:"host"`
	Path              string   `json:"path"`
	TransportSettings `json:","`
}

// VlessConfigData defines the structure for a VLESS config.
type VlessConfigData struct {
	ID                string   `json:"id"`
	Address           string   `json:"add"`
	Port              any      `json:"port"`
	Encryption        string   `json:"encryption"`
	Flow              string   `json:"flow"`
	TransportSettings `json:","`
}

// ShadowsocksConfigData defines the structure for a Shadowsocks config.
type ShadowsocksConfigData struct {
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
	Password   string `json:"password"`
	Method     string `json:"method"`
}

// TrojanConfigData defines the structure for a Trojan config.
type TrojanConfigData struct {
	Server            string   `json:"server"`
	ServerPort        int      `json:"server_port"`
	Password          string   `json:"password"`
	SNI               string   `json:"sni"`
	TransportSettings `json:","`
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

	if _, err := validateAndDecode(payload.Protocol, payload.ConfigData); err != nil {
		if verr, ok := err.(*ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": verr.Msg})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred during validation"})
		}
		return
	}

	userIDVal, _ := c.Get(middleware.ContextUserIDKey)
	userID := userIDVal.(int64)

	insertSQL := `INSERT INTO configurations (user_id, name, protocol, config_data) VALUES (?, ?, ?, ?)`
	res, err := db.DB.Exec(insertSQL, userID, payload.Name, payload.Protocol, string(payload.ConfigData))
	if err != nil {
		log.Error().Err(err).Msg("Error executing SQL for CreateConfig")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create configuration"})
		return
	}

	newID, _ := res.LastInsertId()

	// Audit log
	details := fmt.Sprintf("Configuration '%s' created with protocol '%s'", payload.Name, payload.Protocol)
	security.LogEvent(c, security.ConfigCreated, newID, details)

	newConfig := db.Configuration{
		ID:         newID,
		UserID:     userID,
		Name:       payload.Name,
		Protocol:   payload.Protocol,
		ConfigData: string(payload.ConfigData),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	c.JSON(http.StatusCreated, newConfig)
}

// ListConfigs retrieves all V2Ray configurations for the authenticated user.
func ListConfigs(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserIDKey)

	querySQL := `SELECT id, user_id, name, protocol, config_data, created_at, updated_at FROM configurations WHERE user_id = ?`
	rows, err := db.DB.Query(querySQL, userID)
	if err != nil {
		log.Error().Err(err).Int64("user_id", userID.(int64)).Msg("Error querying configs for user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve configurations"})
		return
	}
	defer rows.Close()

	configs := []db.Configuration{}
	for rows.Next() {
		var config db.Configuration
		if err := rows.Scan(&config.ID, &config.UserID, &config.Name, &config.Protocol, &config.ConfigData, &config.CreatedAt, &config.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("Error scanning config row")
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
		log.Error().Err(err).Str("config_id", configID).Int64("user_id", userID.(int64)).Msg("Error getting config")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve configuration"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// UpdateConfigPayload defines the structure for updating a V2Ray config.
type UpdateConfigPayload struct {
	Name       *string         `json:"name" binding:"min=3,max=50"`
	ConfigData *json.RawMessage `json:"config_data"`
}

// UpdateConfig updates a specific V2Ray configuration.
func UpdateConfig(c *gin.Context) {
	configID := c.Param("id")
	userID, _ := c.Get(middleware.ContextUserIDKey)

	var existingConfig db.Configuration
	err := db.DB.QueryRow("SELECT id, user_id, name, protocol, config_data, created_at, updated_at FROM configurations WHERE id = ? AND user_id = ?", configID, userID).Scan(&existingConfig.ID, &existingConfig.UserID, &existingConfig.Name, &existingConfig.Protocol, &existingConfig.ConfigData, &existingConfig.CreatedAt, &existingConfig.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found or access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve configuration for update"})
		return
	}

	var payload UpdateConfigPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	originalName := existingConfig.Name
	if payload.Name != nil {
		existingConfig.Name = *payload.Name
	}
	if payload.ConfigData != nil {
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

	updateSQL := `UPDATE configurations SET name = ?, config_data = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND user_id = ?`
	_, err = db.DB.Exec(updateSQL, existingConfig.Name, existingConfig.ConfigData, configID, userID)
	if err != nil {
		log.Error().Err(err).Str("config_id", configID).Msg("Error executing update for config")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update configuration"})
		return
	}

	// Audit log
	configIDInt, _ := strconv.ParseInt(configID, 10, 64)
	details := fmt.Sprintf("Configuration '%s' (was '%s') updated", existingConfig.Name, originalName)
	security.LogEvent(c, security.ConfigUpdated, configIDInt, details)

	c.JSON(http.StatusOK, existingConfig)
}

// DeleteConfig deletes a specific V2Ray configuration.
func DeleteConfig(c *gin.Context) {
	configID := c.Param("id")
	userID, _ := c.Get(middleware.ContextUserIDKey)

	var configName string
	err := db.DB.QueryRow("SELECT name FROM configurations WHERE id = ? AND user_id = ?", configID, userID).Scan(&configName)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found or access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve configuration for deletion"})
		return
	}

	deleteSQL := `DELETE FROM configurations WHERE id = ? AND user_id = ?`
	_, err = db.DB.Exec(deleteSQL, configID, userID)
	if err != nil {
		log.Error().Err(err).Str("config_id", configID).Msg("Error deleting config")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete configuration"})
		return
	}

	// Audit log
	configIDInt, _ := strconv.ParseInt(configID, 10, 64)
	details := fmt.Sprintf("Configuration '%s' deleted", configName)
	security.LogEvent(c, security.ConfigDeleted, configIDInt, details)

	c.Status(http.StatusNoContent)
}

// --- V2Ray Process Management Handlers ---

func StartV2Ray(c *gin.Context) {
	if err := v2ray.Start(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "V2Ray service started successfully (mocked)."})
}

func StopV2Ray(c *gin.Context) {
	if err := v2ray.Stop(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "V2Ray service stopped successfully (mocked)."})
}

func GetV2RayStatus(c *gin.Context) {
	isRunning, pid := v2ray.Status()
	status := "stopped"
	if isRunning {
		status = "running"
	}
	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"pid":    pid,
	})
}
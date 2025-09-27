package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"k2ray/internal/api/middleware"
	"k2ray/internal/db"
	"k2ray/internal/security"
	"k2ray/internal/v2ray"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// CreateConfigPayload defines the structure for creating a new V2Ray config.
type CreateConfigPayload struct {
	Name       string      `json:"name" binding:"required,min=3,max=50"`
	Protocol   string      `json:"protocol" binding:"required,oneof=vmess vless shadowsocks trojan"`
	ConfigData interface{} `json:"config_data" binding:"required" swaggertype:"object"`
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

// CreateConfig godoc
// @Summary Create a new V2Ray configuration
// @Description Creates a new V2Ray configuration for the authenticated user.
// @Tags Configs
// @Accept  json
// @Produce  json
// @Param   config body CreateConfigPayload true "New Configuration Details"
// @Success 201 {object} db.Configuration
// @Failure 400 {object} middleware.ErrorResponse "Invalid request payload or config data"
// @Failure 500 {object} middleware.ErrorResponse "Failed to create configuration"
// @Security ApiKeyAuth
// @Router /configs [post]
func CreateConfig(c *gin.Context) {
	var payload CreateConfigPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.Error(err)
		return
	}

	// Marshal the interface{} back to JSON bytes for validation and storage
	configDataBytes, err := json.Marshal(payload.ConfigData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format for config_data"})
		return
	}

	if _, err := validateAndDecode(payload.Protocol, configDataBytes); err != nil {
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
	res, err := db.DB.Exec(insertSQL, userID, payload.Name, payload.Protocol, string(configDataBytes))
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
		ConfigData: string(configDataBytes),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	c.JSON(http.StatusCreated, newConfig)
}

// PaginatedConfigsResponse is the structured response for a list of configs with pagination.
type PaginatedConfigsResponse struct {
	Data       []db.Configuration `json:"data"`
	Pagination PaginationMeta     `json:"pagination"`
}

// ListConfigs godoc
// @Summary List V2Ray configurations
// @Description Retrieves a paginated list of V2Ray configurations for the authenticated user.
// @Tags Configs
// @Accept  json
// @Produce  json
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param sort_by query string false "Field to sort by (id, name, protocol, created_at, updated_at)" default(id)
// @Param order query string false "Sort order (ASC, DESC)" default(ASC)
// @Param name query string false "Filter by configuration name (partial match)"
// @Param protocol query string false "Filter by protocol (vmess, vless, etc.)"
// @Success 200 {object} PaginatedConfigsResponse
// @Failure 500 {object} middleware.ErrorResponse "Failed to retrieve configurations"
// @Security ApiKeyAuth
// @Router /configs [get]
func ListConfigs(c *gin.Context) {
	// 1. Get authenticated user ID
	userID, _ := c.Get(middleware.ContextUserIDKey)

	// 2. Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sort_by", "id")
	order := strings.ToUpper(c.DefaultQuery("order", "ASC"))
	filterName := c.Query("name")
	filterProtocol := c.Query("protocol")

	// 3. Validate and sanitize inputs
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}
	allowedSortColumns := map[string]bool{"id": true, "name": true, "protocol": true, "created_at": true, "updated_at": true}
	if !allowedSortColumns[sortBy] {
		sortBy = "id"
	}

	// 4. Build the database query
	var args []interface{}
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString("FROM configurations WHERE user_id = ?")
	args = append(args, userID)

	if filterName != "" {
		queryBuilder.WriteString(" AND name LIKE ?")
		args = append(args, "%"+filterName+"%")
	}
	if filterProtocol != "" {
		queryBuilder.WriteString(" AND protocol = ?")
		args = append(args, filterProtocol)
	}

	// 5. Get total count for pagination
	var totalItems int
	countQuery := "SELECT COUNT(*) " + queryBuilder.String()
	err := db.DB.QueryRow(countQuery, args...).Scan(&totalItems)
	if err != nil {
		log.Error().Err(err).Msg("Error counting configurations")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve configurations"})
		return
	}

	// 6. Execute main query
	offset := (page - 1) * limit
	selectQuery := fmt.Sprintf("SELECT id, user_id, name, protocol, config_data, created_at, updated_at %s ORDER BY %s %s LIMIT ? OFFSET ?", queryBuilder.String(), sortBy, order)
	rows, err := db.DB.Query(selectQuery, append(args, limit, offset)...)
	if err != nil {
		log.Error().Err(err).Msg("Error querying configurations with pagination")
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

	// 7. Construct response
	response := PaginatedConfigsResponse{
		Data: configs,
		Pagination: PaginationMeta{
			TotalItems:   totalItems,
			TotalPages:   int(math.Ceil(float64(totalItems) / float64(limit))),
			CurrentPage:  page,
			ItemsPerPage: limit,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetConfig godoc
// @Summary Get a single V2Ray configuration
// @Description Retrieves details for a single V2Ray configuration by its ID.
// @Tags Configs
// @Accept  json
// @Produce  json
// @Param id path int true "Configuration ID"
// @Success 200 {object} db.Configuration
// @Failure 404 {object} middleware.ErrorResponse "Configuration not found or access denied"
// @Failure 500 {object} middleware.ErrorResponse "Failed to retrieve configuration"
// @Security ApiKeyAuth
// @Router /configs/{id} [get]
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
	Name       *string      `json:"name" binding:"min=3,max=50"`
	ConfigData *interface{} `json:"config_data" swaggertype:"object"`
}

// UpdateConfig godoc
// @Summary Update a V2Ray configuration
// @Description Updates a specific V2Ray configuration for the authenticated user.
// @Tags Configs
// @Accept  json
// @Produce  json
// @Param id path int true "Configuration ID"
// @Param config body UpdateConfigPayload true "Updated Configuration Details"
// @Success 200 {object} db.Configuration
// @Failure 400 {object} middleware.ErrorResponse "Invalid request payload or config data"
// @Failure 404 {object} middleware.ErrorResponse "Configuration not found or access denied"
// @Failure 500 {object} middleware.ErrorResponse "Failed to update configuration"
// @Security ApiKeyAuth
// @Router /configs/{id} [put]
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
		c.Error(err)
		return
	}

	originalName := existingConfig.Name
	if payload.Name != nil {
		existingConfig.Name = *payload.Name
	}
	if payload.ConfigData != nil {
		// Marshal the interface{} back to JSON bytes
		configDataBytes, err := json.Marshal(payload.ConfigData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format for config_data"})
			return
		}

		if _, err := validateAndDecode(existingConfig.Protocol, configDataBytes); err != nil {
			if verr, ok := err.(*ValidationError); ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": verr.Msg})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred during validation"})
			}
			return
		}
		existingConfig.ConfigData = string(configDataBytes)
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

// BulkDeleteConfigs godoc
// @Summary Bulk delete V2Ray configurations
// @Description Deletes multiple V2Ray configurations at once based on a list of IDs.
// @Tags Configs
// @Accept  json
// @Produce  json
// @Param ids body BulkDeleteRequest true "Configuration IDs to delete"
// @Success 200 {object} map[string]interface{} "message: Configurations deleted successfully, deleted_count: count"
// @Failure 400 {object} middleware.ErrorResponse "Invalid request payload"
// @Failure 500 {object} middleware.ErrorResponse "Failed to delete configurations"
// @Security ApiKeyAuth
// @Router /configs/bulk-delete [post]
func BulkDeleteConfigs(c *gin.Context) {
	var req BulkDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	userID, _ := c.Get(middleware.ContextUserIDKey)

	// Build the IN clause for the SQL query
	query := "DELETE FROM configurations WHERE user_id = ? AND id IN (?" + strings.Repeat(",?", len(req.IDs)-1) + ")"
	args := make([]interface{}, len(req.IDs)+1)
	args[0] = userID
	for i, id := range req.IDs {
		args[i+1] = id
	}

	res, err := db.DB.Exec(query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bulk delete configurations")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete configurations"})
		return
	}

	rowsAffected, _ := res.RowsAffected()

	// Audit log
	details := fmt.Sprintf("Bulk deleted %d configurations with IDs: %v", rowsAffected, req.IDs)
	security.LogEvent(c, security.ConfigDeleted, 0, details) // ID 0 for system/bulk action

	c.JSON(http.StatusOK, gin.H{
		"message":       "Configurations deleted successfully",
		"deleted_count": rowsAffected,
	})
}

// DeleteConfig godoc
// @Summary Delete a V2Ray configuration
// @Description Deletes a single V2Ray configuration by its ID.
// @Tags Configs
// @Accept  json
// @Produce  json
// @Param id path int true "Configuration ID"
// @Success 204 "No Content"
// @Failure 404 {object} middleware.ErrorResponse "Configuration not found or access denied"
// @Failure 500 {object} middleware.ErrorResponse "Failed to delete configuration"
// @Security ApiKeyAuth
// @Router /configs/{id} [delete]
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
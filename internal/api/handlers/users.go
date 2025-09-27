package handlers

import (
	"database/sql"
	"fmt"
	"k2ray/internal/api/middleware"
	"k2ray/internal/db"
	"k2ray/internal/security"
	"k2ray/internal/utils"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// CreateUserRequest defines the payload for creating a new user.
type CreateUserRequest struct {
	Username string      `json:"username" binding:"required,min=3,max=30"`
	Password string      `json:"password" binding:"required,min=8,max=100"`
	Role     db.UserRole `json:"role" binding:"required,oneof=admin user"`
}

// UpdateUserRequest defines the payload for updating a user's role.
type UpdateUserRequest struct {
	Role db.UserRole `json:"role" binding:"required,oneof=admin user"`
}

// UserResponse is the sanitized user object returned by the API.
type UserResponse struct {
	ID       int64       `json:"id"`
	Username string      `json:"username"`
	Role     db.UserRole `json:"role"`
}

// sanitizeUser creates a UserResponse from a db.User to hide sensitive fields.
func sanitizeUser(user db.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user with a username, password, and role. Only accessible by admins.
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   user body CreateUserRequest true "New User Details"
// @Success 201 {object} UserResponse
// @Failure 400 {object} middleware.ErrorResponse "Invalid request payload"
// @Failure 409 {object} middleware.ErrorResponse "Username already exists"
// @Failure 500 {object} middleware.ErrorResponse "Failed to create user"
// @Security ApiKeyAuth
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Str("username", req.Username).Msg("Error hashing password for new user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	insertSQL := `INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)`
	res, err := db.DB.Exec(insertSQL, req.Username, passwordHash, req.Role)
	if err != nil {
		// Handle potential unique constraint violation for username
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	newID, _ := res.LastInsertId()
	newUser := db.User{ID: newID, Username: req.Username, Role: req.Role}

	// Audit log
	security.LogEvent(c, security.UserCreated, newID, fmt.Sprintf("New user '%s' created with role '%s'", req.Username, req.Role))

	c.JSON(http.StatusCreated, sanitizeUser(newUser))
}

// PaginatedUsersResponse is the structured response for a list of users with pagination.
type PaginatedUsersResponse struct {
	Data       []UserResponse `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// ListUsers godoc
// @Summary List users
// @Description Retrieves a paginated list of users with optional filtering and sorting.
// @Tags Users
// @Accept  json
// @Produce  json
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param sort_by query string false "Field to sort by (id, username, role)" default(id)
// @Param order query string false "Sort order (ASC, DESC)" default(ASC)
// @Param role query string false "Filter by user role (admin, user)"
// @Param username query string false "Filter by username (partial match)"
// @Success 200 {object} PaginatedUsersResponse
// @Failure 500 {object} middleware.ErrorResponse "Failed to retrieve users"
// @Security ApiKeyAuth
// @Router /users [get]
func ListUsers(c *gin.Context) {
	// 1. Parse query parameters for pagination, sorting, and filtering
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sort_by", "id")
	order := strings.ToUpper(c.DefaultQuery("order", "ASC"))
	filterRole := c.Query("role")
	filterUsername := c.Query("username")

	// 2. Validate and sanitize inputs
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10 // Default to 10 if limit is invalid or too large
	}
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}
	// Whitelist sortable columns to prevent SQL injection
	allowedSortColumns := map[string]bool{"id": true, "username": true, "role": true}
	if !allowedSortColumns[sortBy] {
		sortBy = "id"
	}

	// 3. Build the database query dynamically
	var args []interface{}
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString("FROM users WHERE 1=1")

	if filterRole != "" {
		queryBuilder.WriteString(" AND role = ?")
		args = append(args, filterRole)
	}
	if filterUsername != "" {
		queryBuilder.WriteString(" AND username LIKE ?")
		args = append(args, "%"+filterUsername+"%")
	}

	// 4. Get the total count for pagination
	var totalItems int
	countQuery := "SELECT COUNT(*) " + queryBuilder.String()
	err := db.DB.QueryRow(countQuery, args...).Scan(&totalItems)
	if err != nil {
		log.Error().Err(err).Msg("Error counting users")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	// 5. Execute the main query to get the paginated data
	offset := (page - 1) * limit
	selectQuery := fmt.Sprintf("SELECT id, username, role %s ORDER BY %s %s LIMIT ? OFFSET ?", queryBuilder.String(), sortBy, order)
	rows, err := db.DB.Query(selectQuery, append(args, limit, offset)...)
	if err != nil {
		log.Error().Err(err).Msg("Error querying users with pagination")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	defer rows.Close()

	users := []UserResponse{}
	for rows.Next() {
		var user db.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Role); err != nil {
			log.Error().Err(err).Msg("Error scanning user row")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user data"})
			return
		}
		users = append(users, sanitizeUser(user))
	}

	// 6. Construct the paginated response
	response := PaginatedUsersResponse{
		Data: users,
		Pagination: PaginationMeta{
			TotalItems:   totalItems,
			TotalPages:   int(math.Ceil(float64(totalItems) / float64(limit))),
			CurrentPage:  page,
			ItemsPerPage: limit,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetUser godoc
// @Summary Get a single user
// @Description Retrieves details for a single user by their ID.
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} middleware.ErrorResponse "Invalid user ID"
// @Failure 404 {object} middleware.ErrorResponse "User not found"
// @Failure 500 {object} middleware.ErrorResponse "Failed to retrieve user"
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user db.User
	err = db.DB.QueryRow("SELECT id, username, role FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	c.JSON(http.StatusOK, sanitizeUser(user))
}

// UpdateUser godoc
// @Summary Update a user's role
// @Description Updates the role of a specific user. Admins cannot change their own role.
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param role body UpdateUserRequest true "New Role"
// @Success 200 {object} map[string]string "message: User updated successfully"
// @Failure 400 {object} middleware.ErrorResponse "Invalid user ID or request payload"
// @Failure 403 {object} middleware.ErrorResponse "Forbidden action"
// @Failure 500 {object} middleware.ErrorResponse "Failed to update user"
// @Security ApiKeyAuth
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	targetUserID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	// Prevent an admin from changing their own role to non-admin
	actorUserID, _ := c.Get(middleware.ContextUserIDKey)
	if actorUserID.(int64) == targetUserID && req.Role != db.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admins cannot remove their own admin role"})
		return
	}

	updateSQL := `UPDATE users SET role = ? WHERE id = ?`
	_, err = db.DB.Exec(updateSQL, req.Role, targetUserID)
	if err != nil {
		log.Error().Err(err).Int64("target_user_id", targetUserID).Msg("Failed to update user role")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Audit log
	details := fmt.Sprintf("User's role updated to '%s'", req.Role)
	security.LogEvent(c, security.UserUpdated, targetUserID, details)

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// BulkDeleteUsers godoc
// @Summary Bulk delete users
// @Description Deletes multiple users at once based on a list of IDs.
// @Tags Users
// @Accept  json
// @Produce  json
// @Param ids body BulkDeleteRequest true "User IDs to delete"
// @Success 200 {object} map[string]interface{} "message: Users deleted successfully, deleted_count: count"
// @Failure 400 {object} middleware.ErrorResponse "Invalid request payload"
// @Failure 403 {object} middleware.ErrorResponse "Forbidden action"
// @Failure 500 {object} middleware.ErrorResponse "Failed to delete users"
// @Security ApiKeyAuth
// @Router /users/bulk-delete [post]
func BulkDeleteUsers(c *gin.Context) {
	var req BulkDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	actorUserID, _ := c.Get(middleware.ContextUserIDKey)
	for _, id := range req.IDs {
		if id == actorUserID.(int64) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You cannot delete your own account as part of a bulk operation"})
			return
		}
	}

	// Build the IN clause for the SQL query
	query := "DELETE FROM users WHERE id IN (?" + strings.Repeat(",?", len(req.IDs)-1) + ")"
	args := make([]interface{}, len(req.IDs))
	for i, id := range req.IDs {
		args[i] = id
	}

	res, err := db.DB.Exec(query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bulk delete users")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete users"})
		return
	}

	rowsAffected, _ := res.RowsAffected()

	// Audit log
	details := fmt.Sprintf("Bulk deleted %d users with IDs: %v", rowsAffected, req.IDs)
	security.LogEvent(c, security.UserDeleted, 0, details) // ID 0 for system/bulk action

	c.JSON(http.StatusOK, gin.H{
		"message":       "Users deleted successfully",
		"deleted_count": rowsAffected,
	})
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Deletes a single user by their ID.
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} middleware.ErrorResponse "Invalid user ID"
// @Failure 403 {object} middleware.ErrorResponse "Forbidden action"
// @Failure 404 {object} middleware.ErrorResponse "User not found"
// @Failure 500 {object} middleware.ErrorResponse "Failed to delete user"
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	targetUserID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Prevent a user from deleting themselves
	actorUserID, _ := c.Get(middleware.ContextUserIDKey)
	if actorUserID.(int64) == targetUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot delete your own account"})
		return
	}

	deleteSQL := `DELETE FROM users WHERE id = ?`
	res, err := db.DB.Exec(deleteSQL, targetUserID)
	if err != nil {
		log.Error().Err(err).Int64("target_user_id", targetUserID).Msg("Failed to delete user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Audit log
	security.LogEvent(c, security.UserDeleted, targetUserID, "User account deleted")

	c.Status(http.StatusNoContent)
}

// GetMe retrieves the currently authenticated user's details.
func GetMe(c *gin.Context) {
	userID, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		// This should technically be caught by the AuthMiddleware, but we check as a safeguard.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user db.User
	err := db.DB.QueryRow("SELECT id, username, role FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			// This case is unlikely if middleware is correct but handled for robustness.
			c.JSON(http.StatusNotFound, gin.H{"error": "Authenticated user not found in database"})
			return
		}
		log.Error().Err(err).Msgf("Failed to retrieve data for user ID: %v", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information"})
		return
	}

	c.JSON(http.StatusOK, sanitizeUser(user))
}

// ErrorResponse is a generic error response.
type ErrorResponse struct {
	Error string `json:"error"`
}
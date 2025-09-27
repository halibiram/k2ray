package handlers

import (
	"database/sql"
	"fmt"
	"k2ray/internal/api/middleware"
	"k2ray/internal/db"
	"k2ray/internal/security"
	"k2ray/internal/utils"
	"net/http"
	"strconv"

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

// CreateUser handles the creation of a new user.
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
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

// ListUsers retrieves a list of all users.
func ListUsers(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, username, role FROM users")
	if err != nil {
		log.Error().Err(err).Msg("Error querying users")
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

	c.JSON(http.StatusOK, users)
}

// GetUser retrieves a single user by their ID.
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

// UpdateUser updates a user's role.
func UpdateUser(c *gin.Context) {
	targetUserID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
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

// DeleteUser deletes a user by their ID.
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
package handlers

import (
	"database/sql"
	"fmt"
	"k2ray/internal/api/middleware"
	"k2ray/internal/auth"
	"k2ray/internal/db"
	"k2ray/internal/security"
	"k2ray/internal/twofactor"
	"k2ray/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// LoginPayload defines the expected JSON structure for a login request.
type LoginPayload struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=8"`
}

// Login2FAPayload defines the structure for the 2FA verification step.
type Login2FAPayload struct {
	TwoFactorToken string `json:"two_factor_token" binding:"required"`
	Code           string `json:"code" binding:"required,numeric,len=6"`
}

// RefreshPayload defines the expected JSON structure for a token refresh request.
type RefreshPayload struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Login is the handler for the user authentication endpoint.
func Login(c *gin.Context) {
	var payload LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	ip := c.ClientIP()
	username := payload.Username

	if security.IsLockedOut(username) || security.IsLockedOut(ip) {
		details := fmt.Sprintf("Attempted login for locked out user '%s'", username)
		security.LogEvent(c, security.LoginFailure, 0, details)
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many failed login attempts. Please try again later."})
		return
	}

	user := &db.User{}
	err := db.DB.QueryRow("SELECT id, username, password_hash, role, two_factor_enabled FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.TwoFactorEnabled)
	if err != nil {
		if err == sql.ErrNoRows {
			security.RecordFailedAttempt(username)
			security.RecordFailedAttempt(ip)
			security.LogEvent(c, security.LoginFailure, 0, "Invalid username")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
		log.Error().Err(err).Str("username", username).Msg("Database error on login")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if !utils.CheckPasswordHash(payload.Password, user.PasswordHash) {
		security.RecordFailedAttempt(user.Username)
		security.RecordFailedAttempt(ip)
		security.LogEvent(c, security.LoginFailure, user.ID, "Invalid password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if user.TwoFactorEnabled {
		twoFactorToken, err := auth.Generate2FAToken(user.ID, user.Username)
		if err != nil {
			log.Error().Err(err).Str("username", username).Msg("2FA token generation error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not initiate 2FA process"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":          "2FA code required",
			"two_factor_token": twoFactorToken,
		})
		return
	}

	security.ResetAttempts(user.Username)
	security.ResetAttempts(ip)
	security.LogEvent(c, security.LoginSuccess, user.ID, "Login successful, no 2FA")

	accessToken, refreshToken, err := auth.GenerateTokens(*user)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("Token generation error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate authentication tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Login2FA handles the second step of authentication for users with 2FA enabled.
func Login2FA(c *gin.Context) {
	var payload Login2FAPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	ip := c.ClientIP()

	claims, err := auth.Validate2FAToken(payload.TwoFactorToken)
	if err != nil {
		security.LogEvent(c, security.TwoFactorFailure, 0, "Invalid 2FA token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired 2FA token"})
		return
	}

	username := claims.Username
	userID := claims.UserID

	if security.IsLockedOut(username) || security.IsLockedOut(ip) {
		details := fmt.Sprintf("Attempted 2FA for locked out user '%s'", username)
		security.LogEvent(c, security.TwoFactorFailure, userID, details)
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many failed login attempts. Please try again later."})
		return
	}

	user := &db.User{}
	var twoFactorSecret sql.NullString
	err = db.DB.QueryRow("SELECT id, username, role, two_factor_secret FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Role, &twoFactorSecret)
	if err != nil || !twoFactorSecret.Valid {
		log.Error().Err(err).Int64("user_id", userID).Msg("Could not retrieve 2FA secret")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify 2FA code"})
		return
	}

	if !twofactor.ValidateCode(twoFactorSecret.String, payload.Code) {
		security.RecordFailedAttempt(username)
		security.RecordFailedAttempt(ip)
		security.LogEvent(c, security.TwoFactorFailure, userID, "Invalid 2FA code")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid 2FA code"})
		return
	}

	security.ResetAttempts(username)
	security.ResetAttempts(ip)
	security.LogEvent(c, security.TwoFactorSuccess, userID, "2FA verification successful")
	security.LogEvent(c, security.LoginSuccess, userID, "Login successful with 2FA")

	accessToken, refreshToken, err := auth.GenerateTokens(*user)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("Token generation error after 2FA")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate authentication tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Refresh is the handler for refreshing JWTs. It implements token rotation.
func Refresh(c *gin.Context) {
	var payload RefreshPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	claims, err := auth.ValidateToken(payload.RefreshToken)
	if err != nil {
		security.LogEvent(c, security.TokenRefreshFailure, 0, "Invalid refresh token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	isBlocklisted, err := db.IsTokenBlocklisted(claims.ID)
	if err != nil {
		log.Error().Err(err).Str("jti", claims.ID).Msg("Error checking token blocklist")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify token status"})
		return
	}
	if isBlocklisted {
		security.LogEvent(c, security.TokenRefreshFailure, claims.UserID, "Attempted to use a blocklisted refresh token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Refresh token has already been used"})
		return
	}

	if err := db.BlocklistToken(claims.ID, claims.ExpiresAt.Time); err != nil {
		log.Error().Err(err).Str("jti", claims.ID).Msg("Error blocklisting token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not process refresh token"})
		return
	}

	user := db.User{ID: claims.UserID, Username: claims.Username, Role: claims.Role}
	newAccessToken, newRefreshToken, err := auth.GenerateTokens(user)
	if err != nil {
		log.Error().Err(err).Str("username", claims.Username).Msg("Token generation error on refresh")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate new tokens"})
		return
	}

	security.LogEvent(c, security.TokenRefreshSuccess, claims.UserID, "Token refreshed successfully")
	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

// Logout invalidates the user's current access token by adding it to the blocklist.
func Logout(c *gin.Context) {
	jti, ok := c.Get(middleware.ContextTokenJTIKey)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JTI not found in context"})
		return
	}

	expiresAt, ok := c.Get(middleware.ContextTokenExpiresAtKey)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Expiration not found in context"})
		return
	}

	if err := db.BlocklistToken(jti.(string), expiresAt.(time.Time)); err != nil {
		log.Error().Err(err).Str("jti", jti.(string)).Msg("Error blocklisting token on logout")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not process logout"})
		return
	}

	userID, _ := c.Get(middleware.ContextUserIDKey)
	security.LogEvent(c, security.LogoutSuccess, userID.(int64), "User logged out successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
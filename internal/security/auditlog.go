package security

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// AuditEventType defines the type of event being logged.
type AuditEventType string

const (
	// Authentication Events
	LoginSuccess      AuditEventType = "LOGIN_SUCCESS"
	LoginFailure      AuditEventType = "LOGIN_FAILURE"
	LogoutSuccess     AuditEventType = "LOGOUT_SUCCESS"
	TokenRefreshSuccess AuditEventType = "TOKEN_REFRESH_SUCCESS"
	TokenRefreshFailure AuditEventType = "TOKEN_REFRESH_FAILURE"
	TwoFactorSuccess  AuditEventType = "2FA_SUCCESS"
	TwoFactorFailure  AuditEventType = "2FA_FAILURE"

	// User Management Events
	UserCreated AuditEventType = "USER_CREATED"
	UserUpdated AuditEventType = "USER_UPDATED"
	UserDeleted AuditEventType = "USER_DELETED"

	// Config Management Events
	ConfigCreated AuditEventType = "CONFIG_CREATED"
	ConfigUpdated AuditEventType = "CONFIG_UPDATED"
	ConfigDeleted AuditEventType = "CONFIG_DELETED"
)

// AuditEvent represents a security-sensitive event that should be logged.
type AuditEvent struct {
	Timestamp time.Time      `json:"timestamp"`
	Type      AuditEventType `json:"type"`
	UserID    int64          `json:"user_id,omitempty"`   // The user who performed the action.
	TargetID  int64          `json:"target_id,omitempty"` // The user/object that was affected.
	ClientIP  string         `json:"client_ip"`
	Details   string         `json:"details,omitempty"` // Additional context about the event.
}

// LogEvent creates and logs a new audit event.
// It's designed to be called from within a Gin context to automatically capture IP and UserID.
func LogEvent(c *gin.Context, eventType AuditEventType, targetID int64, details string) {
	var userID int64
	// Try to get UserID from the context (it will be present for authenticated routes).
	if id, exists := c.Get("user_id"); exists {
		if u, ok := id.(int64); ok {
			userID = u
		}
	}

	event := AuditEvent{
		Timestamp: time.Now().UTC(),
		Type:      eventType,
		UserID:    userID,
		TargetID:  targetID,
		ClientIP:  c.ClientIP(),
		Details:   details,
	}

	// Log the event as a structured object.
	log.Info().
		Str("log_type", "audit").
		Object("event", event).
		Msg("Audit event recorded")
}

// MarshalZerologObject implements the zerolog.LogObjectMarshaler interface for AuditEvent.
func (e AuditEvent) MarshalZerologObject(ze *zerolog.Event) {
	ze.Time("timestamp", e.Timestamp)
	ze.Str("type", string(e.Type))
	if e.UserID != 0 {
		ze.Int64("user_id", e.UserID)
	}
	if e.TargetID != 0 {
		ze.Int64("target_id", e.TargetID)
	}
	ze.Str("client_ip", e.ClientIP)
	if e.Details != "" {
		ze.Str("details", e.Details)
	}
}
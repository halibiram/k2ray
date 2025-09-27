package db

import (
	"database/sql"
	"time"
)

// User represents a user in the system.
type User struct {
	ID                     int64
	Username               string
	PasswordHash           string
	TwoFactorSecret        sql.NullString
	TwoFactorEnabled       bool
	TwoFactorRecoveryCodes sql.NullString
}

// Configuration represents a V2Ray configuration stored in the database.
type Configuration struct {
	ID          int64
	UserID      int64
	Name        string
	Protocol    string
	ConfigData  string // Stored as a JSON string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Log represents a system or application log entry.
type Log struct {
	ID        int64
	Level     string
	Message   string
	Source    sql.NullString
	CreatedAt time.Time
}
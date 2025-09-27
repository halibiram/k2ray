package db

import (
	"database/sql"
	"time"
)

// User represents a user in the system.
// The struct tags for JSON are not included yet as this model is currently
// for database interaction only. API-specific models can be defined later if needed.
type User struct {
	ID                     int64
	Username               string
	PasswordHash           string
	TwoFactorSecret        sql.NullString
	TwoFactorEnabled       bool
	TwoFactorRecoveryCodes sql.NullString
}

// V2rayConfig represents a V2Ray configuration stored in the database.
type V2rayConfig struct {
	ID          int64
	UserID      int64
	Name        string
	Protocol    string
	ConfigData  string // Stored as a JSON string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

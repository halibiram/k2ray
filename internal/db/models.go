package db

// User represents a user in the system.
// The struct tags for JSON are not included yet as this model is currently
// for database interaction only. API-specific models can be defined later if needed.
type User struct {
	ID           int64
	Username     string
	PasswordHash string
}

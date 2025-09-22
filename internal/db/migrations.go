package db

import (
	"log"
)

// RunMigrations executes all necessary database migrations.
// In a larger application, this could be handled by a more robust migration library,
// but for this project, a simple function is sufficient.
func RunMigrations() {
	// SQL statement to create the 'users' table.
	// "IF NOT EXISTS" prevents an error if the table is already there.
	// The username is set to be UNIQUE to prevent duplicate user accounts.
	createUsersTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"username" TEXT NOT NULL UNIQUE,
		"password_hash" TEXT NOT NULL
	);`

	log.Println("Running database migration for 'users' table...")
	statement, err := DB.Prepare(createUsersTableSQL)
	if err != nil {
		log.Fatalf("Fatal error preparing users table migration: %v", err)
	}
	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		log.Fatalf("Fatal error executing users table migration: %v", err)
	}
	log.Println("'users' table migration successful.")
}

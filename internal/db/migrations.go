package db

import (
	"log"
)

// RunMigrations executes all necessary database migrations.
// In a larger application, this could be handled by a more robust migration library,
// but for this project, a simple function is sufficient.
func RunMigrations() {
	runMigration("users", `
		CREATE TABLE IF NOT EXISTS users (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"username" TEXT NOT NULL UNIQUE,
			"password_hash" TEXT NOT NULL
		);
	`)

	runMigration("revoked_tokens", `
		CREATE TABLE IF NOT EXISTS revoked_tokens (
			"jti" TEXT NOT NULL PRIMARY KEY,
			"expires_at" INTEGER NOT NULL
		);
	`)

	runMigration("revoked_tokens_index", `
		CREATE INDEX IF NOT EXISTS idx_revoked_tokens_expires_at ON revoked_tokens (expires_at);
	`)

	runMigration("v2ray_configs", `
		CREATE TABLE IF NOT EXISTS v2ray_configs (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"user_id" INTEGER NOT NULL,
			"name" TEXT NOT NULL,
			"protocol" TEXT NOT NULL,
			"config_data" TEXT NOT NULL,
			"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			"updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`)

	runMigration("system_settings", `
		CREATE TABLE IF NOT EXISTS system_settings (
			"key" TEXT NOT NULL PRIMARY KEY,
			"value" TEXT NOT NULL
		);
	`)
}

// runMigration is a helper function to execute a single migration statement.
func runMigration(name, sql string) {
	log.Printf("Running database migration for '%s'...", name)
	statement, err := DB.Prepare(sql)
	if err != nil {
		log.Fatalf("Fatal error preparing migration '%s': %v", name, err)
	}
	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		log.Fatalf("Fatal error executing migration '%s': %v", name, err)
	}
	log.Printf("Migration '%s' successful.", name)
}

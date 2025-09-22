package db

import (
	"k2ray/internal/utils"
	"log"
)

// SeedDatabase seeds the database with initial required data, such as a default admin user,
// but only if the data doesn't already exist.
func SeedDatabase() {
	var userCount int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		log.Fatalf("Fatal error querying user count for seeding: %v", err)
		return
	}

	// If no users exist in the database, create a default admin user.
	if userCount == 0 {
		log.Println("No users found in the database. Seeding with default admin user...")

		// Default credentials. These should be changed by the user immediately.
		defaultUsername := "admin"
		defaultPassword := "password"

		hashedPassword, err := utils.HashPassword(defaultPassword)
		if err != nil {
			log.Fatalf("Fatal error hashing default admin password: %v", err)
		}

		insertSQL := `INSERT INTO users (username, password_hash) VALUES (?, ?)`
		statement, err := DB.Prepare(insertSQL)
		if err != nil {
			log.Fatalf("Fatal error preparing insert statement for default user: %v", err)
		}
		defer statement.Close()

		_, err = statement.Exec(defaultUsername, hashedPassword)
		if err != nil {
			log.Fatalf("Fatal error executing insert statement for default user: %v", err)
		}

		log.Printf("Successfully created default admin user.")
		log.Printf("  Username: %s", defaultUsername)
		log.Printf("  Password: %s", defaultPassword)
		log.Println("IMPORTANT: It is strongly recommended to change this password after the first login.")
	}
}

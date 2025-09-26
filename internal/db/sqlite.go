package db

import (
	"database/sql"
	"k2ray/internal/config"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3" // The SQLite driver, imported for its side effects
)

var (
	// DB is the global database connection pool.
	DB *sql.DB
	// once ensures the database is initialized only once.
	once sync.Once
)

// InitDB initializes the database connection using the URL from the loaded configuration.
// It uses a sync.Once to ensure that the initialization code runs exactly once,
// making it safe to be called from multiple goroutines.
func InitDB() {
	once.Do(func() {
		var err error

		if config.AppConfig == nil {
			log.Fatal("Configuration is not loaded. Call config.LoadConfig() first.")
		}

		dbURL := config.AppConfig.DatabaseURL
		if dbURL == "" {
			log.Fatal("DATABASE_URL is not set in the configuration.")
		}

		// sql.Open just validates its arguments, it doesn't create a connection.
		DB, err = sql.Open("sqlite3", dbURL)
		if err != nil {
			log.Fatalf("Fatal error opening database connection: %v", err)
		}

		// DB.Ping() is used to verify that the connection to the database is alive.
		if err = DB.Ping(); err != nil {
			log.Fatalf("Fatal error connecting to database: %v", err)
		}

		log.Println("Database connection established successfully.")
	})
}

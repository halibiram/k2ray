package db

import (
	"database/sql"
	"errors"
	"k2ray/internal/config"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3" // The SQLite driver
	"github.com/rs/zerolog/log"
)

var (
	// DB is the global database connection pool.
	DB *sql.DB
	// once ensures the database is initialized only once.
	once sync.Once
)

// InitDB initializes the database connection and runs migrations.
// It uses a sync.Once to ensure that this process runs exactly once.
func InitDB() {
	once.Do(func() {
		var err error

		if config.AppConfig == nil {
			// This can happen in tests, so we load a default config.
			config.LoadConfig("")
		}

		dbURL := config.AppConfig.DatabaseURL
		if dbURL == "" {
			log.Fatal().Msg("DATABASE_URL is not set in the configuration.")
		}

		// Open the database connection.
		DB, err = sql.Open("sqlite3", dbURL)
		if err != nil {
			log.Fatal().Err(err).Msg("Fatal error opening database connection")
		}

		// Ping the database to verify the connection is alive.
		if err = DB.Ping(); err != nil {
			log.Fatal().Err(err).Msg("Fatal error connecting to database")
		}
		log.Info().Msg("Database connection established successfully.")

		// Run database migrations from the embedded filesystem.
		log.Info().Msg("Running database migrations...")
		if err := runMigrations(DB); err != nil {
			log.Fatal().Err(err).Msg("Fatal error running database migrations")
		}
		log.Info().Msg("Database migrations completed successfully.")
	})
}

// runMigrations applies all available "up" migrations using the embedded SQL files.
func runMigrations(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	// Use the embedded filesystem. The path is relative to the go:embed directive.
	source, err := iofs.New(MigrationsFS, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"iofs", // Specify the iofs source driver
		source,
		"sqlite3", // The database name
		driver,
	)
	if err != nil {
		return err
	}

	// Apply all available up migrations.
	// migrate.ErrNoChange is a normal condition, not an error.
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
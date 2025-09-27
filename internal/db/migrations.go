package db

import "embed"

// MigrationsFS holds the embedded SQL migration files.
// The "migrations" path is relative to this file's directory.
//go:embed migrations/*.sql
var MigrationsFS embed.FS
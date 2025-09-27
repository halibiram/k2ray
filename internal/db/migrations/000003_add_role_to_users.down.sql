-- Revert the changes from 000003_add_role_to_users.up.sql

-- SQLite does not directly support DROP COLUMN.
-- The common workaround is to create a new table without the column,
-- copy the data, and then replace the old table.

-- 1. Rename the existing table
ALTER TABLE users RENAME TO users_old;

-- 2. Create the new table without the 'role' column
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    two_factor_secret TEXT,
    two_factor_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    two_factor_recovery_codes TEXT
);

-- 3. Copy the data from the old table to the new one
INSERT INTO users (id, username, password_hash, two_factor_secret, two_factor_enabled, two_factor_recovery_codes)
SELECT id, username, password_hash, two_factor_secret, two_factor_enabled, two_factor_recovery_codes FROM users_old;

-- 4. Drop the old table
DROP TABLE users_old;
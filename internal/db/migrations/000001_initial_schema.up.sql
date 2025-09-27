-- +migrate Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE users (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "username" TEXT NOT NULL UNIQUE,
    "password_hash" TEXT NOT NULL,
    "two_factor_secret" TEXT,
    "two_factor_enabled" INTEGER NOT NULL DEFAULT 0,
    "two_factor_recovery_codes" TEXT
);

CREATE TABLE revoked_tokens (
    "jti" TEXT NOT NULL PRIMARY KEY,
    "expires_at" INTEGER NOT NULL
);

CREATE INDEX idx_revoked_tokens_expires_at ON revoked_tokens (expires_at);

CREATE TABLE configurations (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "user_id" INTEGER NOT NULL,
    "name" TEXT NOT NULL,
    "protocol" TEXT NOT NULL,
    "config_data" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE settings (
    "key" TEXT NOT NULL PRIMARY KEY,
    "value" TEXT NOT NULL
);

CREATE TABLE logs (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "level" TEXT NOT NULL,
    "message" TEXT NOT NULL,
    "source" TEXT,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
-- SQL in this section is executed when the migration is rolled back.
-- This section is intentionally left blank in the up file.
-- The rollback logic is in the down file.
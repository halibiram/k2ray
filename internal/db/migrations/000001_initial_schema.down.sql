-- +migrate Up
-- SQL in this section is executed when the migration is applied.
-- This section is intentionally left blank in the down file.
-- The application logic is in the up file.

-- +migrate Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS logs;
DROP TABLE IF EXISTS settings;
DROP TABLE IF EXISTS configurations;
DROP TABLE IF EXISTS revoked_tokens;
DROP TABLE IF EXISTS users;
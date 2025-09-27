-- +migrate Up
-- This migration seeds the database with a default admin user if one does not already exist.
-- The password is 'password', pre-hashed with bcrypt.
-- IMPORTANT: This password should be changed immediately after the first login.

INSERT INTO users (username, password_hash)
SELECT 'admin', '$2a$10$t.he/.I2.P93J9iQ25R12uBwJg22.E..e.brDbE3zJ/a2L0A5v35q'
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin');

-- +migrate Down
-- This section is intentionally left blank in the up file.
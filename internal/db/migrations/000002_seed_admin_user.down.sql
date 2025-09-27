-- +migrate Up
-- This section is intentionally left blank in the down file.

-- +migrate Down
-- This migration removes the default admin user.

DELETE FROM users WHERE username = 'admin';
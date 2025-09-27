-- +migrate Up
-- SQL in this section is executed when the migration is applied.

-- Index on user_id in the configurations table for faster lookups of a user's configurations.
CREATE INDEX idx_configurations_user_id ON configurations (user_id);

-- Index on created_at in the logs table for faster sorting and filtering of logs by time.
CREATE INDEX idx_logs_created_at ON logs (created_at);

-- Index on level in the logs table for faster filtering of logs by log level.
CREATE INDEX idx_logs_level ON logs (level);
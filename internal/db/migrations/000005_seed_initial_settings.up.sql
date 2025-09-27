-- +migrate Up
-- SQL in this section is executed when the migration is applied.

-- Seed initial application settings
INSERT INTO settings ("key", "value") VALUES
('app_name', 'K2Ray Panel'),
('default_language', 'en'),
('allow_registrations', 'true');
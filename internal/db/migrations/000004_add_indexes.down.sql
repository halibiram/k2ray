-- +migrate Down
-- SQL in this section is executed when the migration is rolled back.

DROP INDEX idx_configurations_user_id;
DROP INDEX idx_logs_created_at;
DROP INDEX idx_logs_level;
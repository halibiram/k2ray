-- +migrate Down
-- SQL in this section is executed when the migration is rolled back.

DELETE FROM settings WHERE "key" IN ('app_name', 'default_language', 'allow_registrations');
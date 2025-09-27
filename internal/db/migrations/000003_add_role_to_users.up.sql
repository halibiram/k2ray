-- Add a 'role' column to the 'users' table
ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'user';

-- Optionally, you could update existing users to have a default role.
-- For example, to make the first user an admin:
-- UPDATE users SET role = 'admin' WHERE id = 1;
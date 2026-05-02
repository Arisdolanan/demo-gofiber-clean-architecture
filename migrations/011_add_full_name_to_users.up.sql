-- Migration: Add full_name to users table
-- Version: 011

ALTER TABLE users ADD COLUMN full_name VARCHAR(255);

-- Update existing full_name from related tables if possible (optional but good for consistency)
-- This is a simple heuristic; in a real scenario, you might want more complex logic.
UPDATE users u
SET full_name = t.full_name
FROM teachers t
WHERE u.id = t.user_id AND u.full_name IS NULL;

UPDATE users u
SET full_name = s.full_name
FROM staff s
WHERE u.id = s.user_id AND u.full_name IS NULL;

-- Migration: Add phone_number and must_change_password to users
-- Version: 012

ALTER TABLE users ADD COLUMN phone_number VARCHAR(20);
ALTER TABLE users ADD COLUMN must_change_password BOOLEAN DEFAULT FALSE;

-- Update existing users: set phone_number from related tables if possible
-- (Heuristic for Staff)
UPDATE users u
SET phone_number = s.phone
FROM staff s
WHERE u.id = s.user_id AND u.phone_number IS NULL;

-- (Heuristic for Teachers)
UPDATE users u
SET phone_number = t.phone
FROM teachers t
WHERE u.id = t.user_id AND u.phone_number IS NULL;

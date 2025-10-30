-- Migration: Drop users table
-- Version: 001
-- Description: Rollback - Drop users table

-- Drop indexes first
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_created_by;
DROP INDEX IF EXISTS idx_users_updated_by;
DROP INDEX IF EXISTS idx_users_deleted_by;

-- Drop table
DROP TABLE IF EXISTS users;
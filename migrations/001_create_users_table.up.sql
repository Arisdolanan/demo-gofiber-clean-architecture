-- Migration: Create users table
-- Version: 001
-- Description: Create users table for authentication

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    updated_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    deleted_by BIGINT REFERENCES users(id) ON DELETE SET NULL
);

-- Create index for email lookup
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL;

-- Create index for username lookup
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username) WHERE deleted_at IS NULL;

-- Create index for soft delete
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Add indexes for audit trail fields
CREATE INDEX IF NOT EXISTS idx_users_created_by ON users(created_by);
CREATE INDEX IF NOT EXISTS idx_users_updated_by ON users(updated_by);
CREATE INDEX IF NOT EXISTS idx_users_deleted_by ON users(deleted_by);

-- Add comment to table
COMMENT ON TABLE users IS 'Users table for authentication and user management';
COMMENT ON COLUMN users.id IS 'Primary key, auto-incrementing';
COMMENT ON COLUMN users.username IS 'Unique username for user identification';
COMMENT ON COLUMN users.email IS 'Unique email address for user authentication';
COMMENT ON COLUMN users.password IS 'Hashed password using bcrypt';
COMMENT ON COLUMN users.created_at IS 'Timestamp when user was created';
COMMENT ON COLUMN users.updated_at IS 'Timestamp when user was last updated';
COMMENT ON COLUMN users.deleted_at IS 'Soft delete timestamp, NULL means active user';
COMMENT ON COLUMN users.created_by IS 'ID of the user who created this record';
COMMENT ON COLUMN users.updated_by IS 'ID of the user who last updated this record';
COMMENT ON COLUMN users.deleted_by IS 'ID of the user who deleted this record';
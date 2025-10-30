-- Migration: Drop refresh_tokens table
-- Version: 002
-- Description: Rollback - Drop refresh_tokens table

-- Drop indexes first
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_refresh_tokens_token_hash;
DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_refresh_tokens_is_revoked;

-- Drop table
DROP TABLE IF EXISTS refresh_tokens; 
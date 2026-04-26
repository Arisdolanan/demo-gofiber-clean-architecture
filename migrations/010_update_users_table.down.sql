-- Down Migration: Update Users Table

ALTER TABLE users DROP COLUMN IF EXISTS last_login_at;
ALTER TABLE users DROP COLUMN IF EXISTS is_active;
ALTER TABLE users DROP COLUMN IF EXISTS user_type;
DROP TYPE IF EXISTS user_type;
ALTER TABLE users DROP COLUMN IF EXISTS school_id;

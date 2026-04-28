-- Up Migration: Update Users Table

ALTER TABLE users 
ADD COLUMN IF NOT EXISTS school_id BIGINT REFERENCES schools(id) ON DELETE SET NULL;

CREATE TYPE user_type AS ENUM ('super_admin', 'admin', 'teacher', 'student', 'parent', 'staff');
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS user_type user_type DEFAULT 'student';

ALTER TABLE users 
ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true,
ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMP WITH TIME ZONE;

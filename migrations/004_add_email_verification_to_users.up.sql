-- Add email verification fields to users table
ALTER TABLE users ADD COLUMN email_verified_at TIMESTAMP DEFAULT NULL;
ALTER TABLE users ADD COLUMN email_verification_token VARCHAR(255) DEFAULT NULL;
ALTER TABLE users ADD COLUMN email_verification_expires_at TIMESTAMP DEFAULT NULL;
ALTER TABLE users ADD COLUMN password_reset_token VARCHAR(255) DEFAULT NULL;
ALTER TABLE users ADD COLUMN password_reset_expires_at TIMESTAMP DEFAULT NULL;

-- Add indexes for better performance
CREATE INDEX idx_users_email_verified_at ON users(email_verified_at);
CREATE INDEX idx_users_email_verification_token ON users(email_verification_token) WHERE email_verification_token IS NOT NULL;
CREATE INDEX idx_users_email_verification_expires_at ON users(email_verification_expires_at) WHERE email_verification_expires_at IS NOT NULL;
CREATE INDEX idx_users_password_reset_token ON users(password_reset_token) WHERE password_reset_token IS NOT NULL;
CREATE INDEX idx_users_password_reset_expires_at ON users(password_reset_expires_at) WHERE password_reset_expires_at IS NOT NULL;

-- Add comments for documentation
COMMENT ON COLUMN users.email_verified_at IS 'Timestamp when email was verified, NULL means not verified';
COMMENT ON COLUMN users.email_verification_token IS 'Token for email verification, NULL when verified';
COMMENT ON COLUMN users.email_verification_expires_at IS 'Expiry timestamp for email verification token';
COMMENT ON COLUMN users.password_reset_token IS 'Token for password reset, NULL when not requested';
COMMENT ON COLUMN users.password_reset_expires_at IS 'Expiry timestamp for password reset token';
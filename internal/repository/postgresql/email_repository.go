package postgresql

import (
	"database/sql"
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/jmoiron/sqlx"
)

// EmailRepository interface for email token operations
type EmailRepository interface {
	// Email verification token operations (stored in users table)
	SetVerificationToken(userID int64, token string, expiresAt time.Time) error
	GetUserByVerificationToken(token string) (*entity.User, error)
	ClearVerificationToken(userID int64) error

	// Password reset token operations (stored in users table)
	SetPasswordResetToken(userID int64, token string, expiresAt time.Time) error
	GetUserByPasswordResetToken(token string) (*entity.User, error)
	ClearPasswordResetToken(userID int64) error

	// User email verification status
	MarkEmailAsVerified(userID int64) error

	// Cleanup operations
	DeleteExpiredTokens() error
}

// emailRepository implements EmailRepository
type emailRepository struct {
	db *sqlx.DB
}

// NewEmailRepository creates a new email repository
func NewEmailRepository(db *sqlx.DB) EmailRepository {
	return &emailRepository{db: db}
}

// SetVerificationToken sets email verification token for a user
func (r *emailRepository) SetVerificationToken(userID int64, token string, expiresAt time.Time) error {
	query := `
		UPDATE users 
		SET email_verification_token = $1, email_verification_expires_at = $2, updated_at = $3
		WHERE id = $4`

	_, err := r.db.Exec(query, token, expiresAt, time.Now(), userID)
	return err
}

// GetUserByVerificationToken gets user by verification token
func (r *emailRepository) GetUserByVerificationToken(token string) (*entity.User, error) {
	query := `
		SELECT id, username, email, password, email_verified_at, 
		       email_verification_token, email_verification_expires_at,
		       password_reset_token, password_reset_expires_at,
		       created_at, updated_at, deleted_at
		FROM users 
		WHERE email_verification_token = $1 
		  AND email_verification_expires_at > $2
		  AND deleted_at IS NULL`

	var user entity.User
	err := r.db.Get(&user, query, token, time.Now())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// ClearVerificationToken clears email verification token for a user
func (r *emailRepository) ClearVerificationToken(userID int64) error {
	query := `
		UPDATE users 
		SET email_verification_token = NULL, email_verification_expires_at = NULL, updated_at = $1
		WHERE id = $2`

	_, err := r.db.Exec(query, time.Now(), userID)
	return err
}

// SetPasswordResetToken sets password reset token for a user
func (r *emailRepository) SetPasswordResetToken(userID int64, token string, expiresAt time.Time) error {
	query := `
		UPDATE users 
		SET password_reset_token = $1, password_reset_expires_at = $2, updated_at = $3
		WHERE id = $4`

	_, err := r.db.Exec(query, token, expiresAt, time.Now(), userID)
	return err
}

// GetUserByPasswordResetToken gets user by password reset token
func (r *emailRepository) GetUserByPasswordResetToken(token string) (*entity.User, error) {
	query := `
		SELECT id, username, email, password, email_verified_at, 
		       email_verification_token, email_verification_expires_at,
		       password_reset_token, password_reset_expires_at,
		       created_at, updated_at, deleted_at
		FROM users 
		WHERE password_reset_token = $1 
		  AND password_reset_expires_at > $2
		  AND deleted_at IS NULL`

	var user entity.User
	err := r.db.Get(&user, query, token, time.Now())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// ClearPasswordResetToken clears password reset token for a user
func (r *emailRepository) ClearPasswordResetToken(userID int64) error {
	query := `
		UPDATE users 
		SET password_reset_token = NULL, password_reset_expires_at = NULL, updated_at = $1
		WHERE id = $2`

	_, err := r.db.Exec(query, time.Now(), userID)
	return err
}

// MarkEmailAsVerified marks user's email as verified and clears verification token
func (r *emailRepository) MarkEmailAsVerified(userID int64) error {
	query := `
		UPDATE users 
		SET email_verified_at = $1, 
		    email_verification_token = NULL, 
		    email_verification_expires_at = NULL,
		    updated_at = $1
		WHERE id = $2`

	_, err := r.db.Exec(query, time.Now(), userID)
	return err
}

// DeleteExpiredTokens clears all expired tokens from users table
func (r *emailRepository) DeleteExpiredTokens() error {
	now := time.Now()

	// Clear expired email verification tokens
	query1 := `
		UPDATE users 
		SET email_verification_token = NULL, email_verification_expires_at = NULL, updated_at = $1
		WHERE email_verification_expires_at < $1`

	// Clear expired password reset tokens
	query2 := `
		UPDATE users 
		SET password_reset_token = NULL, password_reset_expires_at = NULL, updated_at = $1
		WHERE password_reset_expires_at < $1`

	if _, err := r.db.Exec(query1, now); err != nil {
		return err
	}

	if _, err := r.db.Exec(query2, now); err != nil {
		return err
	}

	return nil
}

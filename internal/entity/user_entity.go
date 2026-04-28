package entity

import (
	"time"
)

type UserType string

const (
	UserSuperAdmin UserType = "super_admin"
	UserAdmin      UserType = "admin"
	UserTeacher    UserType = "teacher"
	UserStudent    UserType = "student"
	UserParent     UserType = "parent"
)

type User struct {
	ID                         int64      `json:"id" db:"id"`
	SchoolID                   *int64     `json:"school_id,omitempty" db:"school_id"`
	Username                   string     `json:"username" db:"username"`
	Email                      string     `json:"email" db:"email"`
	Password                   string     `json:"password" db:"password"`
	UserType                   UserType   `json:"user_type" db:"user_type"`
	IsActive                   bool       `json:"is_active" db:"is_active"`
	LastLoginAt                *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt                  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                  time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                  *time.Time `json:"deleted_at" db:"deleted_at"`
	EmailVerifiedAt            *time.Time `json:"email_verified_at" db:"email_verified_at"`
	EmailVerificationToken     *string    `json:"email_verification_token,omitempty" db:"email_verification_token"`
	EmailVerificationExpiresAt *time.Time `json:"email_verification_expires_at,omitempty" db:"email_verification_expires_at"`
	PasswordResetToken         *string    `json:"password_reset_token,omitempty" db:"password_reset_token"`
	PasswordResetExpiresAt     *time.Time `json:"password_reset_expires_at,omitempty" db:"password_reset_expires_at"`
	CreatedBy                  *int64     `json:"created_by,omitempty" db:"created_by"`
	UpdatedBy                  *int64     `json:"updated_by,omitempty" db:"updated_by"`
	DeletedBy                  *int64     `json:"deleted_by,omitempty" db:"deleted_by"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) GetID() int64 {
	return u.ID
}

func (u *User) GetUsername() string {
	return u.Username
}

// IsEmailVerified checks if user's email is verified
func (u *User) IsEmailVerified() bool {
	return u.EmailVerifiedAt != nil
}

// GetEmail returns user email
func (u *User) GetEmail() string {
	return u.Email
}

// HasValidEmailVerificationToken checks if user has a valid email verification token
func (u *User) HasValidEmailVerificationToken() bool {
	return u.EmailVerificationToken != nil &&
		u.EmailVerificationExpiresAt != nil &&
		u.EmailVerificationExpiresAt.After(time.Now())
}

// HasValidPasswordResetToken checks if user has a valid password reset token
func (u *User) HasValidPasswordResetToken() bool {
	return u.PasswordResetToken != nil &&
		u.PasswordResetExpiresAt != nil &&
		u.PasswordResetExpiresAt.After(time.Now())
}

// ClearEmailVerificationToken clears the email verification token
func (u *User) ClearEmailVerificationToken() {
	u.EmailVerificationToken = nil
	u.EmailVerificationExpiresAt = nil
}

// ClearPasswordResetToken clears the password reset token
func (u *User) ClearPasswordResetToken() {
	u.PasswordResetToken = nil
	u.PasswordResetExpiresAt = nil
}

// UserListResponse represents the response for user list
type UserListResponse struct {
	Users      []User `json:"users"`
	TotalCount int64  `json:"total_count"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

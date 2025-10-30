package entity

// EmailVerificationRequest represents email verification request
type EmailVerificationRequest struct {
	Token string `json:"token" validate:"required"`
}

// PasswordResetRequest represents password reset request
type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// PasswordResetConfirmRequest represents password reset confirmation request
type PasswordResetConfirmRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// ResendVerificationRequest represents resend verification email request
type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// EmailTemplate represents email template data
type EmailTemplate struct {
	To       string
	Subject  string
	Template string
	Data     map[string]interface{}
}

// EmailVerificationData represents data for email verification template
type EmailVerificationData struct {
	UserEmail       string
	VerificationURL string
	AppName         string
	ExpirationHours int
}

// PasswordResetData represents data for password reset template
type PasswordResetData struct {
	UserEmail      string
	ResetURL       string
	AppName        string
	ExpirationMins int
}

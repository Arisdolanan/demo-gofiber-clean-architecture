package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/infrastructure/email"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/repository/postgresql"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/configuration"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/sirupsen/logrus"
)

// EmailUsecase interface for email operations
type EmailUsecase interface {
	// Email verification
	SendVerificationEmail(userID int64, userEmail string) error
	VerifyEmail(token string) error
	ResendVerificationEmail(userEmail string) error

	// Password reset
	SendPasswordResetEmail(userEmail string) error
	ResetPassword(token, newPassword string) error

	// Cleanup operations
	CleanupExpiredTokens() error
}

// emailUsecase implements EmailUsecase
type emailUsecase struct {
	emailRepo    postgresql.EmailRepository
	userRepo     postgresql.UserRepository
	emailService email.EmailService
	logger       *logrus.Logger
	emailConfig  configuration.EmailConfig
}

// NewEmailUsecase creates a new email use case
func NewEmailUsecase(
	emailRepo postgresql.EmailRepository,
	userRepo postgresql.UserRepository,
	emailService email.EmailService,
	logger *logrus.Logger,
) EmailUsecase {
	return &emailUsecase{
		emailRepo:    emailRepo,
		userRepo:     userRepo,
		emailService: emailService,
		logger:       logger,
		emailConfig:  configuration.GetEmailConfig(),
	}
}

// generateSecureToken generates a cryptographically secure random token
func (uc *emailUsecase) generateSecureToken() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// SendVerificationEmail sends an email verification email to the user
func (uc *emailUsecase) SendVerificationEmail(userID int64, userEmail string) error {
	// Generate secure token
	token, err := uc.generateSecureToken()
	if err != nil {
		uc.logger.Errorf("Failed to generate verification token: %v", err)
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Set verification token in users table
	expiresAt := time.Now().Add(time.Duration(uc.emailConfig.VerificationTokenExpiry) * time.Second)
	if err := uc.emailRepo.SetVerificationToken(userID, token, expiresAt); err != nil {
		uc.logger.Errorf("Failed to set verification token: %v", err)
		return fmt.Errorf("failed to set verification token: %w", err)
	}

	// Send verification email
	if err := uc.emailService.SendVerificationEmail(userEmail, token); err != nil {
		uc.logger.Errorf("Failed to send verification email to %s: %v", userEmail, err)
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	uc.logger.Infof("Verification email sent successfully to %s", userEmail)
	return nil
}

// VerifyEmail verifies a user's email using the provided token
func (uc *emailUsecase) VerifyEmail(token string) error {
	// Get user by verification token
	user, err := uc.emailRepo.GetUserByVerificationToken(token)
	if err != nil {
		uc.logger.Errorf("Failed to get user by verification token: %v", err)
		return fmt.Errorf("failed to get user by verification token: %w", err)
	}

	if user == nil {
		uc.logger.Warn("Invalid or expired verification token provided")
		return fmt.Errorf("invalid or expired verification token")
	}

	// Check if email is already verified
	if user.IsEmailVerified() {
		uc.logger.Warnf("Email is already verified for user ID: %d", user.ID)
		return fmt.Errorf("email is already verified")
	}

	// Mark email as verified and clear verification token
	if err := uc.emailRepo.MarkEmailAsVerified(user.ID); err != nil {
		uc.logger.Errorf("Failed to mark email as verified: %v", err)
		return fmt.Errorf("failed to mark email as verified: %w", err)
	}

	uc.logger.Infof("Email verified successfully for user ID: %d", user.ID)
	return nil
}

// ResendVerificationEmail resends verification email for a user
func (uc *emailUsecase) ResendVerificationEmail(userEmail string) error {
	// Get user by email
	user, err := uc.userRepo.FindByEmail(userEmail)
	if err != nil {
		uc.logger.Errorf("Failed to find user by email %s: %v", userEmail, err)
		return fmt.Errorf("user not found")
	}

	if user == nil {
		uc.logger.Warnf("User not found for email: %s", userEmail)
		return fmt.Errorf("user not found")
	}

	// Check if email is already verified
	if user.IsEmailVerified() {
		uc.logger.Warnf("Email is already verified for user: %s", userEmail)
		return fmt.Errorf("email is already verified")
	}

	// Send verification email
	return uc.SendVerificationEmail(user.ID, userEmail)
}

// SendPasswordResetEmail sends a password reset email to the user
func (uc *emailUsecase) SendPasswordResetEmail(userEmail string) error {
	// Get user by email
	user, err := uc.userRepo.FindByEmail(userEmail)
	if err != nil {
		uc.logger.Errorf("Failed to find user by email %s: %v", userEmail, err)
		return fmt.Errorf("user not found")
	}

	if user == nil {
		uc.logger.Warnf("User not found for email: %s", userEmail)
		return fmt.Errorf("user not found")
	}

	// Generate secure token
	token, err := uc.generateSecureToken()
	if err != nil {
		uc.logger.Errorf("Failed to generate reset token: %v", err)
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Set password reset token in users table
	expiresAt := time.Now().Add(time.Duration(uc.emailConfig.ResetTokenExpiry) * time.Second)
	if err := uc.emailRepo.SetPasswordResetToken(user.ID, token, expiresAt); err != nil {
		uc.logger.Errorf("Failed to set password reset token: %v", err)
		return fmt.Errorf("failed to set password reset token: %w", err)
	}

	// Send password reset email
	if err := uc.emailService.SendPasswordResetEmail(userEmail, token); err != nil {
		uc.logger.Errorf("Failed to send password reset email to %s: %v", userEmail, err)
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	uc.logger.Infof("Password reset email sent successfully to %s", userEmail)
	return nil
}

// ResetPassword resets a user's password using the provided token
func (uc *emailUsecase) ResetPassword(token, newPassword string) error {
	// Get user by password reset token
	user, err := uc.emailRepo.GetUserByPasswordResetToken(token)
	if err != nil {
		uc.logger.Errorf("Failed to get user by password reset token: %v", err)
		return fmt.Errorf("failed to get user by password reset token: %w", err)
	}

	if user == nil {
		uc.logger.Warn("Invalid or expired password reset token provided")
		return fmt.Errorf("invalid or expired password reset token")
	}

	// Validate new password complexity
	if err := utils.ValidatePasswordComplexity(newPassword); err != nil {
		uc.logger.Errorf("Password complexity validation failed: %v", err)
		return err
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		uc.logger.Errorf("Failed to hash new password: %v", err)
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update user's password
	user.Password = hashedPassword
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(user); err != nil {
		uc.logger.Errorf("Failed to update user password: %v", err)
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Clear password reset token
	if err := uc.emailRepo.ClearPasswordResetToken(user.ID); err != nil {
		uc.logger.Errorf("Failed to clear password reset token: %v", err)
		return fmt.Errorf("failed to clear password reset token: %w", err)
	}

	uc.logger.Infof("Password reset successfully for user ID: %d", user.ID)
	return nil
}

// CleanupExpiredTokens removes expired tokens from the database
func (uc *emailUsecase) CleanupExpiredTokens() error {
	// Delete expired tokens from users table
	if err := uc.emailRepo.DeleteExpiredTokens(); err != nil {
		uc.logger.Errorf("Failed to delete expired tokens: %v", err)
		return fmt.Errorf("failed to delete expired tokens: %w", err)
	}

	uc.logger.Info("Expired tokens cleanup completed successfully")
	return nil
}

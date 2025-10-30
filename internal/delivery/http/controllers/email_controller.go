package controllers

import (
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/usecase"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type EmailController struct {
	emailUsecase usecase.EmailUsecase
	validate     *validator.Validate
	log          *logrus.Logger
}

func NewEmailController(emailUsecase usecase.EmailUsecase, validate *validator.Validate, log *logrus.Logger) *EmailController {
	return &EmailController{
		emailUsecase: emailUsecase,
		validate:     validate,
		log:          log,
	}
}

// VerifyEmail handles email verification
// @Summary Verify user email
// @Description Verify user's email address using verification token
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} response.HTTPSuccessResponse "Email verified successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid or expired token"
// @Failure 500 {object} response.HTTPErrorResponse "Internal server error"
// @Router /api/v1/auth/verify-email [get]
func (c *EmailController) VerifyEmail(ctx *fiber.Ctx) error {
	token := ctx.Query("token")
	if token == "" {
		c.log.Error("Verification token is required")
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Verification token is required",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: "Token parameter is missing"},
			},
		})
	}

	// Verify email
	if err := c.emailUsecase.VerifyEmail(token); err != nil {
		c.log.Errorf("Email verification failed: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Email verification failed",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Email verified successfully",
		Data:    nil,
	})
}

// ResendVerificationEmail handles resending verification email
// @Summary Resend verification email
// @Description Resend email verification link to user's email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.ResendVerificationRequest true "User email"
// @Success 200 {object} response.HTTPSuccessResponse "Verification email sent successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid request or email already verified"
// @Failure 500 {object} response.HTTPErrorResponse "Internal server error"
// @Router /api/v1/auth/resend-verification [post]
func (c *EmailController) ResendVerificationEmail(ctx *fiber.Ctx) error {
	var req entity.ResendVerificationRequest

	if err := ctx.BodyParser(&req); err != nil {
		c.log.Errorf("Error parsing request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: err.Error()},
			},
		})
	}

	// Validate request
	if err := c.validate.Struct(req); err != nil {
		c.log.Errorf("Validation error: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: err.Error()},
			},
		})
	}

	// Resend verification email
	if err := c.emailUsecase.ResendVerificationEmail(req.Email); err != nil {
		c.log.Errorf("Failed to resend verification email: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Failed to resend verification email",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Verification email sent successfully",
		Data:    nil,
	})
}

// RequestPasswordReset handles password reset request
// @Summary Request password reset
// @Description Send password reset email to user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.PasswordResetRequest true "User email"
// @Success 200 {object} response.HTTPSuccessResponse "Password reset email sent successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid request"
// @Failure 500 {object} response.HTTPErrorResponse "Internal server error"
// @Router /api/v1/auth/forgot-password [post]
func (c *EmailController) RequestPasswordReset(ctx *fiber.Ctx) error {
	var req entity.PasswordResetRequest

	if err := ctx.BodyParser(&req); err != nil {
		c.log.Errorf("Error parsing request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: err.Error()},
			},
		})
	}

	// Validate request
	if err := c.validate.Struct(req); err != nil {
		c.log.Errorf("Validation error: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: err.Error()},
			},
		})
	}

	// Send password reset email
	if err := c.emailUsecase.SendPasswordResetEmail(req.Email); err != nil {
		c.log.Errorf("Failed to send password reset email: %v", err)
		// For security reasons, we don't reveal if the user exists or not
		// We always return success to prevent email enumeration attacks
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "If the email exists in our system, a password reset link has been sent",
		Data:    nil,
	})
}

// ResetPassword handles password reset with token
// @Summary Reset password
// @Description Reset user password using reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.PasswordResetConfirmRequest true "Reset token and new password"
// @Success 200 {object} response.HTTPSuccessResponse "Password reset successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid or expired token"
// @Failure 500 {object} response.HTTPErrorResponse "Internal server error"
// @Router /api/v1/auth/reset-password [post]
func (c *EmailController) ResetPassword(ctx *fiber.Ctx) error {
	var req entity.PasswordResetConfirmRequest

	if err := ctx.BodyParser(&req); err != nil {
		c.log.Errorf("Error parsing request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: err.Error()},
			},
		})
	}

	// Validate request
	if err := c.validate.Struct(req); err != nil {
		c.log.Errorf("Validation error: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: err.Error()},
			},
		})
	}

	// Reset password
	if err := c.emailUsecase.ResetPassword(req.Token, req.NewPassword); err != nil {
		c.log.Errorf("Password reset failed: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Password reset failed",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Password reset successfully",
		Data:    nil,
	})
}

// CleanupExpiredTokens handles cleanup of expired tokens (admin endpoint)
// @Summary Cleanup expired tokens
// @Description Remove expired email verification and password reset tokens
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.HTTPSuccessResponse "Cleanup completed successfully"
// @Failure 500 {object} response.HTTPErrorResponse "Internal server error"
// @Router /api/v1/admin/cleanup-tokens [post]
func (c *EmailController) CleanupExpiredTokens(ctx *fiber.Ctx) error {
	if err := c.emailUsecase.CleanupExpiredTokens(); err != nil {
		c.log.Errorf("Failed to cleanup expired tokens: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to cleanup expired tokens",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Cleanup completed successfully",
		Data:    nil,
	})
}

package controllers

import (
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/usecase"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	authUsecase usecase.AuthUsecase
	validate    *validator.Validate
	log         *logrus.Logger
}

func NewAuthController(authUsecase usecase.AuthUsecase, validate *validator.Validate, log *logrus.Logger) *AuthController {
	return &AuthController{
		authUsecase: authUsecase,
		validate:    validate,
		log:         log,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.RegisterRequest true "Registration information"
// @Success 201 {object} response.HTTPSuccessResponse "User registered successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid request body or validation failed"
// @Failure 500 {object} response.HTTPErrorResponse "Registration failed"
// @Router /api/v1/auth/register [post]
func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var req entity.RegisterRequest

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

	ipAddress := ctx.IP()
	userAgent := ctx.Get("User-Agent")

	if err := c.authUsecase.Register(req.Email, req.Password, req.SchoolID, req.UserType, ipAddress, userAgent); err != nil {
		c.log.Errorf("Registration error: %v", err)

		// Check if it's a password complexity error
		if pwdErr, ok := err.(*utils.PasswordComplexityError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
				Status:  fiber.StatusBadRequest,
				Message: "Password complexity requirements not met",
				Errors: func() []response.JSONError {
					errors := make([]response.JSONError, len(pwdErr.Errors))
					for i, errMsg := range pwdErr.Errors {
						errors[i] = response.JSONError{Status: fiber.StatusBadRequest, Message: errMsg}
					}
					return errors
				}(),
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Registration failed",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusCreated,
		Message: "User registered successfully",
		Data:    nil,
	})
}

// Login handles user authentication
// @Summary Login a user
// @Description Authenticate a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.LoginRequest true "Login credentials"
// @Success 200 {object} response.HTTPSuccessResponse{data=entity.AuthToken} "Login successful"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid request body or validation failed"
// @Failure 401 {object} response.HTTPErrorResponse "Login failed"
// @Router /api/v1/auth/login [post]
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req entity.LoginRequest

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

	ipAddress := ctx.IP()
	userAgent := ctx.Get("User-Agent")

	authToken, err := c.authUsecase.Login(req.Email, req.Password, ipAddress, userAgent)
	if err != nil {
		c.log.Errorf("Login error: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Login failed",
			Errors: []response.JSONError{
				{Status: fiber.StatusUnauthorized, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Login successful",
		Data:    authToken,
	})
}

// Verify handles token verification
// @Summary Verify user token
// @Description Verify the JWT token and return user information
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.HTTPSuccessResponse{data=entity.User} "Token verified successfully"
// @Failure 401 {object} response.HTTPErrorResponse "Unauthorized"
// @Router /api/v1/auth/verify [get]
func (c *AuthController) Verify(ctx *fiber.Ctx) error {
	// Get token from context (set by middleware)
	token := ctx.Locals("token").(string)

	user, err := c.authUsecase.VerifyToken(token)
	if err != nil {
		c.log.Errorf("Token verification error: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Errors: []response.JSONError{
				{Status: fiber.StatusUnauthorized, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Token verified successfully",
		Data:    user,
	})
}

// Refresh handles token refresh
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.HTTPSuccessResponse{data=entity.AuthToken} "Token refreshed successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid request body or validation failed"
// @Failure 401 {object} response.HTTPErrorResponse "Invalid refresh token"
// @Router /api/v1/auth/refresh [post]
func (c *AuthController) Refresh(ctx *fiber.Ctx) error {
	var req entity.RefreshTokenRequest

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

	authToken, err := c.authUsecase.RefreshToken(req.RefreshToken)
	if err != nil {
		c.log.Errorf("Refresh token error: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Token refresh failed",
			Errors: []response.JSONError{
				{Status: fiber.StatusUnauthorized, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Token refreshed successfully",
		Data:    authToken,
	})
}

// Logout handles user logout
// @Summary Logout user
// @Description Logout user and invalidate both access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.HTTPSuccessResponse "Logout successful"
// @Failure 401 {object} response.HTTPErrorResponse "Unauthorized"
// @Failure 500 {object} response.HTTPErrorResponse "Logout failed"
// @Router /api/v1/auth/logout [post]
func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	// Get user and token from context (set by middleware)
	userID := ctx.Locals("user_id").(int64)
	accessToken := ctx.Locals("token").(string)

	if err := c.authUsecase.Logout(userID, accessToken); err != nil {
		c.log.Errorf("Logout error: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Logout failed",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Logout successful",
		Data:    nil,
	})
}

// SwitchSchool handles switching the active school context
// @Summary Switch active school
// @Description Switch the active school context and get a new token
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body entity.SwitchSchoolRequest true "Target school ID"
// @Success 200 {object} response.HTTPSuccessResponse{data=entity.AuthToken} "School switched successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid request body"
// @Failure 401 {object} response.HTTPErrorResponse "Unauthorized"
// @Failure 403 {object} response.HTTPErrorResponse "Access denied"
// @Router /api/v1/auth/switch-school [post]
func (c *AuthController) SwitchSchool(ctx *fiber.Ctx) error {
	var req entity.SwitchSchoolRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Errors: []response.JSONError{{Status: fiber.StatusBadRequest, Message: err.Error()}},
		})
	}

	userID := ctx.Locals("user_id").(int64)
	ipAddress := ctx.IP()
	userAgent := ctx.Get("User-Agent")

	authToken, err := c.authUsecase.SwitchSchool(ctx.Context(), userID, req.SchoolID, ipAddress, userAgent)
	if err != nil {
		c.log.Errorf("Switch school error: %v", err)
		return ctx.Status(fiber.StatusForbidden).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusForbidden,
			Message: "Failed to switch school",
			Errors: []response.JSONError{{Status: fiber.StatusForbidden, Message: err.Error()}},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "School switched successfully",
		Data:    authToken,
	})
}

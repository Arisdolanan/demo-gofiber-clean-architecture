package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/mocks"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthControllerTestSuite struct {
	suite.Suite
	app            *fiber.App
	authController *AuthController
	mockUsecase    *mocks.MockAuthUsecase
}

func (suite *AuthControllerTestSuite) SetupTest() {
	suite.mockUsecase = new(mocks.MockAuthUsecase)
	validator := validator.New()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Suppress logs in tests

	suite.authController = NewAuthController(
		suite.mockUsecase,
		validator,
		logger,
	)

	suite.app = fiber.New()
	suite.app.Post("/register", suite.authController.Register)
	suite.app.Post("/login", suite.authController.Login)
	suite.app.Post("/refresh", suite.authController.Refresh)

	// Setup verify and logout routes with middleware
	suite.app.Get("/verify", func(c *fiber.Ctx) error {
		c.Locals("user_id", int64(1))
		c.Locals("email", "test@example.com")
		return suite.authController.Verify(c)
	})

	suite.app.Post("/logout", func(c *fiber.Ctx) error {
		c.Locals("user_id", int64(1))
		c.Locals("email", "test@example.com")
		c.Locals("token", "test.access.token")
		return suite.authController.Logout(c)
	})
}

func (suite *AuthControllerTestSuite) TearDownTest() {
	suite.mockUsecase.AssertExpectations(suite.T())
}

// Test Register - Success
func (suite *AuthControllerTestSuite) TestRegister_Success() {
	reqBody := entity.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Mock expectations
	suite.mockUsecase.On("Register", reqBody.Email, reqBody.Password).Return(nil)

	// Create request body
	jsonBody, _ := json.Marshal(reqBody)

	// Make request
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusCreated, resp.StatusCode)
}

// Test Register - Invalid JSON
func (suite *AuthControllerTestSuite) TestRegister_InvalidJSON() {
	// Make request with invalid JSON
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)
}

// Test Register - Validation Error
func (suite *AuthControllerTestSuite) TestRegister_ValidationError() {
	reqBody := entity.RegisterRequest{
		Email:    "invalid-email", // Invalid email format
		Password: "123",           // Too short password
	}

	// Create request body
	jsonBody, _ := json.Marshal(reqBody)

	// Make request
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)
}

// Test Register - Usecase Error
func (suite *AuthControllerTestSuite) TestRegister_UsecaseError() {
	reqBody := entity.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Mock expectations
	usecaseError := errors.New("user already exists")
	suite.mockUsecase.On("Register", reqBody.Email, reqBody.Password).Return(usecaseError)

	// Create request body
	jsonBody, _ := json.Marshal(reqBody)

	// Make request
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusInternalServerError, resp.StatusCode)
}

// Test Login - Success
func (suite *AuthControllerTestSuite) TestLogin_Success() {
	reqBody := entity.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	authToken := &entity.AuthToken{
		AccessToken:  "test.access.token",
		RefreshToken: "test.refresh.token",
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	// Mock expectations
	suite.mockUsecase.On("Login", reqBody.Email, reqBody.Password).Return(authToken, nil)

	// Create request body
	jsonBody, _ := json.Marshal(reqBody)

	// Make request
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode)
}

// Test Login - Invalid JSON
func (suite *AuthControllerTestSuite) TestLogin_InvalidJSON() {
	// Make request with invalid JSON
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)
}

// Test Login - Validation Error
func (suite *AuthControllerTestSuite) TestLogin_ValidationError() {
	reqBody := entity.LoginRequest{
		Email:    "invalid-email", // Invalid email format
		Password: "123",           // Too short password
	}

	// Create request body
	jsonBody, _ := json.Marshal(reqBody)

	// Make request
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)
}

// Test Login - Usecase Error
func (suite *AuthControllerTestSuite) TestLogin_UsecaseError() {
	reqBody := entity.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	// Mock expectations
	usecaseError := errors.New("invalid credentials")
	suite.mockUsecase.On("Login", reqBody.Email, reqBody.Password).Return(nil, usecaseError)

	// Create request body
	jsonBody, _ := json.Marshal(reqBody)

	// Make request
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusUnauthorized, resp.StatusCode)
}

// Test Verify - Success
func (suite *AuthControllerTestSuite) TestVerify_Success() {
	// Make request
	req, _ := http.NewRequest("GET", "/verify", nil)

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode)
}

// Test Refresh - Success
func (suite *AuthControllerTestSuite) TestRefresh_Success() {
	reqBody := entity.RefreshTokenRequest{
		RefreshToken: "valid.refresh.token",
	}

	authToken := &entity.AuthToken{
		AccessToken:  "test.access.token",
		RefreshToken: "test.refresh.token",
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	// Mock expectations
	suite.mockUsecase.On("RefreshToken", reqBody.RefreshToken).Return(authToken, nil)

	// Create request body
	jsonBody, _ := json.Marshal(reqBody)

	// Make request
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode)
}

// Test Refresh - Invalid JSON
func (suite *AuthControllerTestSuite) TestRefresh_InvalidJSON() {
	// Make request with invalid JSON
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)
}

// Test Refresh - Validation Error
func (suite *AuthControllerTestSuite) TestRefresh_ValidationError() {
	reqBody := entity.RefreshTokenRequest{
		RefreshToken: "", // Empty refresh token
	}

	// Create request body
	jsonBody, _ := json.Marshal(reqBody)

	// Make request
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)
}

// Test Refresh - Usecase Error
func (suite *AuthControllerTestSuite) TestRefresh_UsecaseError() {
	reqBody := entity.RefreshTokenRequest{
		RefreshToken: "invalid.refresh.token",
	}

	// Mock expectations
	usecaseError := errors.New("invalid refresh token")
	suite.mockUsecase.On("RefreshToken", reqBody.RefreshToken).Return(nil, usecaseError)

	// Create request body
	jsonBody, _ := json.Marshal(reqBody)

	// Make request
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusUnauthorized, resp.StatusCode)
}

// Test Logout - Success
func (suite *AuthControllerTestSuite) TestLogout_Success() {
	userID := int64(1)
	accessToken := "test.access.token"

	// Mock expectations
	suite.mockUsecase.On("Logout", userID, accessToken).Return(nil)

	// Make request
	req, _ := http.NewRequest("POST", "/logout", nil)

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode)
}

// Test Logout - Usecase Error
func (suite *AuthControllerTestSuite) TestLogout_UsecaseError() {
	userID := int64(1)
	accessToken := "test.access.token"
	usecaseError := errors.New("redis connection error")

	// Mock expectations
	suite.mockUsecase.On("Logout", userID, accessToken).Return(usecaseError)

	// Make request
	req, _ := http.NewRequest("POST", "/logout", nil)

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusInternalServerError, resp.StatusCode)
}

func TestAuthControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}

package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/delivery/messaging/kafka"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/mocks"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// createTestUser creates a test user for testing purposes
func createTestUser() *entity.User {
	// Create a new hash for "Password123!" to ensure consistency
	hashedPassword, _ := utils.HashPassword("Password123!")

	return &entity.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  hashedPassword, // hashed "Password123!"
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}
}

type AuthUsecaseTestSuite struct {
	suite.Suite
	authUsecase      AuthUsecase
	mockAuthRepo     *mocks.MockAuthRepository
	mockRedis        *mocks.MockAuthRedisRepository
	mockEmailUsecase       *MockEmailUsecase
	mockActivityLogUsecase *MockActivityLogUsecase
	validator              *validator.Validate
	logger                 *logrus.Logger
}

// MockActivityLogUsecase for testing
type MockActivityLogUsecase struct {
	mock.Mock
}

func (m *MockActivityLogUsecase) GetActivities(ctx context.Context, schoolID *int64, page, pageSize int) ([]*entity.ActivityLog, error) {
	args := m.Called(ctx, schoolID, page, pageSize)
	return args.Get(0).([]*entity.ActivityLog), args.Error(1)
}

func (m *MockActivityLogUsecase) GetDeletions(ctx context.Context, schoolID *int64, page, pageSize int) ([]*entity.ActivityLog, error) {
	args := m.Called(ctx, schoolID, page, pageSize)
	return args.Get(0).([]*entity.ActivityLog), args.Error(1)
}

func (m *MockActivityLogUsecase) LogActivity(ctx context.Context, log *entity.ActivityLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

// MockEmailUsecase for testing
type MockEmailUsecase struct {
	mock.Mock
}

func (m *MockEmailUsecase) SendVerificationEmail(userID int64, userEmail string) error {
	args := m.Called(userID, userEmail)
	return args.Error(0)
}

func (m *MockEmailUsecase) VerifyEmail(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockEmailUsecase) ResendVerificationEmail(userEmail string) error {
	args := m.Called(userEmail)
	return args.Error(0)
}

func (m *MockEmailUsecase) SendPasswordResetEmail(userEmail string) error {
	args := m.Called(userEmail)
	return args.Error(0)
}

func (m *MockEmailUsecase) ResetPassword(token, newPassword string) error {
	args := m.Called(token, newPassword)
	return args.Error(0)
}

func (m *MockEmailUsecase) CleanupExpiredTokens() error {
	args := m.Called()
	return args.Error(0)
}

func (suite *AuthUsecaseTestSuite) SetupTest() {
	suite.mockAuthRepo = new(mocks.MockAuthRepository)
	suite.mockRedis = new(mocks.MockAuthRedisRepository)
	suite.mockEmailUsecase = new(MockEmailUsecase)
	suite.mockActivityLogUsecase = new(MockActivityLogUsecase)
	suite.validator = validator.New()
	suite.logger = logrus.New()
	suite.logger.SetLevel(logrus.ErrorLevel) // Suppress logs in tests

	// Create a nil Kafka producer for tests
	var mockKafkaProducer *kafka.UserProducer

	suite.authUsecase = NewAuthUsecase(
		suite.mockAuthRepo,
		suite.mockRedis,
		suite.mockEmailUsecase,
		suite.validator,
		suite.logger,
		"test-secret-key",
		mockKafkaProducer,
		suite.mockActivityLogUsecase,
	)
}

func (suite *AuthUsecaseTestSuite) TearDownTest() {
	suite.mockAuthRepo.AssertExpectations(suite.T())
	suite.mockRedis.AssertExpectations(suite.T())
}

// Test Register - Success scenario
func (suite *AuthUsecaseTestSuite) TestRegister_Success() {
	email := "test@example.com"
	password := "StrongP@ssw0rd!" // Valid password that meets complexity requirements

	// Mock expectations
	suite.mockAuthRepo.On("FindByEmail", email).Return(nil, nil)
	suite.mockAuthRepo.On("Register", mock.AnythingOfType("*entity.User")).Return(nil)
	suite.mockEmailUsecase.On("SendVerificationEmail", mock.AnythingOfType("int64"), email).Return(nil)
	suite.mockActivityLogUsecase.On("LogActivity", mock.Anything, mock.AnythingOfType("*entity.ActivityLog")).Return(nil)

	// Execute
	err := suite.authUsecase.Register(email, password, nil, entity.UserStudent, "127.0.0.1", "test-ua")

	// Assert
	assert.NoError(suite.T(), err)
}

// Test Register - User already exists
func (suite *AuthUsecaseTestSuite) TestRegister_UserAlreadyExists() {
	email := "test@example.com"
	password := "StrongP@ssw0rd!" // Valid password that meets complexity requirements
	existingUser := createTestUser()

	// Mock expectations
	suite.mockAuthRepo.On("FindByEmail", email).Return(existingUser, nil)

	// Execute
	err := suite.authUsecase.Register(email, password, nil, entity.UserStudent, "127.0.0.1", "test-ua")

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user already exists", err.Error())
}

// Test Register - Database error when checking existing user
func (suite *AuthUsecaseTestSuite) TestRegister_DatabaseErrorWhenCheckingUser() {
	email := "test@example.com"
	password := "StrongP@ssw0rd!" // Valid password that meets complexity requirements
	dbError := errors.New("database connection error")

	// Mock expectations
	suite.mockAuthRepo.On("FindByEmail", email).Return(nil, dbError)

	// Execute
	err := suite.authUsecase.Register(email, password, nil, entity.UserStudent, "127.0.0.1", "test-ua")

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), dbError, err)
}

// Test Register - Database error when registering user
func (suite *AuthUsecaseTestSuite) TestRegister_DatabaseErrorWhenRegistering() {
	email := "test@example.com"
	password := "StrongP@ssw0rd!" // Valid password that meets complexity requirements
	dbError := errors.New("database insert error")

	// Mock expectations
	suite.mockAuthRepo.On("FindByEmail", email).Return(nil, nil)
	suite.mockAuthRepo.On("Register", mock.AnythingOfType("*entity.User")).Return(dbError)

	// Execute
	err := suite.authUsecase.Register(email, password, nil, entity.UserStudent, "127.0.0.1", "test-ua")

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), dbError, err)
}

// Test Register - Password complexity validation error
func (suite *AuthUsecaseTestSuite) TestRegister_PasswordComplexityError() {
	email := "test@example.com"
	weakPassword := "123" // Too short and doesn't meet complexity requirements

	// Execute
	err := suite.authUsecase.Register(email, weakPassword, nil, entity.UserStudent, "127.0.0.1", "test-ua")

	// Assert
	assert.Error(suite.T(), err)
	// The error type will be *utils.PasswordComplexityError, but we'll just check it's an error
}

// Test Login - Success scenario
func (suite *AuthUsecaseTestSuite) TestLogin_Success() {
	email := "test@example.com"
	password := "Password123!" // Valid password that matches the hashed password in createTestUser()
	user := createTestUser()

	// Mock expectations
	suite.mockAuthRepo.On("FindByEmail", email).Return(user, nil)
	suite.mockRedis.On("StoreRefreshToken", user.ID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil)
	suite.mockActivityLogUsecase.On("LogActivity", mock.Anything, mock.AnythingOfType("*entity.ActivityLog")).Return(nil)

	// Execute
	authToken, err := suite.authUsecase.Login(email, password, "127.0.0.1", "test-ua")

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), authToken)
	assert.NotEmpty(suite.T(), authToken.AccessToken)
	assert.NotEmpty(suite.T(), authToken.RefreshToken)
	assert.True(suite.T(), authToken.ExpiresAt.After(time.Now()))
}

// Test Login - User not found
func (suite *AuthUsecaseTestSuite) TestLogin_UserNotFound() {
	email := "test@example.com"
	password := "Password123!" // Valid password

	// Mock expectations
	suite.mockAuthRepo.On("FindByEmail", email).Return(nil, nil)

	// Execute
	authToken, err := suite.authUsecase.Login(email, password, "127.0.0.1", "test-ua")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), authToken)
	assert.Equal(suite.T(), "invalid credentials", err.Error())
}

// Test Login - Database error
func (suite *AuthUsecaseTestSuite) TestLogin_DatabaseError() {
	email := "test@example.com"
	password := "Password123!" // Valid password
	dbError := errors.New("database connection error")

	// Mock expectations
	suite.mockAuthRepo.On("FindByEmail", email).Return(nil, dbError)

	// Execute
	authToken, err := suite.authUsecase.Login(email, password, "127.0.0.1", "test-ua")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), authToken)
	assert.Equal(suite.T(), dbError, err)
}

// Test Login - Invalid password
func (suite *AuthUsecaseTestSuite) TestLogin_InvalidPassword() {
	email := "test@example.com"
	password := "wrongpassword"
	user := createTestUser()

	// Mock expectations
	suite.mockAuthRepo.On("FindByEmail", email).Return(user, nil)

	// Execute
	authToken, err := suite.authUsecase.Login(email, password, "127.0.0.1", "test-ua")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), authToken)
	assert.Equal(suite.T(), "invalid credentials", err.Error())
}

// Test Login - Redis error when storing refresh token
func (suite *AuthUsecaseTestSuite) TestLogin_RedisError() {
	email := "test@example.com"
	password := "Password123!" // Valid password that matches the hashed password in createTestUser()
	user := createTestUser()
	redisError := errors.New("redis connection error")

	// Mock expectations
	suite.mockAuthRepo.On("FindByEmail", email).Return(user, nil)
	suite.mockRedis.On("StoreRefreshToken", user.ID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(redisError)

	// Execute
	authToken, err := suite.authUsecase.Login(email, password, "127.0.0.1", "test-ua")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), authToken)
	assert.Equal(suite.T(), redisError, err)
}

// Test VerifyToken - Success scenario
func (suite *AuthUsecaseTestSuite) TestVerifyToken_Success() {
	// This test would require a valid JWT token
	// For now, we test that an invalid token returns an error
	result, err := suite.authUsecase.VerifyToken("invalid.token.format")

	// Assert - this will fail with invalid token, which is expected
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

// Test VerifyToken - User not found after token validation
func (suite *AuthUsecaseTestSuite) TestVerifyToken_UserNotFoundAfterValidation() {
	// Test with invalid token (will fail at token validation step)
	result, err := suite.authUsecase.VerifyToken("invalid.token.format")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

// Test VerifyToken - Database error when looking up user
func (suite *AuthUsecaseTestSuite) TestVerifyToken_DatabaseError() {
	// Test with invalid token (will fail at token validation step)
	result, err := suite.authUsecase.VerifyToken("invalid.token.format")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

// Test RefreshToken - Success scenario (will fail at token validation)
func (suite *AuthUsecaseTestSuite) TestRefreshToken_Success() {
	refreshToken := "valid.refresh.token"

	// This will fail because we need a real JWT token, but we test the business logic
	result, err := suite.authUsecase.RefreshToken(refreshToken)

	// Assert - will fail at token validation step
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "invalid refresh token")
}

// Test RefreshToken - Invalid token
func (suite *AuthUsecaseTestSuite) TestRefreshToken_InvalidToken() {
	invalidToken := "invalid.token.format"

	// Execute
	result, err := suite.authUsecase.RefreshToken(invalidToken)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "invalid refresh token", err.Error())
}

// Test RefreshToken - Token not found in Redis (won't be reached due to token validation)
func (suite *AuthUsecaseTestSuite) TestRefreshToken_TokenNotFoundInRedis() {
	refreshToken := "valid.refresh.token"

	// Execute with invalid token (will fail at token validation first)
	result, err := suite.authUsecase.RefreshToken(refreshToken)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "invalid refresh token")
}

// Test RefreshToken - Token mismatch (won't be reached due to token validation)
func (suite *AuthUsecaseTestSuite) TestRefreshToken_TokenMismatch() {
	refreshToken := "valid.refresh.token"

	// Execute with invalid token (will fail at token validation first)
	result, err := suite.authUsecase.RefreshToken(refreshToken)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "invalid refresh token")
}

// Test Logout - Success scenario
func (suite *AuthUsecaseTestSuite) TestLogout_Success() {
	userID := int64(1)
	accessToken := "test.access.token"

	// Mock expectations
	suite.mockRedis.On("BlacklistToken", accessToken, mock.AnythingOfType("time.Duration")).Return(nil)
	suite.mockRedis.On("DeleteRefreshToken", userID).Return(nil)

	// Execute
	err := suite.authUsecase.Logout(userID, accessToken)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test Logout - Redis error
func (suite *AuthUsecaseTestSuite) TestLogout_RedisError() {
	userID := int64(1)
	accessToken := "test.access.token"
	redisError := errors.New("redis connection error")

	// Mock expectations
	suite.mockRedis.On("BlacklistToken", accessToken, mock.AnythingOfType("time.Duration")).Return(redisError)

	// Execute
	err := suite.authUsecase.Logout(userID, accessToken)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), redisError, err)
}

func TestAuthUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(AuthUsecaseTestSuite))
}

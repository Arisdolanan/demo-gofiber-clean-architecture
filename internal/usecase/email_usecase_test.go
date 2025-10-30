package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockEmailRepository for testing
type MockEmailRepository struct {
	mock.Mock
}

func (m *MockEmailRepository) SetVerificationToken(userID int64, token string, expiresAt time.Time) error {
	args := m.Called(userID, token, expiresAt)
	return args.Error(0)
}

func (m *MockEmailRepository) GetUserByVerificationToken(token string) (*entity.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockEmailRepository) ClearVerificationToken(userID int64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockEmailRepository) SetPasswordResetToken(userID int64, token string, expiresAt time.Time) error {
	args := m.Called(userID, token, expiresAt)
	return args.Error(0)
}

func (m *MockEmailRepository) GetUserByPasswordResetToken(token string) (*entity.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockEmailRepository) ClearPasswordResetToken(userID int64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockEmailRepository) MarkEmailAsVerified(userID int64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockEmailRepository) DeleteExpiredTokens() error {
	args := m.Called()
	return args.Error(0)
}

// MockUserRepository for testing
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id int64) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// Additional methods to satisfy the interface
func (m *MockUserRepository) Create(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByUsername(username string) (*entity.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(limit, offset int) ([]*entity.User, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) SoftDelete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockEmailService for testing
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendVerificationEmail(userEmail, token string) error {
	args := m.Called(userEmail, token)
	return args.Error(0)
}

func (m *MockEmailService) SendPasswordResetEmail(userEmail, token string) error {
	args := m.Called(userEmail, token)
	return args.Error(0)
}

func (m *MockEmailService) SendEmail(emailTemplate *entity.EmailTemplate) error {
	args := m.Called(emailTemplate)
	return args.Error(0)
}

// EmailUsecaseTestSuite is the test suite for email use case
type EmailUsecaseTestSuite struct {
	suite.Suite
	emailUsecase     EmailUsecase
	mockEmailRepo    *MockEmailRepository
	mockUserRepo     *MockUserRepository
	mockEmailService *MockEmailService
	logger           *logrus.Logger
}

func (suite *EmailUsecaseTestSuite) SetupTest() {
	suite.mockEmailRepo = new(MockEmailRepository)
	suite.mockUserRepo = new(MockUserRepository)
	suite.mockEmailService = new(MockEmailService)
	suite.logger = logrus.New()
	suite.logger.SetLevel(logrus.ErrorLevel) // Suppress logs in tests

	suite.emailUsecase = NewEmailUsecase(
		suite.mockEmailRepo,
		suite.mockUserRepo,
		suite.mockEmailService,
		suite.logger,
	)
}

func (suite *EmailUsecaseTestSuite) TestSendVerificationEmail_Success() {
	userID := int64(1)
	userEmail := "test@example.com"

	// Mock expectations
	suite.mockEmailRepo.On("SetVerificationToken", userID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)
	suite.mockEmailService.On("SendVerificationEmail", userEmail, mock.AnythingOfType("string")).Return(nil)

	// Execute
	err := suite.emailUsecase.SendVerificationEmail(userID, userEmail)

	// Assert
	assert.NoError(suite.T(), err)
	suite.mockEmailRepo.AssertExpectations(suite.T())
	suite.mockEmailService.AssertExpectations(suite.T())
}

func (suite *EmailUsecaseTestSuite) TestSendVerificationEmail_SetTokenError() {
	userID := int64(1)
	userEmail := "test@example.com"

	// Mock expectations
	suite.mockEmailRepo.On("SetVerificationToken", userID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(errors.New("database error"))

	// Execute
	err := suite.emailUsecase.SendVerificationEmail(userID, userEmail)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to set verification token")
	suite.mockEmailRepo.AssertExpectations(suite.T())
}

func (suite *EmailUsecaseTestSuite) TestSendVerificationEmail_SendEmailError() {
	userID := int64(1)
	userEmail := "test@example.com"

	// Mock expectations
	suite.mockEmailRepo.On("SetVerificationToken", userID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)
	suite.mockEmailService.On("SendVerificationEmail", userEmail, mock.AnythingOfType("string")).Return(errors.New("email service error"))

	// Execute
	err := suite.emailUsecase.SendVerificationEmail(userID, userEmail)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to send verification email")
	suite.mockEmailRepo.AssertExpectations(suite.T())
	suite.mockEmailService.AssertExpectations(suite.T())
}

func (suite *EmailUsecaseTestSuite) TestVerifyEmail_Success() {
	token := "valid-token"
	userID := int64(1)
	now := time.Now()

	user := &entity.User{
		ID:                         userID,
		Email:                      "test@example.com",
		EmailVerificationToken:     &token,
		EmailVerificationExpiresAt: &now,
		EmailVerifiedAt:            nil, // Not verified yet
	}

	// Mock expectations
	suite.mockEmailRepo.On("GetUserByVerificationToken", token).Return(user, nil)
	suite.mockEmailRepo.On("MarkEmailAsVerified", userID).Return(nil)

	// Execute
	err := suite.emailUsecase.VerifyEmail(token)

	// Assert
	assert.NoError(suite.T(), err)
	suite.mockEmailRepo.AssertExpectations(suite.T())
}

func (suite *EmailUsecaseTestSuite) TestVerifyEmail_AlreadyVerified() {
	token := "valid-token"
	userID := int64(1)
	now := time.Now()
	verifiedAt := time.Now()

	user := &entity.User{
		ID:                         userID,
		Email:                      "test@example.com",
		EmailVerificationToken:     &token,
		EmailVerificationExpiresAt: &now,
		EmailVerifiedAt:            &verifiedAt, // Already verified
	}

	// Mock expectations
	suite.mockEmailRepo.On("GetUserByVerificationToken", token).Return(user, nil)

	// Execute
	err := suite.emailUsecase.VerifyEmail(token)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "email is already verified")
	suite.mockEmailRepo.AssertExpectations(suite.T())
}

func (suite *EmailUsecaseTestSuite) TestVerifyEmail_GetUserError() {
	token := "invalid-token"

	// Mock expectations
	suite.mockEmailRepo.On("GetUserByVerificationToken", token).Return(nil, errors.New("database error"))

	// Execute
	err := suite.emailUsecase.VerifyEmail(token)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to get user by verification token")
	suite.mockEmailRepo.AssertExpectations(suite.T())
}

func (suite *EmailUsecaseTestSuite) TestVerifyEmail_InvalidToken() {
	token := "invalid-token"

	// Mock expectations
	suite.mockEmailRepo.On("GetUserByVerificationToken", token).Return(nil, nil)

	// Execute
	err := suite.emailUsecase.VerifyEmail(token)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid or expired verification token")
	suite.mockEmailRepo.AssertExpectations(suite.T())
}

func (suite *EmailUsecaseTestSuite) TestResendVerificationEmail_Success() {
	userEmail := "test@example.com"
	userID := int64(1)

	user := &entity.User{
		ID:              userID,
		Email:           userEmail,
		EmailVerifiedAt: nil, // Not verified yet
	}

	// Mock expectations
	suite.mockUserRepo.On("FindByEmail", userEmail).Return(user, nil)
	suite.mockEmailRepo.On("SetVerificationToken", userID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)
	suite.mockEmailService.On("SendVerificationEmail", userEmail, mock.AnythingOfType("string")).Return(nil)

	// Execute
	err := suite.emailUsecase.ResendVerificationEmail(userEmail)

	// Assert
	assert.NoError(suite.T(), err)
	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockEmailRepo.AssertExpectations(suite.T())
	suite.mockEmailService.AssertExpectations(suite.T())
}

func (suite *EmailUsecaseTestSuite) TestResendVerificationEmail_UserNotFound() {
	userEmail := "nonexistent@example.com"

	// Mock expectations
	suite.mockUserRepo.On("FindByEmail", userEmail).Return(nil, nil)

	// Execute
	err := suite.emailUsecase.ResendVerificationEmail(userEmail)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "user not found")
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *EmailUsecaseTestSuite) TestResendVerificationEmail_AlreadyVerified() {
	userEmail := "test@example.com"
	userID := int64(1)
	verifiedAt := time.Now()

	user := &entity.User{
		ID:              userID,
		Email:           userEmail,
		EmailVerifiedAt: &verifiedAt, // Already verified
	}

	// Mock expectations
	suite.mockUserRepo.On("FindByEmail", userEmail).Return(user, nil)

	// Execute
	err := suite.emailUsecase.ResendVerificationEmail(userEmail)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "email is already verified")
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *EmailUsecaseTestSuite) TestSendPasswordResetEmail_Success() {
	userEmail := "test@example.com"
	userID := int64(1)

	user := &entity.User{
		ID:    userID,
		Email: userEmail,
	}

	// Mock expectations
	suite.mockUserRepo.On("FindByEmail", userEmail).Return(user, nil)
	suite.mockEmailRepo.On("SetPasswordResetToken", userID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)
	suite.mockEmailService.On("SendPasswordResetEmail", userEmail, mock.AnythingOfType("string")).Return(nil)

	// Execute
	err := suite.emailUsecase.SendPasswordResetEmail(userEmail)

	// Assert
	assert.NoError(suite.T(), err)
	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockEmailRepo.AssertExpectations(suite.T())
	suite.mockEmailService.AssertExpectations(suite.T())
}

func (suite *EmailUsecaseTestSuite) TestSendPasswordResetEmail_UserNotFound() {
	userEmail := "nonexistent@example.com"

	// Mock expectations
	suite.mockUserRepo.On("FindByEmail", userEmail).Return(nil, nil)

	// Execute
	err := suite.emailUsecase.SendPasswordResetEmail(userEmail)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "user not found")
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *EmailUsecaseTestSuite) TestResetPassword_Success() {
	token := "valid-token"
	newPassword := "New$ecureP@ssw0rd!" // Valid password that meets complexity requirements and doesn't contain forbidden words or sequential characters
	userID := int64(1)
	now := time.Now()

	user := &entity.User{
		ID:                         userID,
		Email:                      "test@example.com",
		PasswordResetToken:         &token,
		PasswordResetExpiresAt:     &now,
		Password:                   "oldHashedPassword",
		EmailVerificationExpiresAt: &now,
	}

	// Mock expectations
	suite.mockEmailRepo.On("GetUserByPasswordResetToken", token).Return(user, nil)
	suite.mockUserRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(nil)
	suite.mockEmailRepo.On("ClearPasswordResetToken", userID).Return(nil)

	// Execute
	err := suite.emailUsecase.ResetPassword(token, newPassword)

	// Assert
	assert.NoError(suite.T(), err)
	suite.mockEmailRepo.AssertExpectations(suite.T())
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *EmailUsecaseTestSuite) TestResetPassword_InvalidToken() {
	token := "invalid-token"

	// Mock expectations
	suite.mockEmailRepo.On("GetUserByPasswordResetToken", token).Return(nil, nil)

	// Execute
	err := suite.emailUsecase.ResetPassword(token, "New$ecureP@ssw0rd!") // Valid password that meets complexity requirements

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid or expired password reset token")
	suite.mockEmailRepo.AssertExpectations(suite.T())
}

func (suite *EmailUsecaseTestSuite) TestResetPassword_UpdateError() {
	token := "valid-token"
	newPassword := "New$ecureP@ssw0rd!" // Valid password that meets complexity requirements
	userID := int64(1)
	now := time.Now()

	user := &entity.User{
		ID:                         userID,
		Email:                      "test@example.com",
		PasswordResetToken:         &token,
		PasswordResetExpiresAt:     &now,
		Password:                   "oldHashedPassword",
		EmailVerificationExpiresAt: &now,
	}

	// Mock expectations
	suite.mockEmailRepo.On("GetUserByPasswordResetToken", token).Return(user, nil)
	suite.mockUserRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(errors.New("database error"))

	// Execute
	err := suite.emailUsecase.ResetPassword(token, newPassword)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to update password")
	suite.mockEmailRepo.AssertExpectations(suite.T())
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *EmailUsecaseTestSuite) TestCleanupExpiredTokens_Success() {
	// Mock expectations
	suite.mockEmailRepo.On("DeleteExpiredTokens").Return(nil)

	// Execute
	err := suite.emailUsecase.CleanupExpiredTokens()

	// Assert
	assert.NoError(suite.T(), err)
	suite.mockEmailRepo.AssertExpectations(suite.T())
}

func (suite *EmailUsecaseTestSuite) TestCleanupExpiredTokens_Error() {
	// Mock expectations
	suite.mockEmailRepo.On("DeleteExpiredTokens").Return(errors.New("database error"))

	// Execute
	err := suite.emailUsecase.CleanupExpiredTokens()

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to delete expired tokens")
	suite.mockEmailRepo.AssertExpectations(suite.T())
}

func TestEmailUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(EmailUsecaseTestSuite))
}

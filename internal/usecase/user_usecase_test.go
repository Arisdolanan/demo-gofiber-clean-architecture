package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/infrastructure/cache"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/mocks"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// createTestUserForUserTests creates a test user for user usecase testing purposes
func createTestUserForUserTests() *entity.User {
	return &entity.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // hashed "password"
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}
}

type UserUsecaseTestSuite struct {
	suite.Suite
	userUsecase  UserUseCase
	mockUserRepo *mocks.MockUserRepository
	mockCache    *cache.RedisCache // Using concrete type since no interface exists
	validator    *validator.Validate
	logger       *logrus.Logger
}

func (suite *UserUsecaseTestSuite) SetupTest() {
	suite.mockUserRepo = new(mocks.MockUserRepository)
	suite.mockCache = cache.NewRedisCache() // Use real cache for now
	suite.validator = validator.New()
	suite.logger = logrus.New()
	suite.logger.SetLevel(logrus.ErrorLevel) // Suppress logs in tests

	suite.userUsecase = NewUserUseCase(
		suite.mockUserRepo,
		suite.mockCache,
		suite.logger,
		suite.validator,
	)
}

func (suite *UserUsecaseTestSuite) TearDownTest() {
	suite.mockUserRepo.AssertExpectations(suite.T())
	// Skip cache expectations since it's a real cache
}

// Test CreateUser - Success scenario
func (suite *UserUsecaseTestSuite) TestCreateUser_Success() {
	user := createTestUserForUserTests()

	// Mock expectations
	suite.mockUserRepo.On("Create", mock.AnythingOfType("*entity.User")).Return(nil)

	// Execute
	err := suite.userUsecase.CreateUser(user)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test CreateUser - Database error
func (suite *UserUsecaseTestSuite) TestCreateUser_DatabaseError() {
	user := createTestUserForUserTests()
	dbError := errors.New("database connection error")

	// Mock expectations
	suite.mockUserRepo.On("Create", mock.AnythingOfType("*entity.User")).Return(dbError)

	// Execute
	err := suite.userUsecase.CreateUser(user)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), dbError, err)
}

// Test GetUserByID - Success scenario
func (suite *UserUsecaseTestSuite) TestGetUserByID_Success() {
	userID := int64(1)
	expectedUser := createTestUserForUserTests()

	// Mock expectations
	suite.mockUserRepo.On("FindByID", userID).Return(expectedUser, nil)

	// Execute
	user, err := suite.userUsecase.GetUserByID(userID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), expectedUser.ID, user.ID)
	assert.Equal(suite.T(), expectedUser.Email, user.Email)
}

// Test GetUserByID - User not found
func (suite *UserUsecaseTestSuite) TestGetUserByID_UserNotFound() {
	userID := int64(1)

	// Mock expectations
	suite.mockUserRepo.On("FindByID", userID).Return(nil, nil)

	// Execute
	user, err := suite.userUsecase.GetUserByID(userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Equal(suite.T(), "user not found", err.Error())
}

// Test GetUserByID - Database error
func (suite *UserUsecaseTestSuite) TestGetUserByID_DatabaseError() {
	userID := int64(1)
	dbError := errors.New("database connection error")

	// Mock expectations
	suite.mockUserRepo.On("FindByID", userID).Return(nil, dbError)

	// Execute
	user, err := suite.userUsecase.GetUserByID(userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Equal(suite.T(), dbError, err)
}

// Test GetUserByEmail - Success scenario
func (suite *UserUsecaseTestSuite) TestGetUserByEmail_Success() {
	email := "test@example.com"
	expectedUser := createTestUserForUserTests()

	// Mock expectations
	suite.mockUserRepo.On("FindByEmail", email).Return(expectedUser, nil)

	// Execute
	user, err := suite.userUsecase.GetUserByEmail(email)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), expectedUser.Email, user.Email)
}

// Test GetUserByEmail - Database error
func (suite *UserUsecaseTestSuite) TestGetUserByEmail_DatabaseError() {
	email := "test@example.com"
	dbError := errors.New("database connection error")

	// Mock expectations
	suite.mockUserRepo.On("FindByEmail", email).Return(nil, dbError)

	// Execute
	user, err := suite.userUsecase.GetUserByEmail(email)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Equal(suite.T(), dbError, err)
}

// Test GetUserByUsername - Success scenario
func (suite *UserUsecaseTestSuite) TestGetUserByUsername_Success() {
	username := "testuser"
	expectedUser := createTestUserForUserTests()

	// Mock expectations
	suite.mockUserRepo.On("FindByUsername", username).Return(expectedUser, nil)

	// Execute
	user, err := suite.userUsecase.GetUserByUsername(username)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), expectedUser.Username, user.Username)
}

// Test GetUserByUsername - Database error
func (suite *UserUsecaseTestSuite) TestGetUserByUsername_DatabaseError() {
	username := "testuser"
	dbError := errors.New("database connection error")

	// Mock expectations
	suite.mockUserRepo.On("FindByUsername", username).Return(nil, dbError)

	// Execute
	user, err := suite.userUsecase.GetUserByUsername(username)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Equal(suite.T(), dbError, err)
}

// Test GetAllUsers - Success scenario
func (suite *UserUsecaseTestSuite) TestGetAllUsers_Success() {
	page := 1
	pageSize := 10
	expectedUsers := []*entity.User{createTestUserForUserTests()}

	// Mock expectations - the usecase internally converts page/pageSize to limit/offset
	suite.mockUserRepo.On("FindAll", 10, 0).Return(expectedUsers, nil)

	// Execute
	users, err := suite.userUsecase.GetAllUsers(page, pageSize)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), users)
	assert.IsType(suite.T(), &entity.UserListResponse{}, users)
	assert.Len(suite.T(), users.Users, 1)
	assert.Equal(suite.T(), expectedUsers[0].ID, users.Users[0].ID)
	assert.Equal(suite.T(), page, users.Page)
	assert.Equal(suite.T(), pageSize, users.PageSize)
}

// Test GetAllUsers - Database error
func (suite *UserUsecaseTestSuite) TestGetAllUsers_DatabaseError() {
	page := 1
	pageSize := 10
	dbError := errors.New("database connection error")

	// Mock expectations - the usecase internally converts page/pageSize to limit/offset
	suite.mockUserRepo.On("FindAll", 10, 0).Return(nil, dbError)

	// Execute
	users, err := suite.userUsecase.GetAllUsers(page, pageSize)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), users)
	assert.Equal(suite.T(), dbError, err)
}

// Test UpdateUser - Success scenario
func (suite *UserUsecaseTestSuite) TestUpdateUser_Success() {
	user := createTestUserForUserTests()

	// Mock expectations
	suite.mockUserRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(nil)

	// Execute
	err := suite.userUsecase.UpdateUser(user)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test UpdateUser - Database error
func (suite *UserUsecaseTestSuite) TestUpdateUser_DatabaseError() {
	user := createTestUserForUserTests()
	dbError := errors.New("database connection error")

	// Mock expectations
	suite.mockUserRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(dbError)

	// Execute
	err := suite.userUsecase.UpdateUser(user)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), dbError, err)
}

// Test DeleteUser - Success scenario
func (suite *UserUsecaseTestSuite) TestDeleteUser_Success() {
	userID := int64(1)

	// Mock expectations
	suite.mockUserRepo.On("Delete", userID).Return(nil)

	// Execute
	err := suite.userUsecase.DeleteUser(userID)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test DeleteUser - Database error
func (suite *UserUsecaseTestSuite) TestDeleteUser_DatabaseError() {
	userID := int64(1)
	dbError := errors.New("database connection error")

	// Mock expectations
	suite.mockUserRepo.On("Delete", userID).Return(dbError)

	// Execute
	err := suite.userUsecase.DeleteUser(userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), dbError, err)
}

// Test SoftDeleteUser - Success scenario
func (suite *UserUsecaseTestSuite) TestSoftDeleteUser_Success() {
	userID := int64(1)

	// Mock expectations
	suite.mockUserRepo.On("SoftDelete", userID).Return(nil)

	// Execute
	err := suite.userUsecase.SoftDeleteUser(userID)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test SoftDeleteUser - Database error
func (suite *UserUsecaseTestSuite) TestSoftDeleteUser_DatabaseError() {
	userID := int64(1)
	dbError := errors.New("database connection error")

	// Mock expectations
	suite.mockUserRepo.On("SoftDelete", userID).Return(dbError)

	// Execute
	err := suite.userUsecase.SoftDeleteUser(userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), dbError, err)
}

func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}

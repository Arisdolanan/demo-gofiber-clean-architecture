package mocks

import (
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/stretchr/testify/mock"
)

// MockAuthUsecase is a mock implementation of AuthUsecase
type MockAuthUsecase struct {
	mock.Mock
}

func (m *MockAuthUsecase) Register(email, password string) error {
	args := m.Called(email, password)
	return args.Error(0)
}

func (m *MockAuthUsecase) Login(email, password string) (*entity.AuthToken, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthToken), args.Error(1)
}

func (m *MockAuthUsecase) VerifyToken(token string) (*entity.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockAuthUsecase) RefreshToken(refreshToken string) (*entity.AuthToken, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthToken), args.Error(1)
}

func (m *MockAuthUsecase) Logout(userID int64, accessToken string) error {
	args := m.Called(userID, accessToken)
	return args.Error(0)
}

func (m *MockAuthUsecase) ValidatePasswordComplexity(password string) error {
	args := m.Called(password)
	return args.Error(0)
}

func (m *MockAuthUsecase) BlacklistToken(token string, expiration time.Duration) error {
	args := m.Called(token, expiration)
	return args.Error(0)
}

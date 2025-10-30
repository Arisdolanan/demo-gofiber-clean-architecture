package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// MockAuthRedisRepository is a mock implementation of AuthRedisRepository
type MockAuthRedisRepository struct {
	mock.Mock
}

func (m *MockAuthRedisRepository) StoreRefreshToken(userID int64, refreshToken string, expiration time.Duration) error {
	args := m.Called(userID, refreshToken, expiration)
	return args.Error(0)
}

func (m *MockAuthRedisRepository) GetRefreshToken(userID int64) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockAuthRedisRepository) DeleteRefreshToken(userID int64) error {
	args := m.Called(userID)
	return args.Error(0)
}

// Token blacklisting methods
func (m *MockAuthRedisRepository) BlacklistToken(tokenString string, expiration time.Duration) error {
	args := m.Called(tokenString, expiration)
	return args.Error(0)
}

func (m *MockAuthRedisRepository) IsTokenBlacklisted(tokenString string) (bool, error) {
	args := m.Called(tokenString)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRedisRepository) BlacklistAllUserTokens(userID int64, expiration time.Duration) error {
	args := m.Called(userID, expiration)
	return args.Error(0)
}

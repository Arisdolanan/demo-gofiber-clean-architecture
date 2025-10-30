package mocks

import (
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/stretchr/testify/mock"
)

// MockAuthRepository is a mock implementation of AuthRepository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) Register(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAuthRepository) FindByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}
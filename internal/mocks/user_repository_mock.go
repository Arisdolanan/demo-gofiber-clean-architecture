package mocks

import (
	"context"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) CreateWithContext(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, schoolID int64, id int64) (*entity.User, error) {
	args := m.Called(ctx, schoolID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (*entity.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context, schoolID int64, limit, offset int) ([]*entity.User, error) {
	args := m.Called(ctx, schoolID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindAllBySchool(ctx context.Context, schoolID int64, limit, offset int) ([]*entity.User, error) {
	args := m.Called(ctx, schoolID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByType(ctx context.Context, schoolID int64, userType entity.UserType, limit, offset int) ([]*entity.User, error) {
	args := m.Called(ctx, schoolID, userType, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, schoolID int64, user *entity.User) error {
	args := m.Called(ctx, schoolID, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateWithContext(ctx context.Context, schoolID int64, user *entity.User) error {
	args := m.Called(ctx, schoolID, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, schoolID int64, id int64) error {
	args := m.Called(ctx, schoolID, id)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteWithContext(ctx context.Context, schoolID int64, id int64) error {
	args := m.Called(ctx, schoolID, id)
	return args.Error(0)
}

func (m *MockUserRepository) SoftDelete(ctx context.Context, schoolID int64, id int64) error {
	args := m.Called(ctx, schoolID, id)
	return args.Error(0)
}

func (m *MockUserRepository) SoftDeleteWithContext(ctx context.Context, schoolID int64, id int64) error {
	args := m.Called(ctx, schoolID, id)
	return args.Error(0)
}
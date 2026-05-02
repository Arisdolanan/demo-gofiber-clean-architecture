package mocks

import (
	"context"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/stretchr/testify/mock"
)

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) Create(ctx context.Context, role *entity.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) Update(ctx context.Context, role *entity.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) FindByID(ctx context.Context, id int64) (*entity.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) FindByCode(ctx context.Context, schoolID *int64, code string) (*entity.Role, error) {
	args := m.Called(ctx, schoolID, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) FindAllBySchool(ctx context.Context, schoolID *int64) ([]*entity.Role, error) {
	args := m.Called(ctx, schoolID)
	return args.Get(0).([]*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleRepository) AssignPermission(ctx context.Context, roleID, permissionID int64) error {
	args := m.Called(ctx, roleID, permissionID)
	return args.Error(0)
}

func (m *MockRoleRepository) RemovePermission(ctx context.Context, roleID, permissionID int64) error {
	args := m.Called(ctx, roleID, permissionID)
	return args.Error(0)
}

func (m *MockRoleRepository) GetRolePermissions(ctx context.Context, roleID int64) ([]*entity.Permission, error) {
	args := m.Called(ctx, roleID)
	return args.Get(0).([]*entity.Permission), args.Error(1)
}

func (m *MockRoleRepository) AssignUserRole(ctx context.Context, schoolID int64, userID, roleID int64) error {
	args := m.Called(ctx, schoolID, userID, roleID)
	return args.Error(0)
}

func (m *MockRoleRepository) RemoveUserRole(ctx context.Context, schoolID int64, userID, roleID int64) error {
	args := m.Called(ctx, schoolID, userID, roleID)
	return args.Error(0)
}

func (m *MockRoleRepository) GetUserRoles(ctx context.Context, schoolID int64, userID int64) ([]*entity.Role, error) {
	args := m.Called(ctx, schoolID, userID)
	return args.Get(0).([]*entity.Role), args.Error(1)
}

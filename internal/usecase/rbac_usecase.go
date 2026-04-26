package usecase

import (
	"context"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/repository/postgresql"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type RBACUsecase interface {
	CreateRole(ctx context.Context, role *entity.Role) error
	GetRolesBySchool(ctx context.Context, schoolID *int64) ([]*entity.Role, error)
	AssignPermissionToRole(ctx context.Context, roleID, permissionID int64) error

	// User assignments
	AssignRoleToUser(ctx context.Context, userID, roleID int64) error
	GetUserRoles(ctx context.Context, userID int64) ([]*entity.Role, error)
	GetUserPermissions(ctx context.Context, userID int64) ([]string, error)
	CheckPermission(ctx context.Context, userID int64, permissionCode string) (bool, error)

	// Permission management
	CreatePermission(ctx context.Context, perm *entity.Permission) error
	GetAllPermissions(ctx context.Context) ([]*entity.Permission, error)
}

type rbacUsecase struct {
	roleRepo postgresql.RoleRepository
	permRepo postgresql.PermissionRepository
	validate *validator.Validate
	log      *logrus.Logger
}

func NewRBACUsecase(roleRepo postgresql.RoleRepository, permRepo postgresql.PermissionRepository, validate *validator.Validate, log *logrus.Logger) RBACUsecase {
	return &rbacUsecase{
		roleRepo: roleRepo,
		permRepo: permRepo,
		validate: validate,
		log:      log,
	}
}

func (uc *rbacUsecase) CreateRole(ctx context.Context, role *entity.Role) error {
	if err := uc.validate.Struct(role); err != nil {
		return err
	}
	return uc.roleRepo.Create(ctx, role)
}

func (uc *rbacUsecase) GetRolesBySchool(ctx context.Context, schoolID *int64) ([]*entity.Role, error) {
	return uc.roleRepo.FindAllBySchool(ctx, schoolID)
}

func (uc *rbacUsecase) AssignPermissionToRole(ctx context.Context, roleID, permissionID int64) error {
	return uc.roleRepo.AssignPermission(ctx, roleID, permissionID)
}

func (uc *rbacUsecase) AssignRoleToUser(ctx context.Context, userID, roleID int64) error {
	return uc.roleRepo.AssignUserRole(ctx, userID, roleID)
}

func (uc *rbacUsecase) GetUserRoles(ctx context.Context, userID int64) ([]*entity.Role, error) {
	return uc.roleRepo.GetUserRoles(ctx, userID)
}

func (uc *rbacUsecase) GetUserPermissions(ctx context.Context, userID int64) ([]string, error) {
	roles, err := uc.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	permMap := make(map[string]bool)
	for _, role := range roles {
		perms, err := uc.roleRepo.GetRolePermissions(ctx, role.ID)
		if err != nil {
			uc.log.Errorf("Error getting permissions for role %d: %v", role.ID, err)
			continue
		}
		for _, p := range perms {
			permMap[p.PermissionCode] = true
		}
	}

	uniquePerms := make([]string, 0, len(permMap))
	for p := range permMap {
		uniquePerms = append(uniquePerms, p)
	}

	return uniquePerms, nil
}

func (uc *rbacUsecase) CheckPermission(ctx context.Context, userID int64, permissionCode string) (bool, error) {
	perms, err := uc.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, p := range perms {
		if p == permissionCode {
			return true, nil
		}
	}

	return false, nil
}

func (uc *rbacUsecase) CreatePermission(ctx context.Context, perm *entity.Permission) error {
	if err := uc.validate.Struct(perm); err != nil {
		return err
	}
	return uc.permRepo.Create(ctx, perm)
}

func (uc *rbacUsecase) GetAllPermissions(ctx context.Context) ([]*entity.Permission, error) {
	return uc.permRepo.FindAll(ctx)
}

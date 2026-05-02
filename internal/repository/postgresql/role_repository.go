package postgresql

import (
	"context"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/jmoiron/sqlx"
)

type RoleRepository interface {
	Create(ctx context.Context, role *entity.Role) error
	Update(ctx context.Context, role *entity.Role) error
	FindByID(ctx context.Context, id int64) (*entity.Role, error)
	FindByCode(ctx context.Context, schoolID *int64, code string) (*entity.Role, error)
	FindAllBySchool(ctx context.Context, schoolID *int64) ([]*entity.Role, error)
	Delete(ctx context.Context, id int64) error

	// Role-Permission assignments
	AssignPermission(ctx context.Context, roleID, permissionID int64) error
	RemovePermission(ctx context.Context, roleID, permissionID int64) error
	GetRolePermissions(ctx context.Context, roleID int64) ([]*entity.Permission, error)

	// User-Role assignments
	AssignUserRole(ctx context.Context, schoolID int64, userID, roleID int64) error
	RemoveUserRole(ctx context.Context, schoolID int64, userID, roleID int64) error
	GetUserRoles(ctx context.Context, schoolID int64, userID int64) ([]*entity.Role, error)
}

type roleRepository struct {
	roleRepo           *BaseRepository[entity.Role]
	rolePermissionRepo *BaseRepository[entity.RolePermission]
	userRoleRepo       *BaseRepository[entity.UserRole]
	db                 *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) RoleRepository {
	return &roleRepository{
		roleRepo:           NewBaseRepository[entity.Role](db, "roles"),
		rolePermissionRepo: NewBaseRepository[entity.RolePermission](db, "role_permissions"),
		userRoleRepo:       NewBaseRepository[entity.UserRole](db, "user_roles"),
		db:                 db,
	}
}

func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
	return r.roleRepo.Create(ctx, role)
}

func (r *roleRepository) Update(ctx context.Context, role *entity.Role) error {
	return r.roleRepo.Update(ctx, role, "id = $1 AND deleted_at IS NULL", role.ID)
}

func (r *roleRepository) FindByID(ctx context.Context, id int64) (*entity.Role, error) {
	return r.roleRepo.FindByID(ctx, id)
}

func (r *roleRepository) FindByCode(ctx context.Context, schoolID *int64, code string) (*entity.Role, error) {
	if schoolID == nil {
		return r.roleRepo.FindOne(ctx, "school_id IS NULL AND code = $1 AND deleted_at IS NULL", code)
	}
	return r.roleRepo.FindOne(ctx, "school_id = $1 AND code = $2 AND deleted_at IS NULL", *schoolID, code)
}

func (r *roleRepository) FindAllBySchool(ctx context.Context, schoolID *int64) ([]*entity.Role, error) {
	if schoolID == nil {
		return r.roleRepo.FindAll(ctx, "school_id IS NULL AND deleted_at IS NULL")
	}
	return r.roleRepo.FindAll(ctx, "school_id = $1 AND deleted_at IS NULL", *schoolID)
}

func (r *roleRepository) Delete(ctx context.Context, id int64) error {
	return r.roleRepo.SoftDelete(ctx, "id = $1 AND deleted_at IS NULL", id)
}

// Role-Permission
func (r *roleRepository) AssignPermission(ctx context.Context, roleID, permissionID int64) error {
	rp := &entity.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	// Use manual exec to support ON CONFLICT if needed, or just Create
	return r.rolePermissionRepo.Create(ctx, rp)
}

func (r *roleRepository) RemovePermission(ctx context.Context, roleID, permissionID int64) error {
	return r.rolePermissionRepo.Delete(ctx, "role_id = $1 AND permission_id = $2", roleID, permissionID)
}

func (r *roleRepository) GetRolePermissions(ctx context.Context, roleID int64) ([]*entity.Permission, error) {
	var permissions []*entity.Permission
	query := `
		SELECT p.* FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1 AND p.deleted_at IS NULL
	`
	err := r.db.SelectContext(ctx, &permissions, query, roleID)
	return permissions, err
}

// User-Role
func (r *roleRepository) AssignUserRole(ctx context.Context, schoolID int64, userID, roleID int64) error {
	ur := &entity.UserRole{
		SchoolID: schoolID,
		UserID:   userID,
		RoleID:   roleID,
	}
	// Context should contain user_id if called via middleware, which BaseRepository will use for assigned_by
	return r.userRoleRepo.Create(ctx, ur)
}

func (r *roleRepository) RemoveUserRole(ctx context.Context, schoolID int64, userID, roleID int64) error {
	return r.userRoleRepo.Delete(ctx, "school_id = $1 AND user_id = $2 AND role_id = $3", schoolID, userID, roleID)
}

func (r *roleRepository) GetUserRoles(ctx context.Context, schoolID int64, userID int64) ([]*entity.Role, error) {
	var roles []*entity.Role
	query := `
		SELECT r.* FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.school_id = $1 AND ur.user_id = $2 AND r.deleted_at IS NULL
	`
	err := r.db.SelectContext(ctx, &roles, query, schoolID, userID)
	return roles, err
}

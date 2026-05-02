package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	Register(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByID(id int64) (*entity.User, error)
}

type authRepository struct {
	*BaseRepository[entity.User]
}

func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepository{
		BaseRepository: NewBaseRepository[entity.User](db, "users"),
	}
}

func (r *authRepository) Register(user *entity.User) error {
	ctx := context.Background()
	return r.Create(ctx, user)
}

func (r *authRepository) FindByID(id int64) (*entity.User, error) {
	ctx := context.Background()
	return r.loadUserRelations(ctx, "u.id = $1", id)
}

func (r *authRepository) FindByEmail(email string) (*entity.User, error) {
	ctx := context.Background()
	return r.loadUserRelations(ctx, "u.email = $1", email)
}

func (r *authRepository) loadUserRelations(ctx context.Context, where string, arg interface{}) (*entity.User, error) {
	// Use a temporary struct to handle NULL values from LEFT JOIN safely (for Role)
	type userRow struct {
		entity.User
		// Role fields (Nullable)
		RoleID           sql.NullInt64  `db:"role.id"`
		RoleCode         sql.NullString `db:"role.code"`
		RoleName         sql.NullString `db:"role.name"`
		RoleIsSystemRole sql.NullBool   `db:"role.is_system_role"`
	}

	var row userRow
	query := fmt.Sprintf(`
		SELECT 
			u.*,
			r.id AS "role.id", r.code AS "role.code", r.name AS "role.name", r.is_system_role AS "role.is_system_role"
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE %s AND u.deleted_at IS NULL
		LIMIT 1
	`, where)

	err := r.GetContext(ctx, &row, query, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	user := &row.User

	// 1. Map Role if it exists
	if row.RoleID.Valid {
		user.Role = &entity.Role{
			ID:           row.RoleID.Int64,
			Code:         row.RoleCode.String,
			Name:         row.RoleName.String,
			IsSystemRole: row.RoleIsSystemRole.Bool,
		}

		// Fetch Permissions for the user's role
		permQuery := `
			SELECT p.* 
			FROM permissions p
			JOIN role_permissions rp ON p.id = rp.permission_id
			WHERE rp.role_id = $1
		`
		var permissions []entity.Permission
		err = r.SelectContext(ctx, &permissions, permQuery, user.Role.ID)
		if err == nil {
			user.Permissions = permissions
		}
	}

	// 2. Fetch Schools based on school_id array
	if len(user.SchoolID) > 0 {
		schoolQuery := `SELECT * FROM schools WHERE id = ANY($1) AND deleted_at IS NULL ORDER BY id ASC`
		var schools []entity.School
		err = r.SelectContext(ctx, &schools, schoolQuery, user.SchoolID)
		if err == nil {
			user.AccessibleSchools = schools
			// Set the primary school (first one in the list) to the School field for compatibility
			if len(schools) > 0 {
				user.School = &schools[0]
			}
		}
	} else if user.UserType == entity.UserSuperAdmin {
		// Fallback for Super Admin if no specific schools assigned
		schoolQuery := `SELECT * FROM schools WHERE deleted_at IS NULL ORDER BY id ASC`
		var allSchools []entity.School
		err = r.SelectContext(ctx, &allSchools, schoolQuery)
		if err == nil {
			user.AccessibleSchools = allSchools
		}
	}

	return user, nil
}

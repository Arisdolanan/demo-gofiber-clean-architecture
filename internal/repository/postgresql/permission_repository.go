package postgresql

import (
	"context"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/jmoiron/sqlx"
)

type PermissionRepository interface {
	Create(ctx context.Context, perm *entity.Permission) error
	FindAll(ctx context.Context) ([]*entity.Permission, error)
	FindByCode(ctx context.Context, code string) (*entity.Permission, error)
	FindByModule(ctx context.Context, module string) ([]*entity.Permission, error)
}

type permissionRepository struct {
	*BaseRepository[entity.Permission]
}

func NewPermissionRepository(db *sqlx.DB) PermissionRepository {
	return &permissionRepository{
		BaseRepository: NewBaseRepository[entity.Permission](db, "permissions"),
	}
}

func (r *permissionRepository) Create(ctx context.Context, perm *entity.Permission) error {
	return r.BaseRepository.Create(ctx, perm)
}

func (r *permissionRepository) FindAll(ctx context.Context) ([]*entity.Permission, error) {
	return r.BaseRepository.FindAll(ctx, "")
}

func (r *permissionRepository) FindByCode(ctx context.Context, code string) (*entity.Permission, error) {
	return r.BaseRepository.FindOne(ctx, "permission_code = $1", code)
}

func (r *permissionRepository) FindByModule(ctx context.Context, module string) ([]*entity.Permission, error) {
	return r.BaseRepository.FindAll(ctx, "module_name = $1", module)
}

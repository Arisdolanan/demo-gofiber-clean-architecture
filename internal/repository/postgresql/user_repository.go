package postgresql

import (
	"context"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(user *entity.User) error
	CreateWithContext(ctx context.Context, user *entity.User) error
	FindByID(id int64) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByUsername(username string) (*entity.User, error)
	FindAll(limit, offset int) ([]*entity.User, error)
	Update(user *entity.User) error
	UpdateWithContext(ctx context.Context, user *entity.User) error
	Delete(id int64) error
	DeleteWithContext(ctx context.Context, id int64) error
	SoftDelete(id int64) error
	SoftDeleteWithContext(ctx context.Context, id int64) error
}

type userRepository struct {
	*BaseRepository[entity.User]
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository[entity.User](db, "users"),
	}
}

func (r *userRepository) Create(user *entity.User) error {
	ctx := context.Background()
	return r.BaseRepository.Create(ctx, user)
}

func (r *userRepository) CreateWithContext(ctx context.Context, user *entity.User) error {
	return r.BaseRepository.Create(ctx, user)
}

func (r *userRepository) FindByID(id int64) (*entity.User, error) {
	ctx := context.Background()
	return r.BaseRepository.FindByID(ctx, id)
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	ctx := context.Background()
	return r.BaseRepository.FindOne(ctx, "email = $1 AND deleted_at IS NULL", email)
}

func (r *userRepository) FindByUsername(username string) (*entity.User, error) {
	ctx := context.Background()
	return r.BaseRepository.FindOne(ctx, "username = $1 AND deleted_at IS NULL", username)
}

func (r *userRepository) FindAll(limit, offset int) ([]*entity.User, error) {
	ctx := context.Background()
	return r.BaseRepository.FindAllWithPagination(ctx, limit, offset, "deleted_at IS NULL")
}

func (r *userRepository) Update(user *entity.User) error {
	ctx := context.Background()
	return r.BaseRepository.Update(ctx, user, "id = $1 AND deleted_at IS NULL", user.ID)
}

func (r *userRepository) UpdateWithContext(ctx context.Context, user *entity.User) error {
	return r.BaseRepository.Update(ctx, user, "id = $1 AND deleted_at IS NULL", user.ID)
}

func (r *userRepository) Delete(id int64) error {
	ctx := context.Background()
	return r.BaseRepository.Delete(ctx, "id = $1", id)
}

func (r *userRepository) DeleteWithContext(ctx context.Context, id int64) error {
	return r.BaseRepository.Delete(ctx, "id = $1", id)
}

func (r *userRepository) SoftDelete(id int64) error {
	ctx := context.Background()
	return r.BaseRepository.SoftDelete(ctx, "id = $1 AND deleted_at IS NULL", id)
}

func (r *userRepository) SoftDeleteWithContext(ctx context.Context, id int64) error {
	return r.BaseRepository.SoftDelete(ctx, "id = $1 AND deleted_at IS NULL", id)
}

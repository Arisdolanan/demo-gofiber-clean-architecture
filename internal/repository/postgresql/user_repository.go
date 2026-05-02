package postgresql

import (
	"context"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(user *entity.User) error
	CreateWithContext(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, schoolID int64, id int64) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByUsername(username string) (*entity.User, error)
	FindAll(ctx context.Context, schoolID int64, limit, offset int) ([]*entity.User, error)
	FindAllBySchool(ctx context.Context, schoolID int64, limit, offset int) ([]*entity.User, error)
	FindByType(ctx context.Context, schoolID int64, userType entity.UserType, limit, offset int) ([]*entity.User, error)
	Update(ctx context.Context, schoolID int64, user *entity.User) error
	UpdateWithContext(ctx context.Context, schoolID int64, user *entity.User) error
	Delete(ctx context.Context, schoolID int64, id int64) error
	DeleteWithContext(ctx context.Context, schoolID int64, id int64) error
	SoftDelete(ctx context.Context, schoolID int64, id int64) error
	SoftDeleteWithContext(ctx context.Context, schoolID int64, id int64) error
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

func (r *userRepository) FindByID(ctx context.Context, schoolID int64, id int64) (*entity.User, error) {
	return r.BaseRepository.FindOne(ctx, "id = $1 AND $2 = ANY(school_id) AND deleted_at IS NULL", id, schoolID)
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	ctx := context.Background()
	return r.BaseRepository.FindOne(ctx, "email = $1 AND deleted_at IS NULL", email)
}

func (r *userRepository) FindByUsername(username string) (*entity.User, error) {
	ctx := context.Background()
	return r.BaseRepository.FindOne(ctx, "username = $1 AND deleted_at IS NULL", username)
}

func (r *userRepository) FindAll(ctx context.Context, schoolID int64, limit, offset int) ([]*entity.User, error) {
	return r.BaseRepository.FindAllWithPagination(ctx, limit, offset, "$1 = ANY(school_id) AND deleted_at IS NULL", schoolID)
}

func (r *userRepository) FindAllBySchool(ctx context.Context, schoolID int64, limit, offset int) ([]*entity.User, error) {
	return r.BaseRepository.FindAllWithPagination(ctx, limit, offset, "$1 = ANY(school_id) AND deleted_at IS NULL", schoolID)
}

func (r *userRepository) FindByType(ctx context.Context, schoolID int64, userType entity.UserType, limit, offset int) ([]*entity.User, error) {
	return r.BaseRepository.FindAllWithPagination(ctx, limit, offset, "user_type = $1 AND $2 = ANY(school_id) AND deleted_at IS NULL", userType, schoolID)
}

func (r *userRepository) Update(ctx context.Context, schoolID int64, user *entity.User) error {
	return r.BaseRepository.Update(ctx, user, "id = $1 AND $2 = ANY(school_id) AND deleted_at IS NULL", user.ID, schoolID)
}

func (r *userRepository) UpdateWithContext(ctx context.Context, schoolID int64, user *entity.User) error {
	return r.BaseRepository.Update(ctx, user, "id = $1 AND $2 = ANY(school_id) AND deleted_at IS NULL", user.ID, schoolID)
}

func (r *userRepository) Delete(ctx context.Context, schoolID int64, id int64) error {
	return r.BaseRepository.Delete(ctx, "id = $1 AND $2 = ANY(school_id)", id, schoolID)
}

func (r *userRepository) DeleteWithContext(ctx context.Context, schoolID int64, id int64) error {
	return r.BaseRepository.Delete(ctx, "id = $1 AND $2 = ANY(school_id)", id, schoolID)
}

func (r *userRepository) SoftDelete(ctx context.Context, schoolID int64, id int64) error {
	return r.BaseRepository.SoftDelete(ctx, "id = $1 AND $2 = ANY(school_id) AND deleted_at IS NULL", id, schoolID)
}

func (r *userRepository) SoftDeleteWithContext(ctx context.Context, schoolID int64, id int64) error {
	return r.BaseRepository.SoftDelete(ctx, "id = $1 AND $2 = ANY(school_id) AND deleted_at IS NULL", id, schoolID)
}

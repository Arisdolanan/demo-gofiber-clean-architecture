package postgresql

import (
	"context"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	Register(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
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

func (r *authRepository) FindByEmail(email string) (*entity.User, error) {
	ctx := context.Background()
	return r.FindOne(ctx, "email = $1 AND deleted_at IS NULL", email)
}

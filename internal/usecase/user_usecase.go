package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/infrastructure/cache"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/repository/postgresql"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByID(id int64) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
	GetAllUsers(page, pageSize int) (*entity.UserListResponse, error)
	GetUsersBySchool(schoolID int64, page, pageSize int) (*entity.UserListResponse, error)
	GetUsersByType(userType entity.UserType, page, pageSize int) (*entity.UserListResponse, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id int64) error
	SoftDeleteUser(ctx context.Context, id int64) error
}

type userUseCase struct {
	userRepo   postgresql.UserRepository
	redisCache *cache.RedisCache
	log        *logrus.Logger
	validate   *validator.Validate
}

func NewUserUseCase(
	userRepo postgresql.UserRepository,
	redisCache *cache.RedisCache,
	log *logrus.Logger,
	validate *validator.Validate,
) UserUseCase {
	return &userUseCase{
		userRepo:   userRepo,
		redisCache: redisCache,
		log:        log,
		validate:   validate,
	}
}

// CreateUser creates a new user
func (uc *userUseCase) CreateUser(ctx context.Context, user *entity.User) error {
	// Set timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Create user in database
	if err := uc.userRepo.CreateWithContext(ctx, user); err != nil {
		uc.log.Errorf("Error creating user: %v", err)
		return err
	}

	uc.log.Infof("User created successfully: %s", user.Email)
	return nil
}

// GetUserByID retrieves a user by ID
func (uc *userUseCase) GetUserByID(id int64) (*entity.User, error) {
	user, err := uc.userRepo.FindByID(id)
	if err != nil {
		uc.log.Errorf("Error getting user by ID %d: %v", id, err)
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (uc *userUseCase) GetUserByEmail(email string) (*entity.User, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		uc.log.Errorf("Error getting user by email %s: %v", email, err)
		return nil, err
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (uc *userUseCase) GetUserByUsername(username string) (*entity.User, error) {
	user, err := uc.userRepo.FindByUsername(username)
	if err != nil {
		uc.log.Errorf("Error getting user by username %s: %v", username, err)
		return nil, err
	}

	return user, nil
}

// GetAllUsers retrieves all users with pagination
func (uc *userUseCase) GetAllUsers(page, pageSize int) (*entity.UserListResponse, error) {
	// Normalize pagination parameters
	pagination := utils.CalculatePagination(page, pageSize, 0)

	users, err := uc.userRepo.FindAll(pagination.PageSize, pagination.Offset)
	if err != nil {
		uc.log.Errorf("Error getting all users: %v", err)
		return nil, err
	}

	// Get total count for proper pagination calculation
	// Note: You may need to add a CountAll method to your repository
	// For now, we'll use the length of returned users as an approximation
	totalCount := int64(len(users))

	// Recalculate pagination with actual total count
	pagination = utils.CalculatePagination(page, pageSize, totalCount)

	// Convert []*entity.User to []entity.User
	userList := make([]entity.User, len(users))
	for i, user := range users {
		if user != nil {
			userList[i] = *user
		}
	}

	return &entity.UserListResponse{
		Users:      userList,
		TotalCount: totalCount,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: pagination.TotalPages,
	}, nil
}

// GetUsersBySchool retrieves all users for a specific school with pagination
func (uc *userUseCase) GetUsersBySchool(schoolID int64, page, pageSize int) (*entity.UserListResponse, error) {
	pagination := utils.CalculatePagination(page, pageSize, 0)

	users, err := uc.userRepo.FindAllBySchool(schoolID, pagination.PageSize, pagination.Offset)
	if err != nil {
		uc.log.Errorf("Error getting users for school %d: %v", schoolID, err)
		return nil, err
	}

	totalCount := int64(len(users))
	pagination = utils.CalculatePagination(page, pageSize, totalCount)

	userList := make([]entity.User, len(users))
	for i, user := range users {
		if user != nil {
			userList[i] = *user
		}
	}

	return &entity.UserListResponse{
		Users:      userList,
		TotalCount: totalCount,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: pagination.TotalPages,
	}, nil
}

// GetUsersByType retrieves all users of a specific type with pagination
func (uc *userUseCase) GetUsersByType(userType entity.UserType, page, pageSize int) (*entity.UserListResponse, error) {
	pagination := utils.CalculatePagination(page, pageSize, 0)

	users, err := uc.userRepo.FindByType(userType, pagination.PageSize, pagination.Offset)
	if err != nil {
		uc.log.Errorf("Error getting users by type %s: %v", userType, err)
		return nil, err
	}

	totalCount := int64(len(users))
	pagination = utils.CalculatePagination(page, pageSize, totalCount)

	userList := make([]entity.User, len(users))
	for i, user := range users {
		if user != nil {
			userList[i] = *user
		}
	}

	return &entity.UserListResponse{
		Users:      userList,
		TotalCount: totalCount,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: pagination.TotalPages,
	}, nil
}

// UpdateUser updates an existing user
func (uc *userUseCase) UpdateUser(ctx context.Context, user *entity.User) error {
	// Set updated timestamp
	user.UpdatedAt = time.Now()

	// Update user in database
	if err := uc.userRepo.UpdateWithContext(ctx, user); err != nil {
		uc.log.Errorf("Error updating user %d: %v", user.ID, err)
		return err
	}

	uc.log.Infof("User updated successfully: %d", user.ID)
	return nil
}

// DeleteUser permanently deletes a user
func (uc *userUseCase) DeleteUser(ctx context.Context, id int64) error {
	if err := uc.userRepo.DeleteWithContext(ctx, id); err != nil {
		uc.log.Errorf("Error deleting user %d: %v", id, err)
		return err
	}

	uc.log.Infof("User deleted successfully: %d", id)
	return nil
}

// SoftDeleteUser soft deletes a user
func (uc *userUseCase) SoftDeleteUser(ctx context.Context, id int64) error {
	if err := uc.userRepo.SoftDeleteWithContext(ctx, id); err != nil {
		uc.log.Errorf("Error soft deleting user %d: %v", id, err)
		return err
	}

	uc.log.Infof("User soft deleted successfully: %d", id)
	return nil
}

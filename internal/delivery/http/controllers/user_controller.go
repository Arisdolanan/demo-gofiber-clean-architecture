package controllers

import (
	"context"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/usecase"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log     *logrus.Logger
	UseCase usecase.UserUseCase
}

func NewUserController(useCase usecase.UserUseCase, logger *logrus.Logger) *UserController {
	return &UserController{
		Log:     logger,
		UseCase: useCase,
	}
}

// CreateUser handles user creation
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body entity.User true "User information"
// @Success 201 {object} response.HTTPSuccessResponse{data=entity.User} "User created successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid request body"
// @Failure 500 {object} response.HTTPErrorResponse "Failed to create user"
// @Router /api/v1/users [post]
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	var user entity.User

	if err := ctx.BodyParser(&user); err != nil {
		c.Log.Errorf("Error parsing request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: err.Error()},
			},
		})
	}

	// Create context with user ID from authenticated user (if available)
	reqCtx := context.Background()
	if userID := ctx.Locals("user_id"); userID != nil {
		if userIDVal, ok := userID.(int64); ok {
			reqCtx = context.WithValue(reqCtx, "user_id", userIDVal)
		}
	}

	// Ensure user is created within the context of the current school
	schoolID, err := utils.GetSchoolIDFromToken(ctx)
	if err == nil && schoolID > 0 {
		// Only add schoolID if it's not already in the array
		hasSchool := false
		for _, id := range user.SchoolID {
			if id == schoolID {
				hasSchool = true
				break
			}
		}
		if !hasSchool {
			user.SchoolID = append(user.SchoolID, schoolID)
		}
	}

	if err := c.UseCase.CreateUser(reqCtx, &user); err != nil {
		c.Log.Errorf("Error creating user: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to create user",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusCreated,
		Message: "User created successfully",
		Data:    user,
	})
}

// GetUserByID handles getting user by ID
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.HTTPSuccessResponse{data=entity.User} "User retrieved successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid user ID"
// @Failure 404 {object} response.HTTPErrorResponse "User not found"
// @Router /api/v1/users/{id} [get]
func (c *UserController) GetUserByID(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid user ID",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: "Invalid user ID format"},
			},
		})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)
	user, err := c.UseCase.GetUserByID(ctx.Context(), schoolID, id)
	if err != nil {
		c.Log.Errorf("Error getting user by ID %d: %v", id, err)
		return ctx.Status(fiber.StatusNotFound).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusNotFound,
			Message: "User not found",
			Errors: []response.JSONError{
				{Status: fiber.StatusNotFound, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// GetAllUsers handles getting all users with pagination
// @Summary Get all users
// @Description Get all users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param limit query int false "Limit (backward compatibility)" default(10)
// @Param offset query int false "Offset (backward compatibility)" default(0)
// @Success 200 {object} response.HTTPSuccessResponse{data=entity.UserListResponse} "Users retrieved successfully"
// @Failure 500 {object} response.HTTPErrorResponse "Failed to retrieve users"
// @Router /api/v1/users [get]
func (c *UserController) GetAllUsers(ctx *fiber.Ctx) error {
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)
	params := utils.ParsePaginationQuery(ctx)
	// Use page/pageSize approach for the usecase
	users, err := c.UseCase.GetAllUsers(ctx.Context(), schoolID, params.Page, params.PageSize)
	if err != nil {
		c.Log.Errorf("Error getting all users: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve users",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

// UpdateUser handles user update
// @Summary Update user
// @Description Update an existing user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body entity.User true "User information"
// @Success 200 {object} response.HTTPSuccessResponse{data=entity.User} "User updated successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid request body"
// @Failure 404 {object} response.HTTPErrorResponse "User not found"
// @Failure 500 {object} response.HTTPErrorResponse "Failed to update user"
// @Router /api/v1/users/{id} [put]
func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid user ID",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: "Invalid user ID format"},
			},
		})
	}

	var user entity.User
	user.ID = id

	if err := ctx.BodyParser(&user); err != nil {
		c.Log.Errorf("Error parsing request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: err.Error()},
			},
		})
	}

	// Create context with user ID from authenticated user
	reqCtx := context.Background()
	if userID := ctx.Locals("user_id"); userID != nil {
		if userIDVal, ok := userID.(int64); ok {
			reqCtx = context.WithValue(reqCtx, "user_id", userIDVal)
		}
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.UseCase.UpdateUser(reqCtx, schoolID, &user); err != nil {
		c.Log.Errorf("Error updating user %d: %v", id, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to update user",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "User updated successfully",
		Data:    user,
	})
}

// DeleteUser handles user deletion
// @Summary Delete user
// @Description Permanently delete a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.HTTPSuccessResponse "User deleted successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid user ID"
// @Failure 404 {object} response.HTTPErrorResponse "User not found"
// @Failure 500 {object} response.HTTPErrorResponse "Failed to delete user"
// @Router /api/v1/users/{id} [delete]
func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid user ID",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: "Invalid user ID format"},
			},
		})
	}

	// Create context with user ID from authenticated user
	reqCtx := context.Background()
	if userID := ctx.Locals("user_id"); userID != nil {
		if userIDVal, ok := userID.(int64); ok {
			reqCtx = context.WithValue(reqCtx, "user_id", userIDVal)
		}
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.UseCase.DeleteUser(reqCtx, schoolID, id); err != nil {
		c.Log.Errorf("Error deleting user %d: %v", id, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to delete user",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "User deleted successfully",
	})
}

// SoftDeleteUser handles user soft deletion
// @Summary Soft delete user
// @Description Soft delete a user (mark as deleted without removing from database)
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.HTTPSuccessResponse "User soft deleted successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid user ID"
// @Failure 404 {object} response.HTTPErrorResponse "User not found"
// @Failure 500 {object} response.HTTPErrorResponse "Failed to soft delete user"
// @Router /api/v1/users/{id}/soft-delete [delete]
func (c *UserController) SoftDeleteUser(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid user ID",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: "Invalid user ID format"},
			},
		})
	}

	// Create context with user ID from authenticated user
	reqCtx := context.Background()
	if userID := ctx.Locals("user_id"); userID != nil {
		if userIDVal, ok := userID.(int64); ok {
			reqCtx = context.WithValue(reqCtx, "user_id", userIDVal)
		}
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.UseCase.SoftDeleteUser(reqCtx, schoolID, id); err != nil {
		c.Log.Errorf("Error soft deleting user %d: %v", id, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to soft delete user",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "User soft deleted successfully",
	})
}

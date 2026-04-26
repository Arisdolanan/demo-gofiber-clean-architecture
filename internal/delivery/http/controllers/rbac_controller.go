package controllers

import (
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/usecase"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type RBACController struct {
	usecase usecase.RBACUsecase
	log     *logrus.Logger
}

func NewRBACController(uc usecase.RBACUsecase, log *logrus.Logger) *RBACController {
	return &RBACController{
		usecase: uc,
		log:     log,
	}
}

// CreateRole handles role creation
// @Summary Create a new role
// @Tags rbac
// @Accept json
// @Produce json
// @Param role body entity.Role true "Role details"
// @Success 201 {object} response.HTTPSuccessResponse
// @Router /api/v1/rbac/roles [post]
func (c *RBACController) CreateRole(ctx *fiber.Ctx) error {
	var role entity.Role
	if err := ctx.BodyParser(&role); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}

	if err := c.usecase.CreateRole(ctx.Context(), &role); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusCreated,
		Message: "Role created successfully",
	})
}

// GetRoles retrieves roles for a school
// @Summary Get roles
// @Tags rbac
// @Param school_id query int false "School ID (optional)"
// @Success 200 {object} response.HTTPSuccessResponse{data=[]entity.Role}
// @Router /api/v1/rbac/roles [get]
func (c *RBACController) GetRoles(ctx *fiber.Ctx) error {
	var schoolID *int64
	idStr := ctx.Query("school_id")
	if idStr != "" {
		id, err := utils.ParseInt64(idStr)
		if err == nil {
			schoolID = &id
		}
	}

	roles, err := c.usecase.GetRolesBySchool(ctx.Context(), schoolID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Roles retrieved successfully",
		Data:    roles,
	})
}

// AssignPermission binds a permission to a role
// @Summary Assign permission to role
// @Tags rbac
// @Param id path int true "Role ID"
// @Param permission_id query int true "Permission ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/rbac/roles/{id}/permissions [post]
func (c *RBACController) AssignPermission(ctx *fiber.Ctx) error {
	roleID, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Role ID"})
	}

	permIDStr := ctx.Query("permission_id")
	permID, err := utils.ParseInt64(permIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Permission ID"})
	}

	if err := c.usecase.AssignPermissionToRole(ctx.Context(), roleID, permID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Permission assigned successfully",
	})
}

// AssignRoleToUser assigns a specific role to a user
// @Summary Assign role to user
// @Tags rbac
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param role_id body int true "Role ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/rbac/users/{user_id}/roles [post]
func (c *RBACController) AssignRoleToUser(ctx *fiber.Ctx) error {
	userID, err := utils.ParseInt64FromParam(ctx, "user_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid User ID"})
	}

	var body struct {
		RoleID int64 `json:"role_id"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}

	if err := c.usecase.AssignRoleToUser(ctx.Context(), userID, body.RoleID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Role assigned to user successfully",
	})
}

// GetUserRoles retrieves all roles assigned to a user
// @Summary Get user roles
// @Tags rbac
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/rbac/users/{user_id}/roles [get]
func (c *RBACController) GetUserRoles(ctx *fiber.Ctx) error {
	userID, err := utils.ParseInt64FromParam(ctx, "user_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid User ID"})
	}

	roles, err := c.usecase.GetUserRoles(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "User roles retrieved successfully",
		Data:    roles,
	})
}

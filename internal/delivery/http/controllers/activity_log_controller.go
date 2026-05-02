package controllers

import (
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/usecase"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ActivityLogController struct {
	Log     *logrus.Logger
	UseCase usecase.ActivityLogUsecase
}

func NewActivityLogController(useCase usecase.ActivityLogUsecase, logger *logrus.Logger) *ActivityLogController {
	return &ActivityLogController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *ActivityLogController) GetActivities(ctx *fiber.Ctx) error {
	params := utils.ParsePaginationQuery(ctx)
	
	schoolIDFromToken, _ := utils.GetSchoolIDFromToken(ctx)
	userType, _ := ctx.Locals("user_type").(string)

	var schoolID *int64
	if userType == "super_admin" {
		// Superadmin can filter by school_id or see all (null)
		if sID := ctx.QueryInt("school_id", 0); sID != 0 {
			val := int64(sID)
			schoolID = &val
		}
	} else {
		// Others only see their own school
		schoolID = &schoolIDFromToken
	}

	logs, err := c.UseCase.GetActivities(ctx.Context(), schoolID, params.Page, params.PageSize)
	if err != nil {
		c.Log.Errorf("Error getting activity logs: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve activity logs",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Activity logs retrieved successfully",
		Data:    logs,
	})
}

func (c *ActivityLogController) GetDeletions(ctx *fiber.Ctx) error {
	params := utils.ParsePaginationQuery(ctx)
	
	schoolIDFromToken, _ := utils.GetSchoolIDFromToken(ctx)
	userType, _ := ctx.Locals("user_type").(string)

	var schoolID *int64
	if userType == "super_admin" {
		// Superadmin can filter by school_id or see all (null)
		if sID := ctx.QueryInt("school_id", 0); sID != 0 {
			val := int64(sID)
			schoolID = &val
		}
	} else {
		// Others only see their own school
		schoolID = &schoolIDFromToken
	}

	logs, err := c.UseCase.GetDeletions(ctx.Context(), schoolID, params.Page, params.PageSize)
	if err != nil {
		c.Log.Errorf("Error getting deletion logs: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve deletion logs",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Deletion logs retrieved successfully",
		Data:    logs,
	})
}

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

type SettingController struct {
	Log     *logrus.Logger
	UseCase usecase.SettingUseCase
}

func NewSettingController(useCase usecase.SettingUseCase, logger *logrus.Logger) *SettingController {
	return &SettingController{
		Log:     logger,
		UseCase: useCase,
	}
}

// GetSettings handles getting settings (optionally filtered by group)
// @Summary Get settings
// @Description Get settings, optionally filtered by group_name
// @Tags settings
// @Accept json
// @Produce json
// @Param group query string false "Group Name (e.g. umum, backup, integrasi)"
// @Success 200 {object} response.HTTPSuccessResponse{data=[]entity.SettingResponse} "Settings retrieved successfully"
// @Failure 500 {object} response.HTTPErrorResponse "Failed to retrieve settings"
// @Router /api/v1/settings [get]
func (c *SettingController) GetSettings(ctx *fiber.Ctx) error {
	groupName := ctx.Query("group")
	
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)
	
	reqCtx := context.Background()
	var settings []*entity.Setting
	var err error

	if groupName != "" {
		settings, err = c.UseCase.GetSettingsByGroup(reqCtx, schoolID, groupName)
	} else {
		settings, err = c.UseCase.GetAllSettings(reqCtx, schoolID)
	}

	if err != nil {
		c.Log.Errorf("Error getting settings: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve settings",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	// Map to response DTO to hide audit fields
	res := []entity.SettingResponse{}
	for _, s := range settings {
		res = append(res, entity.SettingResponse{
			ID:           s.ID,
			SettingKey:   s.SettingKey,
			SettingValue: s.SettingValue,
			GroupName:    s.GroupName,
			Description:  s.Description,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Settings retrieved successfully",
		Data:    res,
	})
}

// UpdateSettings handles updating/inserting multiple settings
// @Summary Update settings
// @Description Update or insert multiple settings
// @Tags settings
// @Accept json
// @Produce json
// @Param request body entity.SettingUpdateRequest true "Settings update payload"
// @Success 200 {object} response.HTTPSuccessResponse "Settings updated successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid request body"
// @Failure 500 {object} response.HTTPErrorResponse "Failed to update settings"
// @Router /api/v1/settings [put]
func (c *SettingController) UpdateSettings(ctx *fiber.Ctx) error {
	var req entity.SettingUpdateRequest

	if err := ctx.BodyParser(&req); err != nil {
		c.Log.Errorf("Error parsing request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: err.Error()},
			},
		})
	}

	reqCtx := context.Background()
	updatedBy, _ := utils.GetUserIDFromToken(ctx)
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.UseCase.UpdateSettings(reqCtx, schoolID, req, updatedBy); err != nil {
		c.Log.Errorf("Error updating settings: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to update settings",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Settings updated successfully",
	})
}

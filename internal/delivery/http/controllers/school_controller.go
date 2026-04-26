package controllers

import (
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/usecase"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type SchoolController struct {
	usecase usecase.SchoolUsecase
	log     *logrus.Logger
}

func NewSchoolController(uc usecase.SchoolUsecase, log *logrus.Logger) *SchoolController {
	return &SchoolController{
		usecase: uc,
		log:     log,
	}
}

// RegisterSchool handles school onboarding
// @Summary Register a new school
// @Tags schools
// @Accept json
// @Produce json
// @Param school body entity.School true "School information"
// @Success 201 {object} response.HTTPSuccessResponse{data=entity.School}
// @Router /api/v1/schools [post]
func (c *SchoolController) RegisterSchool(ctx *fiber.Ctx) error {
	var school entity.School
	if err := ctx.BodyParser(&school); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	if err := c.usecase.RegisterSchool(ctx.Context(), &school); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusCreated,
		Message: "School registered successfully",
		Data:    school,
	})
}

// GetSchoolByID retrieves school details
// @Summary Get school by ID
// @Tags schools
// @Param id path int true "School ID"
// @Success 200 {object} response.HTTPSuccessResponse{data=entity.School}
// @Router /api/v1/schools/{id} [get]
func (c *SchoolController) GetSchoolByID(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid ID"})
	}

	school, err := c.usecase.GetSchoolByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.HTTPErrorResponse{Status: fiber.StatusNotFound, Message: "School not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "School retrieved successfully",
		Data:    school,
	})
}

// CreatePackage adds a new application package
// @Summary Create application package
// @Tags schools
// @Param package body entity.AppPackage true "Package details"
// @Success 201 {object} response.HTTPSuccessResponse
// @Router /api/v1/schools/packages [post]
func (c *SchoolController) CreatePackage(ctx *fiber.Ctx) error {
	var pkg entity.AppPackage
	if err := ctx.BodyParser(&pkg); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}

	if err := c.usecase.CreatePackage(ctx.Context(), &pkg); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusCreated,
		Message: "Package created successfully",
	})
}

// AssignLicense grants a package license to a school
// @Summary Assign license to school
// @Tags schools
// @Param id path int true "School ID"
// @Param package_code query string true "Package Code"
// @Success 200 {object} response.HTTPSuccessResponse{data=entity.SchoolLicense}
// @Router /api/v1/schools/{id}/license [post]
func (c *SchoolController) AssignLicense(ctx *fiber.Ctx) error {
	schoolID, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid ID"})
	}

	packageCode := ctx.Query("package_code")
	if packageCode == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "package_code is required"})
	}

	license, err := c.usecase.AssignLicense(ctx.Context(), schoolID, packageCode)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "License assigned successfully",
		Data:    license,
	})
}

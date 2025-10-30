package controllers

import (
	"fmt"
	"io"
	"strings"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/usecase"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ExcelController struct {
	excelUsecase usecase.ExcelUsecase
	validate     *validator.Validate
	log          *logrus.Logger
	db           *sqlx.DB
}

func NewExcelController(excelUsecase usecase.ExcelUsecase, validate *validator.Validate, log *logrus.Logger, db *sqlx.DB) *ExcelController {
	return &ExcelController{
		excelUsecase: excelUsecase,
		validate:     validate,
		log:          log,
		db:           db,
	}
}

// ImportExcel handles importing data from an Excel file
// @Summary Import data from Excel file
// @Description Import data from an uploaded Excel file
// @Tags excel
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel file to import"
// @Param sheet formData string false "Sheet name to import from" default(Sheet1)
// @Param table formData string true "Target table name"
// @Success 200 {object} response.HTTPSuccessResponse{data=entity.ExcelData} "Excel data imported successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid request or validation failed"
// @Failure 500 {object} response.HTTPErrorResponse "Excel import failed"
// @Router /api/v1/excel/import [post]
// @Security ApiKeyAuth
func (c *ExcelController) ImportExcel(ctx *fiber.Ctx) error {
	// Get form data
	form, err := ctx.MultipartForm()
	if err != nil {
		c.log.Errorf("Error parsing multipart form: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid form data",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: err.Error()},
			},
		})
	}

	// Get sheet name (optional, defaults to first sheet)
	var sheet string
	if sheetValues := form.Value["sheet"]; len(sheetValues) > 0 {
		sheet = sheetValues[0]
	}

	// Get table name (required)
	tableValues := form.Value["table"]
	if len(tableValues) == 0 {
		c.log.Error("Table name is required")
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Table name is required",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: "Table name is required"},
			},
		})
	}
	table := tableValues[0]

	// Validate table name
	if err := c.validate.Var(table, "required"); err != nil {
		c.log.Errorf("Validation error for table name: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: fmt.Sprintf("Table name validation failed: %v", err)},
			},
		})
	}

	// Get file
	files := form.File["file"]
	if len(files) == 0 {
		c.log.Error("File is required")
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "File is required",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: "File is required"},
			},
		})
	}
	file := files[0]

	// Check file extension
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".xlsx") &&
		!strings.HasSuffix(strings.ToLower(file.Filename), ".xls") {
		c.log.Error("Invalid file extension")
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid file type. Only .xlsx and .xls files are supported",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: "Invalid file type. Only .xlsx and .xls files are supported"},
			},
		})
	}

	// Open file
	fileContent, err := file.Open()
	if err != nil {
		c.log.Errorf("Error opening file: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to open file",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}
	defer fileContent.Close()

	// Read file content
	fileData, err := io.ReadAll(fileContent)
	if err != nil {
		c.log.Errorf("Error reading file: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to read file",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	// Create import request
	req := entity.ExcelImportRequest{
		FileData: fileData,
		Sheet:    sheet,
		Table:    table,
	}

	// Import Excel data
	excelData, err := c.excelUsecase.ImportExcel(req)
	if err != nil {
		c.log.Errorf("Excel import error: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Excel import failed",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Excel data imported successfully",
		Data:    excelData,
	})
}

// ExportExcel handles exporting data to an Excel file
// @Summary Export data to Excel file
// @Description Export data to an Excel file
// @Tags excel
// @Accept json
// @Produce json
// @Param request body entity.ExcelExportRequest true "Export request"
// @Success 200 {object} response.HTTPSuccessResponse{data=entity.ExcelResponse} "Excel file generated successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid request or validation failed"
// @Failure 500 {object} response.HTTPErrorResponse "Excel export failed"
// @Router /api/v1/excel/export [post]
// @Security ApiKeyAuth
func (c *ExcelController) ExportExcel(ctx *fiber.Ctx) error {
	var req entity.ExcelExportRequest

	// Parse request body
	if err := ctx.BodyParser(&req); err != nil {
		c.log.Errorf("Error parsing request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: err.Error()},
			},
		})
	}

	// Validate request
	if err := c.validate.Struct(req); err != nil {
		c.log.Errorf("Validation error: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: err.Error()},
			},
		})
	}

	// Export Excel data
	excelResponse, err := c.excelUsecase.ExportExcel(req)
	if err != nil {
		c.log.Errorf("Excel export error: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Excel export failed",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Excel file generated successfully",
		Data:    excelResponse,
	})
}

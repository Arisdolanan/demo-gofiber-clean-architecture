package controllers

import (
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/template/pdfs"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/usecase"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// TemplatePDFRequest represents the request structure for PDF generation from template
type TemplatePDFRequest struct {
	TemplateName string            `json:"template_name" validate:"required"`
	TemplateData pdfs.TemplateData `json:"template_data" validate:"required"`
	Filename     string            `json:"filename" validate:"required"`
	Title        string            `json:"title,omitempty"`
}

type PDFController struct {
	pdfUsecase usecase.PDFUsecase
	validate   *validator.Validate
	log        *logrus.Logger
}

func NewPDFController(pdfUsecase usecase.PDFUsecase, validate *validator.Validate, log *logrus.Logger) *PDFController {
	return &PDFController{
		pdfUsecase: pdfUsecase,
		validate:   validate,
		log:        log,
	}
}

// GeneratePDF handles PDF generation from HTML content
// @Summary Generate PDF from HTML content
// @Description Generate a PDF file from HTML content
// @Tags pdf
// @Accept json
// @Produce json
// @Param request body entity.PDFRequest true "PDF generation request"
// @Success 200 {object} response.HTTPSuccessResponse{data=entity.PDFResponse} "PDF generated successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid request body or validation failed"
// @Failure 500 {object} response.HTTPErrorResponse "PDF generation failed"
// @Router /api/v1/pdf/generate [post]
// @Security ApiKeyAuth
func (c *PDFController) GeneratePDF(ctx *fiber.Ctx) error {
	var req entity.PDFRequest

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

	// Generate PDF
	pdfResponse, err := c.pdfUsecase.GeneratePDF(req)
	if err != nil {
		c.log.Errorf("PDF generation error: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "PDF generation failed",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "PDF generated successfully",
		Data:    pdfResponse,
	})
}

// GeneratePDFFromTemplate handles PDF generation from a predefined template
// @Summary Generate PDF from template
// @Description Generate a PDF file from a predefined template
// @Tags pdf
// @Accept json
// @Produce json
// @Param request body TemplatePDFRequest true "PDF generation from template request"
// @Success 200 {object} response.HTTPSuccessResponse{data=entity.PDFResponse} "PDF generated successfully from template"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid request body or validation failed"
// @Failure 500 {object} response.HTTPErrorResponse "PDF generation failed"
// @Router /api/v1/pdf/generate-template [post]
// @Security ApiKeyAuth
func (c *PDFController) GeneratePDFFromTemplate(ctx *fiber.Ctx) error {
	var req TemplatePDFRequest

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

	// Generate PDF from template
	pdfResponse, err := c.pdfUsecase.GeneratePDFFromTemplate(req.TemplateName, req.TemplateData, req.Filename)
	if err != nil {
		c.log.Errorf("PDF generation from template error: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "PDF generation from template failed",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "PDF generated successfully from template",
		Data:    pdfResponse,
	})
}

package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/template/pdfs"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// Mock PDFUsecase for testing
type mockPDFUsecase struct{}

func (m *mockPDFUsecase) GeneratePDF(req entity.PDFRequest) (*entity.PDFResponse, error) {
	return &entity.PDFResponse{
		Filename: "test.pdf",
		URL:      "/storage/pdfs/test.pdf",
	}, nil
}

func (m *mockPDFUsecase) GeneratePDFFromTemplate(templateName string, templateData pdfs.TemplateData, filename string) (*entity.PDFResponse, error) {
	return &entity.PDFResponse{
		Filename: "test-template.pdf",
		URL:      "/storage/pdfs/test-template.pdf",
	}, nil
}

func TestPDFController_GeneratePDF(t *testing.T) {
	// Setup
	app := fiber.New()
	log := logrus.New()
	validate := validator.New()
	pdfUsecase := &mockPDFUsecase{}
	controller := NewPDFController(pdfUsecase, validate, log)

	// Register route
	app.Post("/api/v1/pdf/generate", controller.GeneratePDF)

	// Create test request
	req := entity.PDFRequest{
		HTMLContent: "<h1>Hello World</h1><p>This is a test PDF</p>",
		Filename:    "test.pdf",
		Title:       "Test PDF",
		Author:      "GoFiber Clean Architecture",
	}

	// Convert request to JSON
	reqBytes, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/pdf/generate", bytes.NewReader(reqBytes))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create response recorder
	resp, err := app.Test(httpReq)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Parse response
	// Note: We're not parsing the full response here to keep the test simple
	// In a real test, you would parse and verify the response body
	assert.NotNil(t, resp)
}

package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExcelUsecase is a mock implementation of the ExcelUsecase interface
type MockExcelUsecase struct {
	mock.Mock
}

func (m *MockExcelUsecase) ImportExcel(req entity.ExcelImportRequest) (*entity.ExcelData, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ExcelData), args.Error(1)
}

func (m *MockExcelUsecase) ExportExcel(req entity.ExcelExportRequest) (*entity.ExcelResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ExcelResponse), args.Error(1)
}

func (m *MockExcelUsecase) ExportCustomQueryToExcel(db *sqlx.DB, req entity.CustomQueryRequest) (*entity.ExcelResponse, error) {
	args := m.Called(db, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ExcelResponse), args.Error(1)
}

func (m *MockExcelUsecase) ExportQueryResultToExcel(queryResult *entity.QueryResult, sheet, filename string) (*entity.ExcelResponse, error) {
	args := m.Called(queryResult, sheet, filename)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ExcelResponse), args.Error(1)
}

// TestExcelController tests the Excel controller functionality
func TestExcelController(t *testing.T) {
	// Create a logger for testing
	log := logrus.New()

	// Create validator
	validate := validator.New()

	// Create mock usecase
	mockUsecase := new(MockExcelUsecase)

	// Create the Excel controller (db parameter can be nil for testing)
	excelController := NewExcelController(mockUsecase, validate, log, nil)

	// Test that the controller was created successfully
	assert.NotNil(t, excelController)
}

// TestExportExcelInvalidBody tests the ExportExcel endpoint with invalid request body
func TestExportExcelInvalidBody(t *testing.T) {
	// Create a logger for testing
	log := logrus.New()

	// Create validator
	validate := validator.New()

	// Create mock usecase
	mockUsecase := new(MockExcelUsecase)

	// Create the Excel controller (db parameter can be nil for testing)
	excelController := NewExcelController(mockUsecase, validate, log, nil)

	// Create a Fiber app for testing
	app := fiber.New()

	// Create a test route
	app.Post("/export", excelController.ExportExcel)

	// Create a test request with invalid JSON
	req := httptest.NewRequest("POST", "/export", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	// Create a test response recorder
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Check that we get a bad request response
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// TestExportExcelSuccess tests the ExportExcel endpoint with valid request
func TestExportExcelSuccess(t *testing.T) {
	// Create a logger for testing
	log := logrus.New()

	// Create validator
	validate := validator.New()

	// Create mock usecase
	mockUsecase := new(MockExcelUsecase)

	// Create the Excel controller (db parameter can be nil for testing)
	excelController := NewExcelController(mockUsecase, validate, log, nil)

	// Create a Fiber app for testing
	app := fiber.New()

	// Create a test route
	app.Post("/export", excelController.ExportExcel)

	// Create a test request body
	exportReq := entity.ExcelExportRequest{
		Data: []map[string]interface{}{
			{"name": "John", "age": 30},
			{"name": "Jane", "age": 25},
		},
		Sheet:    "Sheet1",
		Filename: "test.xlsx",
		Table:    "users",
	}

	// Marshal the request to JSON
	requestBody, err := json.Marshal(exportReq)
	assert.NoError(t, err)

	// Set up mock expectation
	mockResponse := &entity.ExcelResponse{
		Filename: "test.xlsx",
		URL:      "/storage/private/excels/test.xlsx",
		Rows:     2,
	}
	mockUsecase.On("ExportExcel", mock.AnythingOfType("entity.ExcelExportRequest")).Return(mockResponse, nil)

	// Create a test request
	req := httptest.NewRequest("POST", "/export", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a test response recorder
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Check that we get a success response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse the response
	var response response.HTTPSuccessResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	assert.NoError(t, err)

	// Check response data
	assert.Equal(t, fiber.StatusOK, response.Status)
	assert.Equal(t, "Excel file generated successfully", response.Message)

	// Verify mock expectations
	mockUsecase.AssertExpectations(t)
}

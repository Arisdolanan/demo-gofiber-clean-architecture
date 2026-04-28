package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// ClientErrorLogRequest represents the request body for client error logging
type ClientErrorLogRequest struct {
	Platform    string `json:"platform"`     // web, android, ios, macos, windows, linux
	AppVersion  string `json:"app_version"`  // App version
	Device      string `json:"device"`       // Device info
	ErrorType   string `json:"error_type"`   // Error type/category
	Message     string `json:"message"`      // Error message
	StackTrace  string `json:"stack_trace"`  // Stack trace (optional)
	Timestamp   string `json:"timestamp"`    // Client timestamp
	Endpoint    string `json:"endpoint"`     // API endpoint that failed (optional)
	RequestBody string `json:"request_body"` // Request body that caused error (optional)
}

// ClientLogController handles client-side error logging
type ClientLogController struct {
	log *logrus.Logger
}

// NewClientLogController creates a new instance of ClientLogController
func NewClientLogController(log *logrus.Logger) *ClientLogController {
	return &ClientLogController{
		log: log,
	}
}

// LogError handles error logs from frontend/mobile clients
// @Summary Log client-side error
// @Description Receive and store error logs from frontend/mobile applications
// @Tags client-logs
// @Accept json
// @Produce json
// @Param request body ClientErrorLogRequest true "Client error log data"
// @Success 200 {object} response.HTTPSuccessResponse "Error log recorded successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid request body"
// @Failure 500 {object} response.HTTPErrorResponse "Failed to record error log"
// @Router /api/v1/client-logs/error [post]
func (c *ClientLogController) LogError(ctx *fiber.Ctx) error {
	var req ClientErrorLogRequest

	// Parse request body
	if err := ctx.BodyParser(&req); err != nil {
		c.log.Errorf("Error parsing client error log request: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: err.Error()},
			},
		})
	}

	// Validate required fields
	if req.Platform == "" || req.Message == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Platform and message are required",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: "platform and message fields are required"},
			},
		})
	}

	// Set default timestamp if not provided
	if req.Timestamp == "" {
		req.Timestamp = time.Now().Format(time.RFC3339)
	}

	// Write error log to file
	if err := c.writeErrorLog(req); err != nil {
		c.log.Errorf("Failed to write client error log: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to record error log",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	c.log.Infof("Client error log received: [%s] %s - %s", req.Platform, req.ErrorType, req.Message)

	return ctx.JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Error log recorded successfully",
		Data:    nil,
	})
}

// writeErrorLog writes the error log to a platform-specific log file
func (c *ClientLogController) writeErrorLog(req ClientErrorLogRequest) error {
	// Create log directory if it doesn't exist
	logDir := "storage/logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Determine log file name based on platform
	var logFileName string
	switch req.Platform {
	case "web":
		logFileName = "web.error.log"
	case "android":
		logFileName = "android.error.log"
	case "ios":
		logFileName = "ios.error.log"
	case "macos", "darwin":
		logFileName = "mobile.error.log"
	case "windows":
		logFileName = "desktop.error.log"
	case "linux":
		logFileName = "desktop.error.log"
	default:
		logFileName = "client.error.log"
	}

	logFilePath := filepath.Join(logDir, logFileName)

	// Open log file in append mode, create if not exists
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Format log entry
	logEntry := fmt.Sprintf("[%s] Platform: %s | Version: %s | Device: %s\n"+
		"Error Type: %s\n"+
		"Message: %s\n"+
		"Endpoint: %s\n"+
		"Request Body: %s\n"+
		"Stack Trace:\n%s\n"+
		"----------------------------------------\n",
		req.Timestamp,
		req.Platform,
		req.AppVersion,
		req.Device,
		req.ErrorType,
		req.Message,
		req.Endpoint,
		req.RequestBody,
		req.StackTrace,
	)

	// Write to file
	if _, err := file.WriteString(logEntry); err != nil {
		return fmt.Errorf("failed to write log entry: %w", err)
	}

	return nil
}

// LogInfo handles info/warning logs from frontend/mobile clients
// @Summary Log client-side info/warning
// @Description Receive and store info/warning logs from frontend/mobile applications
// @Tags client-logs
// @Accept json
// @Produce json
// @Param request body ClientErrorLogRequest true "Client info log data"
// @Success 200 {object} response.HTTPSuccessResponse "Info log recorded successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid request body"
// @Failure 500 {object} response.HTTPErrorResponse "Failed to record info log"
// @Router /api/v1/client-logs/info [post]
func (c *ClientLogController) LogInfo(ctx *fiber.Ctx) error {
	var req ClientErrorLogRequest

	// Parse request body
	if err := ctx.BodyParser(&req); err != nil {
		c.log.Errorf("Error parsing client info log request: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: err.Error()},
			},
		})
	}

	// Validate required fields
	if req.Platform == "" || req.Message == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Platform and message are required",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: "platform and message fields are required"},
			},
		})
	}

	// Set default timestamp if not provided
	if req.Timestamp == "" {
		req.Timestamp = time.Now().Format(time.RFC3339)
	}

	// Write info log to file
	if err := c.writeInfoLog(req); err != nil {
		c.log.Errorf("Failed to write client info log: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to record info log",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	c.log.Debugf("Client info log received: [%s] %s", req.Platform, req.Message)

	return ctx.JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Info log recorded successfully",
		Data:    nil,
	})
}

// writeInfoLog writes the info log to a platform-specific log file
func (c *ClientLogController) writeInfoLog(req ClientErrorLogRequest) error {
	// Create log directory if it doesn't exist
	logDir := "storage/logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Determine log file name based on platform
	var logFileName string
	switch req.Platform {
	case "web":
		logFileName = "web.info.log"
	case "android":
		logFileName = "android.info.log"
	case "ios":
		logFileName = "ios.info.log"
	case "macos", "darwin":
		logFileName = "mobile.info.log"
	case "windows":
		logFileName = "desktop.info.log"
	case "linux":
		logFileName = "desktop.info.log"
	default:
		logFileName = "client.info.log"
	}

	logFilePath := filepath.Join(logDir, logFileName)

	// Open log file in append mode, create if not exists
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Format log entry
	logEntry := fmt.Sprintf("[%s] Platform: %s | Version: %s | Device: %s\n"+
		"Message: %s\n"+
		"Details: %s\n"+
		"----------------------------------------\n",
		req.Timestamp,
		req.Platform,
		req.AppVersion,
		req.Device,
		req.Message,
		req.StackTrace,
	)

	// Write to file
	if _, err := file.WriteString(logEntry); err != nil {
		return fmt.Errorf("failed to write log entry: %w", err)
	}

	return nil
}

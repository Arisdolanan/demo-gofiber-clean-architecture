package controllers

import (
	"context"
	"strconv"
	"strings"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/usecase"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type FileController struct {
	Log      *logrus.Logger
	UseCase  usecase.FileUseCase
	Validate *validator.Validate
}

func NewFileController(useCase usecase.FileUseCase, validate *validator.Validate, logger *logrus.Logger) *FileController {
	return &FileController{
		Log:      logger,
		UseCase:  useCase,
		Validate: validate,
	}
}

// UploadFile handles file upload
// @Summary Upload a file
// @Description Upload a file with optional metadata
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param description formData string false "File description"
// @Param category formData string false "File category"
// @Param is_public formData bool false "Whether file is public"
// @Success 201 {object} response.HTTPSuccessResponse{data=entity.File} "File uploaded successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid request or file"
// @Failure 401 {object} response.HTTPErrorResponse "Unauthorized"
// @Security ApiKeyAuth
// @Router /api/v1/files/upload [post]
func (c *FileController) UploadFile(ctx *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Errors:  []response.JSONError{{Status: fiber.StatusUnauthorized, Message: "Invalid token"}},
		})
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "No file uploaded",
			Errors:  []response.JSONError{{Status: fiber.StatusBadRequest, Message: err.Error()}},
		})
	}

	req := &entity.FileUploadRequest{
		Description: ctx.FormValue("description"),
		Category:    ctx.FormValue("category"),
		IsPublic:    ctx.FormValue("is_public") == "true",
	}

	if err := c.Validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Errors:  utils.ValidatorErrors(err),
		})
	}

	uploadedFile, err := c.UseCase.UploadFile(context.Background(), userID, file, req)
	if err != nil {
		c.Log.Errorf("Error uploading file: %v", err)
		if strings.Contains(err.Error(), "file type not allowed") || strings.Contains(err.Error(), "file size exceeds") {
			return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
				Status:  fiber.StatusBadRequest,
				Message: "File validation failed",
				Errors:  []response.JSONError{{Status: fiber.StatusBadRequest, Message: err.Error()}},
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to upload file",
			Errors:  []response.JSONError{{Status: fiber.StatusInternalServerError, Message: err.Error()}},
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusCreated,
		Message: "File uploaded successfully",
		Data:    uploadedFile,
	})
}

// GetUserFiles handles getting user's files with pagination
// @Summary Get user's files
// @Description Get files belonging to the authenticated user with pagination
// @Tags files
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param category query string false "Filter by category"
// @Success 200 {object} response.HTTPSuccessResponse{data=entity.FileListResponse} "Files retrieved successfully"
// @Failure 401 {object} response.HTTPErrorResponse "Unauthorized"
// @Security ApiKeyAuth
// @Router /api/v1/files [get]
func (c *FileController) GetUserFiles(ctx *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Errors:  []response.JSONError{{Status: fiber.StatusUnauthorized, Message: "Invalid token"}},
		})
	}

	// Parse pagination parameters using utils
	params := utils.ParsePaginationQuery(ctx)
	category := ctx.Query("category")

	var files *entity.FileListResponse
	if category != "" {
		files, err = c.UseCase.GetUserFilesByCategory(context.Background(), userID, category, params.Page, params.PageSize)
	} else {
		files, err = c.UseCase.GetUserFiles(context.Background(), userID, params.Page, params.PageSize)
	}

	if err != nil {
		c.Log.Errorf("Error getting user files: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve files",
			Errors:  []response.JSONError{{Status: fiber.StatusInternalServerError, Message: err.Error()}},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Files retrieved successfully",
		Data:    files,
	})
}

// GetPublicFiles handles getting public files with pagination
// @Summary Get public files
// @Description Get all public files with pagination
// @Tags files
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.HTTPSuccessResponse{data=entity.FileListResponse} "Public files retrieved successfully"
// @Router /api/v1/files/public [get]
func (c *FileController) GetPublicFiles(ctx *fiber.Ctx) error {
	// Parse pagination parameters using utils
	params := utils.ParsePaginationQuery(ctx)

	files, err := c.UseCase.GetPublicFiles(context.Background(), params.Page, params.PageSize)
	if err != nil {
		c.Log.Errorf("Error getting public files: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve public files",
			Errors:  []response.JSONError{{Status: fiber.StatusInternalServerError, Message: err.Error()}},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Public files retrieved successfully",
		Data:    files,
	})
}

// GetPrivateFiles handles getting user's private files with pagination
// @Summary Get user's private files
// @Description Get private files belonging to the authenticated user with pagination
// @Tags files
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.HTTPSuccessResponse{data=entity.FileListResponse} "Private files retrieved successfully"
// @Failure 401 {object} response.HTTPErrorResponse "Unauthorized"
// @Security ApiKeyAuth
// @Router /api/v1/files/private [get]
func (c *FileController) GetPrivateFiles(ctx *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Errors:  []response.JSONError{{Status: fiber.StatusUnauthorized, Message: "Invalid token"}},
		})
	}

	// Parse pagination parameters using utils
	params := utils.ParsePaginationQuery(ctx)

	files, err := c.UseCase.GetPrivateFilesByUserID(context.Background(), userID, params.Page, params.PageSize)
	if err != nil {
		c.Log.Errorf("Error getting private files: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve private files",
			Errors:  []response.JSONError{{Status: fiber.StatusInternalServerError, Message: err.Error()}},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Private files retrieved successfully",
		Data:    files,
	})
}

// UpdateFile handles file metadata update
// @Summary Update file metadata
// @Description Update file description, category, and public status
// @Tags files
// @Accept json
// @Produce json
// @Param id path int true "File ID"
// @Param file body entity.FileUpdateRequest true "File update information"
// @Success 200 {object} response.HTTPSuccessResponse{data=entity.File} "File updated successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid request"
// @Failure 401 {object} response.HTTPErrorResponse "Unauthorized"
// @Security ApiKeyAuth
// @Router /api/v1/files/{id} [put]
func (c *FileController) UpdateFile(ctx *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Errors:  []response.JSONError{{Status: fiber.StatusUnauthorized, Message: "Invalid token"}},
		})
	}

	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid file ID",
			Errors:  []response.JSONError{{Status: fiber.StatusBadRequest, Message: "Invalid file ID format"}},
		})
	}

	var req entity.FileUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Errors:  []response.JSONError{{Status: fiber.StatusBadRequest, Message: err.Error()}},
		})
	}

	if err := c.Validate.Struct(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Errors:  utils.ValidatorErrors(err),
		})
	}

	file, err := c.UseCase.UpdateFile(context.Background(), id, userID, &req)
	if err != nil {
		c.Log.Errorf("Error updating file %d: %v", id, err)
		if strings.Contains(err.Error(), "not found") {
			return ctx.Status(fiber.StatusNotFound).JSON(response.HTTPErrorResponse{
				Status:  fiber.StatusNotFound,
				Message: "File not found",
				Errors:  []response.JSONError{{Status: fiber.StatusNotFound, Message: err.Error()}},
			})
		}
		if strings.Contains(err.Error(), "access denied") {
			return ctx.Status(fiber.StatusForbidden).JSON(response.HTTPErrorResponse{
				Status:  fiber.StatusForbidden,
				Message: "Access denied",
				Errors:  []response.JSONError{{Status: fiber.StatusForbidden, Message: err.Error()}},
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to update file",
			Errors:  []response.JSONError{{Status: fiber.StatusInternalServerError, Message: err.Error()}},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "File updated successfully",
		Data:    file,
	})
}

// DeleteFile handles file deletion
// @Summary Delete file
// @Description Delete a file and its metadata
// @Tags files
// @Accept json
// @Produce json
// @Param id path int true "File ID"
// @Success 200 {object} response.HTTPSuccessResponse "File deleted successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid file ID"
// @Failure 401 {object} response.HTTPErrorResponse "Unauthorized"
// @Security ApiKeyAuth
// @Router /api/v1/files/{id} [delete]
func (c *FileController) DeleteFile(ctx *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Errors:  []response.JSONError{{Status: fiber.StatusUnauthorized, Message: "Invalid token"}},
		})
	}

	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid file ID",
			Errors:  []response.JSONError{{Status: fiber.StatusBadRequest, Message: "Invalid file ID format"}},
		})
	}

	err = c.UseCase.DeleteFile(context.Background(), id, userID)
	if err != nil {
		c.Log.Errorf("Error deleting file %d: %v", id, err)
		if strings.Contains(err.Error(), "not found") {
			return ctx.Status(fiber.StatusNotFound).JSON(response.HTTPErrorResponse{
				Status:  fiber.StatusNotFound,
				Message: "File not found",
				Errors:  []response.JSONError{{Status: fiber.StatusNotFound, Message: err.Error()}},
			})
		}
		if strings.Contains(err.Error(), "access denied") {
			return ctx.Status(fiber.StatusForbidden).JSON(response.HTTPErrorResponse{
				Status:  fiber.StatusForbidden,
				Message: "Access denied",
				Errors:  []response.JSONError{{Status: fiber.StatusForbidden, Message: err.Error()}},
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to delete file",
			Errors:  []response.JSONError{{Status: fiber.StatusInternalServerError, Message: err.Error()}},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "File deleted successfully",
		Data:    nil,
	})
}

// DownloadFile handles file download
// @Summary Download file
// @Description Download a file by ID
// @Tags files
// @Accept json
// @Produce application/octet-stream
// @Param id path int true "File ID"
// @Success 200 {file} file "File content"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid file ID"
// @Failure 401 {object} response.HTTPErrorResponse "Unauthorized"
// @Security ApiKeyAuth
// @Router /api/v1/files/{id}/download [get]
func (c *FileController) DownloadFile(ctx *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Errors:  []response.JSONError{{Status: fiber.StatusUnauthorized, Message: "Invalid token"}},
		})
	}

	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid file ID",
			Errors:  []response.JSONError{{Status: fiber.StatusBadRequest, Message: "Invalid file ID format"}},
		})
	}

	fileResponse, err := c.UseCase.DownloadFile(context.Background(), id, userID)
	if err != nil {
		c.Log.Errorf("Error downloading file %d: %v", id, err)
		if strings.Contains(err.Error(), "not found") {
			return ctx.Status(fiber.StatusNotFound).JSON(response.HTTPErrorResponse{
				Status:  fiber.StatusNotFound,
				Message: "File not found",
				Errors:  []response.JSONError{{Status: fiber.StatusNotFound, Message: err.Error()}},
			})
		}
		if strings.Contains(err.Error(), "access denied") {
			return ctx.Status(fiber.StatusForbidden).JSON(response.HTTPErrorResponse{
				Status:  fiber.StatusForbidden,
				Message: "Access denied",
				Errors:  []response.JSONError{{Status: fiber.StatusForbidden, Message: err.Error()}},
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to download file",
			Errors:  []response.JSONError{{Status: fiber.StatusInternalServerError, Message: err.Error()}},
		})
	}

	ctx.Set("Content-Type", fileResponse.MimeType)
	ctx.Set("Content-Disposition", `attachment; filename="`+fileResponse.Filename+`"`)
	ctx.Set("Content-Length", strconv.FormatInt(fileResponse.Size, 10))

	return ctx.Send(fileResponse.Data)
}

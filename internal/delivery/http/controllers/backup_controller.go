package controllers

import (
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/usecase"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type BackupController struct {
	Log     *logrus.Logger
	UseCase usecase.BackupUseCase
}

func NewBackupController(useCase usecase.BackupUseCase, logger *logrus.Logger) *BackupController {
	return &BackupController{
		Log:     logger,
		UseCase: useCase,
	}
}

// GetBackups handles listing backup files
// @Summary List backups
// @Description List all available database backup files
// @Tags backup
// @Accept json
// @Produce json
// @Success 200 {object} response.HTTPSuccessResponse{data=[]entity.BackupRecord} "Backups retrieved successfully"
// @Failure 500 {object} response.HTTPErrorResponse "Failed to retrieve backups"
// @Router /api/v1/backup/list [get]
func (c *BackupController) GetBackups(ctx *fiber.Ctx) error {
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)
	backups, err := c.UseCase.ListBackups(ctx.Context(), schoolID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve backups",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Backups retrieved successfully",
		Data:    backups,
	})
}

// CreateManualBackup handles triggering a manual backup
// @Summary Create backup
// @Description Trigger a manual database backup
// @Tags backup
// @Accept json
// @Produce json
// @Success 200 {object} response.HTTPSuccessResponse{data=entity.BackupRecord} "Backup created successfully"
// @Failure 500 {object} response.HTTPErrorResponse "Failed to create backup"
// @Router /api/v1/backup/manual [post]
func (c *BackupController) CreateManualBackup(ctx *fiber.Ctx) error {
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)
	info, err := c.UseCase.CreateBackup(ctx.Context(), schoolID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to create backup",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Backup created successfully",
		Data:    info,
	})
}

// RestoreBackup handles restoring from a backup file
// @Summary Restore backup
// @Description Restore database from a specified backup file
// @Tags backup
// @Accept json
// @Produce json
// @Param request body entity.RestoreRequest true "Restore request payload"
// @Success 200 {object} response.HTTPSuccessResponse "Database restored successfully"
// @Failure 400 {object} response.HTTPErrorResponse "Invalid request body"
// @Failure 500 {object} response.HTTPErrorResponse "Failed to restore backup"
// @Router /api/v1/backup/restore [post]
func (c *BackupController) RestoreBackup(ctx *fiber.Ctx) error {
	var req entity.RestoreRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Errors: []response.JSONError{
				{Status: fiber.StatusBadRequest, Message: err.Error()},
			},
		})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.UseCase.RestoreBackup(ctx.Context(), schoolID, req.Filename); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to restore backup",
			Errors: []response.JSONError{
				{Status: fiber.StatusInternalServerError, Message: err.Error()},
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Database restored successfully",
	})
}

package controllers

import (
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/usecase"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AcademicController struct {
	usecase usecase.AcademicUsecase
	log     *logrus.Logger
}

func NewAcademicController(uc usecase.AcademicUsecase, log *logrus.Logger) *AcademicController {
	return &AcademicController{
		usecase: uc,
		log:     log,
	}
}

// Sessions

// CreateSession handles academic session creation
// @Summary Create academic session
// @Tags academic
// @Accept json
// @Produce json
// @Param session body entity.AcademicSession true "Session details"
// @Success 201 {object} response.HTTPSuccessResponse
// @Router /api/v1/academic/sessions [post]
func (c *AcademicController) CreateSession(ctx *fiber.Ctx) error {
	var session entity.AcademicSession
	if err := ctx.BodyParser(&session); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.usecase.CreateSession(ctx.Context(), schoolID, &session); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.HTTPSuccessResponse{Status: fiber.StatusCreated, Message: "Session created successfully"})
}

// UpdateSession handles academic session updates
// @Summary Update academic session
// @Tags academic
// @Accept json
// @Produce json
// @Param id path int true "Session ID"
// @Param session body entity.AcademicSession true "Session details"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/academic/sessions/{id} [put]
func (c *AcademicController) UpdateSession(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Session ID"})
	}

	var session entity.AcademicSession
	if err := ctx.BodyParser(&session); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)
	session.ID = id

	if err := c.usecase.UpdateSession(ctx.Context(), schoolID, &session); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Session updated successfully"})
}

// DeleteSession handles academic session deletion
// @Summary Delete academic session
// @Tags academic
// @Param id path int true "Session ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/academic/sessions/{id} [delete]
func (c *AcademicController) DeleteSession(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Session ID"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.usecase.DeleteSession(ctx.Context(), schoolID, id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Session deleted successfully"})
}

// GetSessions retrieves all sessions for a school
// @Summary Get sessions
// @Tags academic
// @Param school_id query int true "School ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/academic/sessions [get]
func (c *AcademicController) GetSessions(ctx *fiber.Ctx) error {
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	sessions, err := c.usecase.GetSessionsBySchool(ctx.Context(), schoolID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Sessions retrieved successfully", Data: sessions})
}

// GetActiveSession retrieves the active academic session
// @Summary Get active session
// @Tags academic
// @Param school_id query int true "School ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/academic/sessions/active [get]
func (c *AcademicController) GetActiveSession(ctx *fiber.Ctx) error {
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	session, err := c.usecase.GetActiveSession(ctx.Context(), schoolID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Active session retrieved successfully", Data: session})
}

// Classes

// CreateClass handles class creation
// @Summary Create school class
// @Tags academic
// @Accept json
// @Produce json
// @Param class body entity.Class true "Class details"
// @Success 201 {object} response.HTTPSuccessResponse
// @Router /api/v1/academic/classes [post]
func (c *AcademicController) CreateClass(ctx *fiber.Ctx) error {
	var class entity.Class
	if err := ctx.BodyParser(&class); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.usecase.CreateClass(ctx.Context(), schoolID, &class); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.HTTPSuccessResponse{Status: fiber.StatusCreated, Message: "Class created successfully"})
}

// UpdateClass handles class updates
// @Summary Update school class
// @Tags academic
// @Accept json
// @Produce json
// @Param id path int true "Class ID"
// @Param class body entity.Class true "Class details"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/academic/classes/{id} [put]
func (c *AcademicController) UpdateClass(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Class ID"})
	}

	var class entity.Class
	if err := ctx.BodyParser(&class); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)
	class.ID = id

	if err := c.usecase.UpdateClass(ctx.Context(), schoolID, &class); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Class updated successfully"})
}

// DeleteClass handles class deletion
// @Summary Delete school class
// @Tags academic
// @Param id path int true "Class ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/academic/classes/{id} [delete]
func (c *AcademicController) DeleteClass(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Class ID"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.usecase.DeleteClass(ctx.Context(), schoolID, id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Class deleted successfully"})
}

// GetClasses retrieves all classes for a school
// @Summary Get classes
// @Tags academic
// @Param school_id query int true "School ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/academic/classes [get]
func (c *AcademicController) GetClasses(ctx *fiber.Ctx) error {
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	classes, err := c.usecase.GetClassesBySchool(ctx.Context(), schoolID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Classes retrieved successfully", Data: classes})
}

// Sections

// CreateSection handles section creation
// @Summary Create class section
// @Tags academic
// @Accept json
// @Produce json
// @Param section body entity.Section true "Section details"
// @Success 201 {object} response.HTTPSuccessResponse
// @Router /api/v1/academic/sections [post]
func (c *AcademicController) CreateSection(ctx *fiber.Ctx) error {
	var section entity.Section
	if err := ctx.BodyParser(&section); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}

	// Optionally verify class belongs to school or set school_id if section had it
	// For now, sections are class-based, and we already filter by class/school in repo
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)
	if err := c.usecase.CreateSection(ctx.Context(), schoolID, &section); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.HTTPSuccessResponse{Status: fiber.StatusCreated, Message: "Section created successfully"})
}

// UpdateSection handles section updates
// @Summary Update class section
// @Tags academic
// @Accept json
// @Produce json
// @Param id path int true "Section ID"
// @Param section body entity.Section true "Section details"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/academic/sections/{id} [put]
func (c *AcademicController) UpdateSection(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Section ID"})
	}

	var section entity.Section
	if err := ctx.BodyParser(&section); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}
	section.ID = id
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)
	if err := c.usecase.UpdateSection(ctx.Context(), schoolID, &section); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Section updated successfully"})
}

// DeleteSection handles section deletion
// @Summary Delete class section
// @Tags academic
// @Param id path int true "Section ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/academic/sections/{id} [delete]
func (c *AcademicController) DeleteSection(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Section ID"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.usecase.DeleteSection(ctx.Context(), schoolID, id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Section deleted successfully"})
}

// GetSections retrieves all sections for a class
// @Summary Get sections by class
// @Tags academic
// @Param class_id query int true "Class ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/academic/sections [get]
func (c *AcademicController) GetSections(ctx *fiber.Ctx) error {
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)
	classIDStr := ctx.Query("class_id")
	sessionIDStr := ctx.Query("academic_session_id")

	// If no filters are provided, return all sections
	if classIDStr == "" && sessionIDStr == "" {
		sections, err := c.usecase.GetAllSections(ctx.Context(), schoolID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
		}
		return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Sections retrieved successfully", Data: sections})
	}

	if sessionIDStr != "" {
		sessionID, err := utils.ParseInt64(sessionIDStr)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Session ID"})
		}
		sections, err := c.usecase.GetSectionsBySession(ctx.Context(), schoolID, sessionID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
		}
		return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Sections retrieved successfully", Data: sections})
	}

	classID, err := utils.ParseInt64(classIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Class ID"})
	}

	sections, err := c.usecase.GetSectionsByClass(ctx.Context(), schoolID, classID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Sections retrieved successfully", Data: sections})
}

// Subjects

// CreateSubject handles subject creation
// @Summary Create school subject
// @Tags academic
// @Accept json
// @Produce json
// @Param subject body entity.Subject true "Subject details"
// @Success 201 {object} response.HTTPSuccessResponse
// @Router /api/v1/academic/subjects [post]
func (c *AcademicController) CreateSubject(ctx *fiber.Ctx) error {
	var subject entity.Subject
	if err := ctx.BodyParser(&subject); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.usecase.CreateSubject(ctx.Context(), schoolID, &subject); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.HTTPSuccessResponse{Status: fiber.StatusCreated, Message: "Subject created successfully"})
}

// UpdateSubject handles subject updates
// @Summary Update school subject
// @Tags academic
// @Accept json
// @Produce json
// @Param id path int true "Subject ID"
// @Param subject body entity.Subject true "Subject details"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/academic/subjects/{id} [put]
func (c *AcademicController) UpdateSubject(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Subject ID"})
	}

	var subject entity.Subject
	if err := ctx.BodyParser(&subject); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)
	subject.ID = id

	if err := c.usecase.UpdateSubject(ctx.Context(), schoolID, &subject); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Subject updated successfully"})
}

// DeleteSubject handles subject deletion
// @Summary Delete school subject
// @Tags academic
// @Param id path int true "Subject ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/academic/subjects/{id} [delete]
func (c *AcademicController) DeleteSubject(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Subject ID"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.usecase.DeleteSubject(ctx.Context(), schoolID, id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Subject deleted successfully"})
}

// GetSubjects retrieves all subjects for a school
// @Summary Get subjects
// @Tags academic
// @Param school_id query int true "School ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/academic/subjects [get]
func (c *AcademicController) GetSubjects(ctx *fiber.Ctx) error {
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	subjects, err := c.usecase.GetSubjectsBySchool(ctx.Context(), schoolID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Subjects retrieved successfully", Data: subjects})
}

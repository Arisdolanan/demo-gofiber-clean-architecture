package controllers

import (
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/usecase"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type OperationController struct {
	usecase usecase.OperationUsecase
	log     *logrus.Logger
}

func NewOperationController(uc usecase.OperationUsecase, log *logrus.Logger) *OperationController {
	return &OperationController{
		usecase: uc,
		log:     log,
	}
}

// CreateSchedule handles schedule creation
// @Summary Create class schedule
// @Tags operations
// @Accept json
// @Produce json
// @Param schedule body entity.Schedule true "Schedule details"
// @Success 201 {object} response.HTTPSuccessResponse
// @Router /api/v1/operations/schedules [post]
func (c *OperationController) CreateSchedule(ctx *fiber.Ctx) error {
	var schedule entity.Schedule
	if err := ctx.BodyParser(&schedule); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}

	if err := c.usecase.CreateSchedule(ctx.Context(), &schedule); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.HTTPSuccessResponse{Status: fiber.StatusCreated, Message: "Schedule created successfully"})
}

// UpdateSchedule handles schedule updates
// @Summary Update class schedule
// @Tags operations
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param schedule body entity.Schedule true "Schedule details"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/operations/schedules/{id} [put]
func (c *OperationController) UpdateSchedule(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Schedule ID"})
	}

	var schedule entity.Schedule
	if err := ctx.BodyParser(&schedule); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}
	schedule.ID = id

	if err := c.usecase.UpdateSchedule(ctx.Context(), &schedule); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Schedule updated successfully"})
}

// DeleteSchedule handles schedule deletion
// @Summary Delete class schedule
// @Tags operations
// @Param id path int true "Schedule ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/operations/schedules/{id} [delete]
func (c *OperationController) DeleteSchedule(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Schedule ID"})
	}

	if err := c.usecase.DeleteSchedule(ctx.Context(), id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Schedule deleted successfully"})
}

// GetAllSchedules retrieves all schedules with filters
// @Summary List all schedules
// @Tags operations
// @Produce json
// @Param section_id query int false "Section ID"
// @Param teacher_id query int false "Teacher ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/operations/schedules [get]
func (c *OperationController) GetAllSchedules(ctx *fiber.Ctx) error {
	filters := make(map[string]interface{})
	if sid := ctx.Query("section_id"); sid != "" {
		filters["section_id"] = sid
	}
	if tid := ctx.Query("teacher_id"); tid != "" {
		filters["teacher_id"] = tid
	}

	schedules, err := c.usecase.GetAllSchedules(ctx.Context(), filters)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Schedules retrieved successfully", Data: schedules})
}

// CreateExam handles exam creation
// @Summary Create a new exam
// @Tags operations
// @Accept json
// @Produce json
// @Param exam body entity.Exam true "Exam details"
// @Success 201 {object} response.HTTPSuccessResponse
// @Router /api/v1/operations/exams [post]
func (c *OperationController) CreateExam(ctx *fiber.Ctx) error {
	var exam entity.Exam
	if err := ctx.BodyParser(&exam); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}

	if err := c.usecase.CreateExam(ctx.Context(), &exam); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.HTTPSuccessResponse{Status: fiber.StatusCreated, Message: "Exam created successfully"})
}

// UpdateExam handles exam updates
// @Summary Update exam details
// @Tags operations
// @Accept json
// @Produce json
// @Param id path int true "Exam ID"
// @Param exam body entity.Exam true "Exam details"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/operations/exams/{id} [put]
func (c *OperationController) UpdateExam(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Exam ID"})
	}

	var exam entity.Exam
	if err := ctx.BodyParser(&exam); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}
	exam.ID = id

	if err := c.usecase.UpdateExam(ctx.Context(), &exam); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Exam updated successfully"})
}

// DeleteExam handles exam deletion
// @Summary Delete an exam
// @Tags operations
// @Param id path int true "Exam ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/operations/exams/{id} [delete]
func (c *OperationController) DeleteExam(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Exam ID"})
	}

	if err := c.usecase.DeleteExam(ctx.Context(), id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Exam deleted successfully"})
}

// GetAllExams retrieves all exams with filters
// @Summary List all exams
// @Tags operations
// @Produce json
// @Param section_id query int false "Section ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/operations/exams [get]
func (c *OperationController) GetAllExams(ctx *fiber.Ctx) error {
	filters := make(map[string]interface{})
	if sid := ctx.Query("section_id"); sid != "" {
		filters["section_id"] = sid
	}

	exams, err := c.usecase.GetAllExams(ctx.Context(), filters)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Exams retrieved successfully", Data: exams})
}

// UpdateMark handles adding/updating exam marks
// @Summary Update student exam mark
// @Tags operations
// @Accept json
// @Produce json
// @Param mark body entity.ExamMark true "Mark details"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/operations/marks [post]
func (c *OperationController) UpdateMark(ctx *fiber.Ctx) error {
	var mark entity.ExamMark
	if err := ctx.BodyParser(&mark); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}

	if err := c.usecase.UpdateMark(ctx.Context(), &mark); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Mark updated successfully"})
}

// GetExamMarks retrieves marks for a specific exam
// @Summary Get marks by Exam
// @Tags operations
// @Produce json
// @Param exam_id path int true "Exam ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/operations/exams/{exam_id}/marks [get]
func (c *OperationController) GetExamMarks(ctx *fiber.Ctx) error {
	examID, err := utils.ParseInt64FromParam(ctx, "exam_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Exam ID"})
	}

	marks, err := c.usecase.GetExamMarks(ctx.Context(), examID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Exam marks retrieved successfully", Data: marks})
}

// GetStudentMarks retrieves marks for a specific student
// @Summary Get marks by Student
// @Tags operations
// @Produce json
// @Param student_id path int true "Student ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/operations/students/{student_id}/marks [get]
func (c *OperationController) GetStudentMarks(ctx *fiber.Ctx) error {
	studentID, err := utils.ParseInt64FromParam(ctx, "student_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Student ID"})
	}

	marks, err := c.usecase.GetStudentMarks(ctx.Context(), studentID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Student marks retrieved successfully", Data: marks})
}

// RecordStudentAttendance handles student attendance logging
// @Summary Record student attendance
// @Tags operations
// @Accept json
// @Produce json
// @Param attendance body entity.StudentAttendance true "Attendance details"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/operations/attendance/student [post]
func (c *OperationController) RecordStudentAttendance(ctx *fiber.Ctx) error {
	var attendance entity.StudentAttendance
	if err := ctx.BodyParser(&attendance); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}

	if err := c.usecase.RecordStudentAttendance(ctx.Context(), &attendance); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Student attendance recorded successfully"})
}

// UpdateStudentAttendance handles student attendance updates
// @Summary Update student attendance
// @Tags operations
// @Accept json
// @Produce json
// @Param id path int true "Attendance ID"
// @Param attendance body entity.StudentAttendance true "Attendance details"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/operations/attendance/student/{id} [put]
func (c *OperationController) UpdateStudentAttendance(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Attendance ID"})
	}

	var attendance entity.StudentAttendance
	if err := ctx.BodyParser(&attendance); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}
	attendance.ID = id

	if err := c.usecase.UpdateStudentAttendance(ctx.Context(), &attendance); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Student attendance updated successfully"})
}

// GetStudentAttendance retrieves attendance history for a student
// @Summary Get student attendance history
// @Tags operations
// @Produce json
// @Param student_id path int true "Student ID"
// @Param start_date query string false "Start Date (YYYY-MM-DD)"
// @Param end_date query string false "End Date (YYYY-MM-DD)"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/operations/students/{student_id}/attendance [get]
func (c *OperationController) GetStudentAttendance(ctx *fiber.Ctx) error {
	studentID, err := utils.ParseInt64FromParam(ctx, "student_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Student ID"})
	}

	filters := make(map[string]interface{})
	if sd := ctx.Query("start_date"); sd != "" {
		filters["start_date"] = sd
	}
	if ed := ctx.Query("end_date"); ed != "" {
		filters["end_date"] = ed
	}

	attendance, err := c.usecase.GetStudentAttendance(ctx.Context(), studentID, filters)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Student attendance retrieved successfully", Data: attendance})
}

// GetAttendanceReport retrieves filtered attendance records
// @Summary Get filtered attendance report
// @Tags operations
// @Produce json
// @Param section_id query int false "Section ID"
// @Param date query string false "Date (YYYY-MM-DD)"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/operations/attendance [get]
func (c *OperationController) GetAttendanceReport(ctx *fiber.Ctx) error {
	sectionIDStr := ctx.Query("section_id")
	date := ctx.Query("date")

	if sectionIDStr != "" && date != "" {
		sectionID, _ := utils.ParseInt64(sectionIDStr)
		attendance, err := c.usecase.GetAttendanceBySection(ctx.Context(), sectionID, date)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
		}
		return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Section attendance retrieved successfully", Data: attendance})
	}

	filters := make(map[string]interface{})
	if sectionIDStr != "" {
		filters["section_id"] = sectionIDStr
	}
	if date != "" {
		filters["date"] = date
	}

	attendance, err := c.usecase.GetAttendanceWithFilters(ctx.Context(), filters)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Attendance report retrieved successfully", Data: attendance})
}

// RecordTeacherAttendance handles teacher attendance logging
// @Summary Record teacher attendance
// @Tags operations
// @Accept json
// @Produce json
// @Param attendance body entity.TeacherAttendance true "Attendance details"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/operations/attendance/teacher [post]
func (c *OperationController) RecordTeacherAttendance(ctx *fiber.Ctx) error {
	var attendance entity.TeacherAttendance
	if err := ctx.BodyParser(&attendance); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}

	if err := c.usecase.RecordTeacherAttendance(ctx.Context(), &attendance); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Teacher attendance recorded successfully"})
}

// GetNotifications retrieves notifications for the current user
// @Summary Get user notifications
// @Tags operations
// @Produce json
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/operations/notifications [get]
func (c *OperationController) GetNotifications(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(int64)
	
	notifications, err := c.usecase.GetUserNotifications(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Notifications retrieved successfully", Data: notifications})
}

// GetReportCard handles fetching student report card
// @Summary Get student report card
// @Description Get consolidated report card for a student in a specific session
// @Tags operations
// @Accept json
// @Produce json
// @Param student_id query int true "Student ID"
// @Param session_id query int true "Session ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Failure 500 {object} response.HTTPErrorResponse
// @Router /api/v1/operations/report-card [get]
func (c *OperationController) GetReportCard(ctx *fiber.Ctx) error {
	studentID, _ := utils.ParseInt64(ctx.Query("student_id"))
	sessionID, _ := utils.ParseInt64(ctx.Query("session_id"))

	report, err := c.usecase.GetStudentReportCard(ctx.Context(), studentID, sessionID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Report card retrieved successfully", Data: report})
}

// GetSectionReportCards handles fetching bulk report cards for a section
// @Summary Get section report cards
// @Description Get report cards for all students in a section for a specific session
// @Tags operations
// @Accept json
// @Produce json
// @Param section_id query int true "Section ID"
// @Param session_id query int true "Session ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Failure 500 {object} response.HTTPErrorResponse
// @Router /api/v1/operations/report-card/section [get]
func (c *OperationController) GetSectionReportCards(ctx *fiber.Ctx) error {
	sectionID, _ := utils.ParseInt64(ctx.Query("section_id"))
	sessionID, _ := utils.ParseInt64(ctx.Query("session_id"))

	reports, err := c.usecase.GetSectionReportCards(ctx.Context(), sectionID, sessionID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Class report cards retrieved successfully", Data: reports})
}

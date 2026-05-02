package controllers

import (
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/usecase"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type PeopleController struct {
	usecase usecase.PeopleUsecase
	log     *logrus.Logger
}

func NewPeopleController(uc usecase.PeopleUsecase, log *logrus.Logger) *PeopleController {
	return &PeopleController{
		usecase: uc,
		log:     log,
	}
}

// CreateTeacher handles teacher onboarding
// @Summary Create school teacher
// @Tags people
// @Accept json
// @Produce json
// @Param teacher body entity.Teacher true "Teacher details"
// @Success 201 {object} response.HTTPSuccessResponse
// @Failure 400 {object} response.HTTPErrorResponse
// @Router /api/v1/people/teachers [post]
func (c *PeopleController) CreateTeacher(ctx *fiber.Ctx) error {
	var teacher entity.Teacher
	if err := ctx.BodyParser(&teacher); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)
	resp, err := c.usecase.CreateTeacher(ctx.Context(), schoolID, &teacher)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusCreated,
		Message: "Teacher created successfully",
		Data:    resp,
	})
}

// UpdateTeacher handles teacher updates
// @Summary Update school teacher
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "Teacher ID"
// @Param teacher body entity.Teacher true "Teacher details"
// @Success 200 {object} response.HTTPSuccessResponse
// @Failure 400 {object} response.HTTPErrorResponse
// @Router /api/v1/people/teachers/{id} [put]
func (c *PeopleController) UpdateTeacher(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Teacher ID"})
	}

	var teacher entity.Teacher
	if err := ctx.BodyParser(&teacher); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}
	teacher.ID = id

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.usecase.UpdateTeacher(ctx.Context(), schoolID, &teacher); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Teacher updated successfully"})
}

// GetTeachers retrieves all teachers for a school
// @Summary Get teachers
// @Tags people
// @Produce json
// @Param school_id query int true "School ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/people/teachers [get]
func (c *PeopleController) GetTeachers(ctx *fiber.Ctx) error {
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	teachers, err := c.usecase.GetTeachersBySchool(ctx.Context(), schoolID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Teachers retrieved successfully", Data: teachers})
}

// GetTeacherByUserID retrieves teacher details by user ID
// @Summary Get teacher by User ID
// @Tags people
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/people/teachers/users/{user_id} [get]
func (c *PeopleController) GetTeacherByUserID(ctx *fiber.Ctx) error {
	userID, err := utils.ParseInt64FromParam(ctx, "user_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid User ID"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	teacher, err := c.usecase.GetTeacherByUserID(ctx.Context(), schoolID, userID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.HTTPErrorResponse{Status: fiber.StatusNotFound, Message: "Teacher not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Teacher retrieved successfully", Data: teacher})
}

// DeleteTeacher handles teacher deletion
// @Summary Delete school teacher
// @Tags people
// @Produce json
// @Param id path int true "Teacher ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/people/teachers/{id} [delete]
func (c *PeopleController) DeleteTeacher(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Teacher ID"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.usecase.DeleteTeacher(ctx.Context(), schoolID, id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Teacher deleted and user deactivated successfully"})
}

// CreateStudent handles student registration
// @Summary Create school student
// @Tags people
// @Accept json
// @Produce json
// @Param student body entity.Student true "Student details (can include parents for unified registration)"
// @Success 201 {object} response.HTTPSuccessResponse
// @Router /api/v1/people/students [post]
func (c *PeopleController) CreateStudent(ctx *fiber.Ctx) error {
	var student entity.Student
	if err := ctx.BodyParser(&student); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)
	resp, err := c.usecase.CreateStudent(ctx.Context(), schoolID, &student)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusCreated,
		Message: "Student registered successfully",
		Data:    resp,
	})
}

// UpdateStudent handles student updates
// @Summary Update school student
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param student body entity.Student true "Student details"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/people/students/{id} [put]
func (c *PeopleController) UpdateStudent(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Student ID"})
	}

	var student entity.Student
	if err := ctx.BodyParser(&student); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}
	student.ID = id

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.usecase.UpdateStudent(ctx.Context(), schoolID, &student); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Student updated successfully"})
}

// GetAllStudents retrieves all students with pagination and filters
// @Summary List all students
// @Tags people
// @Produce json
// @Param school_id query int true "School ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/people/students/list [get]
func (c *PeopleController) GetAllStudents(ctx *fiber.Ctx) error {
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	limit := ctx.QueryInt("limit", 10)
	offset := ctx.QueryInt("offset", 0)

	students, total, err := c.usecase.GetAllStudents(ctx.Context(), schoolID, limit, offset, nil)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Students retrieved successfully",
		Data: fiber.Map{
			"items":  students,
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetStudentsBySection retrieves all students for a class section
// @Summary Get students by Section
// @Tags people
// @Produce json
// @Param section_id path int true "Section ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/people/sections/{section_id}/students [get]
func (c *PeopleController) GetStudentsBySection(ctx *fiber.Ctx) error {
	sectionID, err := utils.ParseInt64FromParam(ctx, "section_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Section ID"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	students, err := c.usecase.GetStudentsBySection(ctx.Context(), schoolID, sectionID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Section students retrieved successfully", Data: students})
}

// GetStudentByID retrieves student details by ID
// @Summary Get student by ID
// @Tags people
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/people/students/{id} [get]
func (c *PeopleController) GetStudentByID(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Student ID"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	student, err := c.usecase.GetStudentByID(ctx.Context(), schoolID, id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Student retrieved successfully", Data: student})
}

// GetStudentByUserID retrieves student details by user ID
// @Summary Get student by User ID
// @Tags people
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/people/students/users/{user_id} [get]
func (c *PeopleController) GetStudentByUserID(ctx *fiber.Ctx) error {
	userID, err := utils.ParseInt64FromParam(ctx, "user_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid User ID"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	student, err := c.usecase.GetStudentByUserID(ctx.Context(), schoolID, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Student retrieved successfully", Data: student})
}

// DeleteStudent handles student deletion
func (c *PeopleController) DeleteStudent(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Student ID"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.usecase.DeleteStudent(ctx.Context(), schoolID, id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Student deleted and user deactivated successfully"})
}

// EnrollStudent handles student enrollment in a section
// @Summary Enroll student in section
// @Tags people
// @Accept json
// @Produce json
// @Param enrollment body entity.StudentSection true "Enrollment details"
// @Success 201 {object} response.HTTPSuccessResponse
// @Router /api/v1/people/enroll [post]
func (c *PeopleController) EnrollStudent(ctx *fiber.Ctx) error {
	var enrollment entity.StudentSection
	if err := ctx.BodyParser(&enrollment); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)
	if err := c.usecase.EnrollStudent(ctx.Context(), schoolID, &enrollment); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.HTTPSuccessResponse{Status: fiber.StatusCreated, Message: "Student enrolled successfully"})
}

// GetStudentSections retrieves sections a student is enrolled in
// @Summary Get student sections
// @Tags people
// @Produce json
// @Param student_id path int true "Student ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/people/students/{student_id}/sections [get]
func (c *PeopleController) GetStudentSections(ctx *fiber.Ctx) error {
	studentID, err := utils.ParseInt64FromParam(ctx, "student_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Student ID"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	sections, err := c.usecase.GetStudentSections(ctx.Context(), schoolID, studentID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Student sections retrieved successfully", Data: sections})
}

// CreateParent handles parent creation
// @Summary Create parent
// @Tags people
// @Accept json
// @Produce json
// @Param parent body entity.Parent true "Parent details"
// @Success 201 {object} response.HTTPSuccessResponse
// @Router /api/v1/people/parents [post]
func (c *PeopleController) CreateParent(ctx *fiber.Ctx) error {
	var parent entity.Parent
	if err := ctx.BodyParser(&parent); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)
	resp, err := c.usecase.CreateParent(ctx.Context(), schoolID, &parent)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusCreated,
		Message: "Parent created successfully",
		Data:    resp,
	})
}

// GetParents retrieves all parents for a school
// @Summary Get parents
// @Tags people
// @Produce json
// @Param school_id query int true "School ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/people/parents [get]
func (c *PeopleController) GetParents(ctx *fiber.Ctx) error {
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	parents, err := c.usecase.GetParentsBySchool(ctx.Context(), schoolID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Parents retrieved successfully", Data: parents})
}

// SearchParents searches parents by name, phone or email
// @Summary Search parents
// @Tags people
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/people/parents/search [get]
func (c *PeopleController) SearchParents(ctx *fiber.Ctx) error {
	query := ctx.Query("q")
	if query == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Query parameter 'q' is required"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	parents, err := c.usecase.SearchParents(ctx.Context(), schoolID, query)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Parents searched successfully", Data: parents})
}

// UpdateParent handles parent updates
// @Summary Update parent
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "Parent ID"
// @Param parent body entity.Parent true "Parent details"
// @Success 200 {object} response.HTTPSuccessResponse
// @Failure 400 {object} response.HTTPErrorResponse
// @Router /api/v1/people/parents/{id} [put]
func (c *PeopleController) UpdateParent(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Parent ID"})
	}

	var parent entity.Parent
	if err := ctx.BodyParser(&parent); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}
	parent.ID = id

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.usecase.UpdateParent(ctx.Context(), schoolID, &parent); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Parent updated successfully"})
}

// DeleteParent handles parent deletion
// @Summary Delete parent
// @Tags people
// @Param id path int true "Parent ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Failure 400 {object} response.HTTPErrorResponse
// @Router /api/v1/people/parents/{id} [delete]
func (c *PeopleController) DeleteParent(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Parent ID"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.usecase.DeleteParent(ctx.Context(), schoolID, id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Parent deleted and user deactivated successfully"})
}

// LinkParentToStudent links a parent to a student
// @Summary Link parent to student
// @Tags people
// @Accept json
// @Produce json
// @Param link body entity.StudentParent true "Link details"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/people/parents/link [post]
func (c *PeopleController) LinkParentToStudent(ctx *fiber.Ctx) error {
	var link entity.StudentParent
	if err := ctx.BodyParser(&link); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)
	if err := c.usecase.LinkParentToStudent(ctx.Context(), schoolID, &link); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Parent linked to student successfully"})
}

// CreateStaff handles staff onboarding
// @Summary Create school staff
// @Tags people
// @Accept json
// @Produce json
// @Param staff body entity.Staff true "Staff details"
// @Success 201 {object} response.HTTPSuccessResponse
// @Failure 400 {object} response.HTTPErrorResponse
// @Router /api/v1/people/staff [post]
func (c *PeopleController) CreateStaff(ctx *fiber.Ctx) error {
	var staff entity.Staff
	if err := ctx.BodyParser(&staff); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)
	resp, err := c.usecase.CreateStaff(ctx.Context(), schoolID, &staff)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.HTTPSuccessResponse{
		Status:  fiber.StatusCreated,
		Message: "Staff created successfully",
		Data:    resp,
	})
}

// UpdateStaff handles staff updates
// @Summary Update school staff
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "Staff ID"
// @Param staff body entity.Staff true "Staff details"
// @Success 200 {object} response.HTTPSuccessResponse
// @Failure 400 {object} response.HTTPErrorResponse
// @Router /api/v1/people/staff/{id} [put]
func (c *PeopleController) UpdateStaff(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Staff ID"})
	}

	var staff entity.Staff
	if err := ctx.BodyParser(&staff); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid body"})
	}
	staff.ID = id

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.usecase.UpdateStaff(ctx.Context(), schoolID, &staff); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Staff updated successfully"})
}

// GetStaff retrieves all staff for a school
// @Summary Get staff
// @Tags people
// @Produce json
// @Param school_id query int true "School ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/people/staff [get]
func (c *PeopleController) GetStaff(ctx *fiber.Ctx) error {
	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	staff, err := c.usecase.GetStaffBySchool(ctx.Context(), schoolID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Staff retrieved successfully", Data: staff})
}

// DeleteStaff handles staff deletion
// @Summary Delete school staff
// @Tags people
// @Param id path int true "Staff ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Failure 400 {object} response.HTTPErrorResponse
// @Router /api/v1/people/staff/{id} [delete]
func (c *PeopleController) DeleteStaff(ctx *fiber.Ctx) error {
	id, err := utils.ParseInt64FromParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid Staff ID"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	if err := c.usecase.DeleteStaff(ctx.Context(), schoolID, id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Staff deleted and user deactivated successfully"})
}

// GetStaffByUserID retrieves staff details by user ID
// @Summary Get staff by User ID
// @Tags people
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} response.HTTPSuccessResponse
// @Router /api/v1/people/staff/users/{user_id} [get]
func (c *PeopleController) GetStaffByUserID(ctx *fiber.Ctx) error {
	userID, err := utils.ParseInt64FromParam(ctx, "user_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.HTTPErrorResponse{Status: fiber.StatusBadRequest, Message: "Invalid User ID"})
	}

	schoolID, _ := utils.GetSchoolIDFromToken(ctx)

	staff, err := c.usecase.GetStaffByUserID(ctx.Context(), schoolID, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.HTTPErrorResponse{Status: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.HTTPSuccessResponse{Status: fiber.StatusOK, Message: "Staff retrieved successfully", Data: staff})
}

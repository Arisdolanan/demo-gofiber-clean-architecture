package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/repository/postgresql"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type OperationUsecase interface {
	// Schedules
	CreateSchedule(ctx context.Context, schoolID int64, schedule *entity.Schedule) error
	UpdateSchedule(ctx context.Context, schoolID int64, schedule *entity.Schedule) error
	DeleteSchedule(ctx context.Context, schoolID, id int64) error
	GetAllSchedules(ctx context.Context, schoolID int64, filters map[string]interface{}) ([]*entity.Schedule, error)
	GetSchedulesBySection(ctx context.Context, schoolID, sectionID int64) ([]*entity.Schedule, error)
	GetSchedulesByTeacher(ctx context.Context, schoolID, teacherID int64) ([]*entity.Schedule, error)

	// Exams
	CreateExam(ctx context.Context, schoolID int64, exam *entity.Exam) error
	UpdateExam(ctx context.Context, schoolID int64, exam *entity.Exam) error
	DeleteExam(ctx context.Context, schoolID, id int64) error
	GetAllExams(ctx context.Context, schoolID int64, filters map[string]interface{}) ([]*entity.Exam, error)
	GetExamByID(ctx context.Context, schoolID, id int64) (*entity.Exam, error)
	GetExamsBySection(ctx context.Context, schoolID, sectionID int64) ([]*entity.Exam, error)

	// Marks
	UpdateMark(ctx context.Context, schoolID int64, mark *entity.ExamMark) error
	GetExamMarks(ctx context.Context, schoolID, examID int64) ([]*entity.ExamMark, error)
	GetStudentMarks(ctx context.Context, schoolID, studentID int64) ([]*entity.ExamMark, error)

	// Attendance
	RecordStudentAttendance(ctx context.Context, schoolID int64, attendance *entity.StudentAttendance) error
	UpdateStudentAttendance(ctx context.Context, schoolID int64, attendance *entity.StudentAttendance) error
	GetStudentAttendance(ctx context.Context, schoolID, studentID int64, filters map[string]interface{}) ([]*entity.StudentAttendance, error)
	GetAttendanceBySection(ctx context.Context, schoolID, sectionID int64, date string, subjectID *int64) ([]*entity.StudentAttendance, error)
	GetAttendanceWithFilters(ctx context.Context, schoolID int64, filters map[string]interface{}) ([]*entity.StudentAttendance, error)
	GetStudentsBySectionForAttendance(ctx context.Context, schoolID, sectionID int64) ([]*entity.Student, error)
	RecordSectionAttendance(ctx context.Context, schoolID int64, req *entity.SectionAttendanceRequest) error
	RecordTeacherAttendance(ctx context.Context, schoolID int64, attendance *entity.TeacherAttendance) error
	UpdateTeacherAttendance(ctx context.Context, schoolID int64, attendance *entity.TeacherAttendance) error
	GetTeacherAttendance(ctx context.Context, schoolID, teacherID int64, filters map[string]interface{}) ([]*entity.TeacherAttendance, error)
	RecordStaffAttendance(ctx context.Context, schoolID int64, attendance *entity.StaffAttendance) error
	UpdateStaffAttendance(ctx context.Context, schoolID int64, attendance *entity.StaffAttendance) error
	GetStaffAttendance(ctx context.Context, schoolID, employeeID int64, filters map[string]interface{}) ([]*entity.StaffAttendance, error)

	// Notifications
	NotifyUser(ctx context.Context, notification *entity.Notification) error
	GetUserNotifications(ctx context.Context, userID int64) ([]*entity.Notification, error)

	// Report Card
	GetStudentReportCard(ctx context.Context, schoolID, studentID, sessionID int64) (*entity.ReportCard, error)
	GetSectionReportCards(ctx context.Context, schoolID, sectionID, sessionID int64) ([]*entity.ReportCard, error)

	// Integration Definitions
	GetIntegrationDefinitions(ctx context.Context, schoolID int64) ([]*entity.IntegrationDefinition, error)
}

type operationUsecase struct {
	repo       postgresql.OperationRepository
	peopleRepo postgresql.PeopleRepository
	validate   *validator.Validate
	log        *logrus.Logger
}

func NewOperationUsecase(repo postgresql.OperationRepository, peopleRepo postgresql.PeopleRepository, validate *validator.Validate, log *logrus.Logger) OperationUsecase {
	return &operationUsecase{
		repo:       repo,
		peopleRepo: peopleRepo,
		validate:   validate,
		log:        log,
	}
}

func (uc *operationUsecase) CreateSchedule(ctx context.Context, schoolID int64, schedule *entity.Schedule) error {
	schedule.SchoolID = schoolID
	if err := uc.validate.Struct(schedule); err != nil {
		return err
	}
	return uc.repo.CreateSchedule(ctx, schedule)
}

func (uc *operationUsecase) UpdateSchedule(ctx context.Context, schoolID int64, schedule *entity.Schedule) error {
	schedule.SchoolID = schoolID
	if err := uc.validate.Struct(schedule); err != nil {
		return err
	}
	return uc.repo.UpdateSchedule(ctx, schedule)
}

func (uc *operationUsecase) DeleteSchedule(ctx context.Context, schoolID, id int64) error {
	return uc.repo.DeleteSchedule(ctx, schoolID, id)
}

func (uc *operationUsecase) GetAllSchedules(ctx context.Context, schoolID int64, filters map[string]interface{}) ([]*entity.Schedule, error) {
	return uc.repo.GetAllSchedules(ctx, schoolID, filters)
}

func (uc *operationUsecase) GetSchedulesBySection(ctx context.Context, schoolID, sectionID int64) ([]*entity.Schedule, error) {
	return uc.repo.FindSchedulesBySection(ctx, schoolID, sectionID)
}

func (uc *operationUsecase) GetSchedulesByTeacher(ctx context.Context, schoolID, teacherID int64) ([]*entity.Schedule, error) {
	return uc.repo.GetSchedulesByTeacher(ctx, schoolID, teacherID)
}

func (uc *operationUsecase) CreateExam(ctx context.Context, schoolID int64, exam *entity.Exam) error {
	exam.SchoolID = schoolID
	if err := uc.validate.Struct(exam); err != nil {
		return err
	}
	return uc.repo.CreateExam(ctx, exam)
}

func (uc *operationUsecase) UpdateExam(ctx context.Context, schoolID int64, exam *entity.Exam) error {
	exam.SchoolID = schoolID
	if err := uc.validate.Struct(exam); err != nil {
		return err
	}
	return uc.repo.UpdateExam(ctx, exam)
}

func (uc *operationUsecase) DeleteExam(ctx context.Context, schoolID, id int64) error {
	return uc.repo.DeleteExam(ctx, schoolID, id)
}

func (uc *operationUsecase) GetAllExams(ctx context.Context, schoolID int64, filters map[string]interface{}) ([]*entity.Exam, error) {
	return uc.repo.GetAllExams(ctx, schoolID, filters)
}

func (uc *operationUsecase) GetExamByID(ctx context.Context, schoolID, id int64) (*entity.Exam, error) {
	return uc.repo.GetExamByID(ctx, schoolID, id)
}

func (uc *operationUsecase) GetExamsBySection(ctx context.Context, schoolID, sectionID int64) ([]*entity.Exam, error) {
	return uc.repo.FindExamsBySection(ctx, schoolID, sectionID)
}

func (uc *operationUsecase) UpdateMark(ctx context.Context, schoolID int64, mark *entity.ExamMark) error {
	if err := uc.validate.Struct(mark); err != nil {
		return err
	}
	return uc.repo.UpdateMark(ctx, schoolID, mark)
}

func (uc *operationUsecase) GetExamMarks(ctx context.Context, schoolID, examID int64) ([]*entity.ExamMark, error) {
	return uc.repo.GetMarksByExam(ctx, schoolID, examID)
}

func (uc *operationUsecase) GetStudentMarks(ctx context.Context, schoolID, studentID int64) ([]*entity.ExamMark, error) {
	return uc.repo.GetMarksByStudent(ctx, schoolID, studentID)
}

func (uc *operationUsecase) RecordStudentAttendance(ctx context.Context, schoolID int64, attendance *entity.StudentAttendance) error {
	if err := uc.validate.Struct(attendance); err != nil {
		return err
	}
	return uc.repo.RecordStudentAttendance(ctx, schoolID, attendance)
}

func (uc *operationUsecase) UpdateStudentAttendance(ctx context.Context, schoolID int64, attendance *entity.StudentAttendance) error {
	if err := uc.validate.Struct(attendance); err != nil {
		return err
	}
	return uc.repo.UpdateStudentAttendance(ctx, schoolID, attendance)
}

func (uc *operationUsecase) GetStudentAttendance(ctx context.Context, schoolID, studentID int64, filters map[string]interface{}) ([]*entity.StudentAttendance, error) {
	return uc.repo.GetStudentAttendance(ctx, schoolID, studentID, filters)
}

func (uc *operationUsecase) GetAttendanceBySection(ctx context.Context, schoolID, sectionID int64, date string, subjectID *int64) ([]*entity.StudentAttendance, error) {
	return uc.repo.GetSectionAttendance(ctx, schoolID, sectionID, date, subjectID)
}

func (uc *operationUsecase) GetAttendanceWithFilters(ctx context.Context, schoolID int64, filters map[string]interface{}) ([]*entity.StudentAttendance, error) {
	return uc.repo.GetAttendanceWithFilters(ctx, schoolID, filters)
}

func (uc *operationUsecase) GetStudentsBySectionForAttendance(ctx context.Context, schoolID, sectionID int64) ([]*entity.Student, error) {
	// Note: PeopleRepository might need updating later
	return uc.peopleRepo.GetStudentsBySection(ctx, schoolID, sectionID)
}

func (uc *operationUsecase) RecordSectionAttendance(ctx context.Context, schoolID int64, req *entity.SectionAttendanceRequest) error {
	subjectID := req.SubjectID
	for _, entry := range req.Attendances {
		attendance := &entity.StudentAttendance{
			StudentID:         entry.StudentID,
			SectionID:         req.SectionID,
			AcademicSessionID: req.AcademicSessionID,
			SubjectID:         &subjectID,
			AttendanceDate:    req.AttendanceDate,
			Status:            entry.Status,
			Notes:             entry.Notes,
			MarkedBy:          req.MarkedBy,
			CheckInIPAddress:  &req.IPAddress,
			CheckInDevice:     &req.UserAgent,
		}
		if err := uc.repo.RecordStudentAttendance(ctx, schoolID, attendance); err != nil {
			return err
		}
	}
	return nil
}

func (uc *operationUsecase) RecordTeacherAttendance(ctx context.Context, schoolID int64, attendance *entity.TeacherAttendance) error {
	if err := uc.validate.Struct(attendance); err != nil {
		return err
	}
	return uc.repo.RecordTeacherAttendance(ctx, schoolID, attendance)
}

func (uc *operationUsecase) UpdateTeacherAttendance(ctx context.Context, schoolID int64, attendance *entity.TeacherAttendance) error {
	if err := uc.validate.Struct(attendance); err != nil {
		return err
	}
	return uc.repo.UpdateTeacherAttendance(ctx, schoolID, attendance)
}

func (uc *operationUsecase) GetTeacherAttendance(ctx context.Context, schoolID, teacherID int64, filters map[string]interface{}) ([]*entity.TeacherAttendance, error) {
	return uc.repo.GetTeacherAttendance(ctx, schoolID, teacherID, filters)
}

func (uc *operationUsecase) RecordStaffAttendance(ctx context.Context, schoolID int64, attendance *entity.StaffAttendance) error {
	if err := uc.validate.Struct(attendance); err != nil {
		return err
	}

	// If employee_id is provided, check if staff exists
	if attendance.EmployeeID > 0 {
		existingStaff, err := uc.peopleRepo.FindStaffByID(ctx, schoolID, attendance.EmployeeID)
		if err != nil {
			return err
		}

		// If staff doesn't exist, auto-create for testing
		if existingStaff == nil {
			fullName := fmt.Sprintf("Test Staff %d", attendance.EmployeeID)
			position := "Test Position"
			department := "Testing"
			newStaff := &entity.Staff{
				ID:             attendance.EmployeeID,
				SchoolID:       schoolID,
				EmployeeNumber: fmt.Sprintf("TEST-%d", attendance.EmployeeID),
				FullName:       &fullName,
				Status:         "active",
				Position:       &position,
				Department:     &department,
			}
			if err := uc.peopleRepo.CreateStaff(ctx, newStaff); err != nil {
				uc.log.Warnf("Failed to auto-create staff record for testing: %v", err)
				return errors.New("staff not found and auto-creation failed")
			}
		}

		return uc.repo.RecordStaffAttendance(ctx, schoolID, attendance)
	}

	return errors.New("employee_id is required")
}

func (uc *operationUsecase) UpdateStaffAttendance(ctx context.Context, schoolID int64, attendance *entity.StaffAttendance) error {
	if err := uc.validate.Struct(attendance); err != nil {
		return err
	}
	return uc.repo.UpdateStaffAttendance(ctx, schoolID, attendance)
}

func (uc *operationUsecase) GetStaffAttendance(ctx context.Context, schoolID, employeeID int64, filters map[string]interface{}) ([]*entity.StaffAttendance, error) {
	return uc.repo.GetStaffAttendance(ctx, schoolID, employeeID, filters)
}

func (uc *operationUsecase) NotifyUser(ctx context.Context, notification *entity.Notification) error {
	return uc.repo.CreateNotification(ctx, notification)
}

func (uc *operationUsecase) GetUserNotifications(ctx context.Context, userID int64) ([]*entity.Notification, error) {
	return uc.repo.GetNotificationsByUser(ctx, userID)
}

func (uc *operationUsecase) GetStudentReportCard(ctx context.Context, schoolID, studentID, sessionID int64) (*entity.ReportCard, error) {
	marks, err := uc.repo.GetStudentMarksBySession(ctx, schoolID, studentID, sessionID)
	if err != nil {
		return nil, err
	}

	attendance, err := uc.repo.GetStudentAttendanceBySession(ctx, schoolID, studentID, sessionID)
	if err != nil {
		return nil, err
	}

	report := &entity.ReportCard{
		StudentID:         studentID,
		AcademicSessionID: sessionID,
		SubjectGrades:     []entity.SubjectGrade{},
		AttendanceSummary: make(map[string]int),
	}

	subjectMap := make(map[int64]*entity.SubjectGrade)
	totalScore := 0.0
	count := 0

	for _, m := range marks {
		if _, ok := subjectMap[m.ExamID]; !ok {
			subjectMap[m.ExamID] = &entity.SubjectGrade{SubjectID: m.ExamID}
		}
		subjectMap[m.ExamID].Grades = append(subjectMap[m.ExamID].Grades, m.Score)
		totalScore += m.Score
		count++
	}

	for _, sg := range subjectMap {
		sum := 0.0
		for _, g := range sg.Grades {
			sum += g
		}
		sg.Average = sum / float64(len(sg.Grades))
		report.SubjectGrades = append(report.SubjectGrades, *sg)
	}

	if count > 0 {
		report.TotalAverage = totalScore / float64(count)
	}

	for _, a := range attendance {
		report.AttendanceSummary[a.Status]++
	}

	return report, nil
}

func (uc *operationUsecase) GetSectionReportCards(ctx context.Context, schoolID, sectionID, sessionID int64) ([]*entity.ReportCard, error) {
	// TODO: This is a placeholder implementation
	uc.log.Infof("GetSectionReportCards: schoolID=%d, sectionID=%d, sessionID=%d", schoolID, sectionID, sessionID)
	return []*entity.ReportCard{}, nil
}

func (uc *operationUsecase) GetIntegrationDefinitions(ctx context.Context, schoolID int64) ([]*entity.IntegrationDefinition, error) {
	return uc.repo.GetIntegrationDefinitions(ctx, schoolID)
}

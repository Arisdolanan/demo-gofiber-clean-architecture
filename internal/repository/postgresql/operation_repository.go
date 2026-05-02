package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/jmoiron/sqlx"
)

type OperationRepository interface {
	// Schedules
	CreateSchedule(ctx context.Context, schedule *entity.Schedule) error
	UpdateSchedule(ctx context.Context, schedule *entity.Schedule) error
	DeleteSchedule(ctx context.Context, schoolID, id int64) error
	GetAllSchedules(ctx context.Context, schoolID int64, filters map[string]interface{}) ([]*entity.Schedule, error)
	FindSchedulesBySection(ctx context.Context, schoolID, sectionID int64) ([]*entity.Schedule, error)
	GetSchedulesByTeacher(ctx context.Context, schoolID, teacherID int64) ([]*entity.Schedule, error)

	// Exams
	CreateExam(ctx context.Context, exam *entity.Exam) error
	UpdateExam(ctx context.Context, exam *entity.Exam) error
	DeleteExam(ctx context.Context, schoolID, id int64) error
	GetAllExams(ctx context.Context, schoolID int64, filters map[string]interface{}) ([]*entity.Exam, error)
	GetExamByID(ctx context.Context, schoolID, id int64) (*entity.Exam, error)
	FindExamsBySection(ctx context.Context, schoolID, sectionID int64) ([]*entity.Exam, error)

	// Marks
	UpdateMark(ctx context.Context, schoolID int64, mark *entity.ExamMark) error
	GetMarksByExam(ctx context.Context, schoolID, examID int64) ([]*entity.ExamMark, error)
	GetMarksByStudent(ctx context.Context, schoolID, studentID int64) ([]*entity.ExamMark, error)

	// Attendance
	RecordStudentAttendance(ctx context.Context, schoolID int64, attendance *entity.StudentAttendance) error
	UpdateStudentAttendance(ctx context.Context, schoolID int64, attendance *entity.StudentAttendance) error
	GetStudentAttendance(ctx context.Context, schoolID, studentID int64, filters map[string]interface{}) ([]*entity.StudentAttendance, error)
	GetSectionAttendance(ctx context.Context, schoolID, sectionID int64, date string, subjectID *int64) ([]*entity.StudentAttendance, error)
	GetAttendanceWithFilters(ctx context.Context, schoolID int64, filters map[string]interface{}) ([]*entity.StudentAttendance, error)
	RecordTeacherAttendance(ctx context.Context, schoolID int64, attendance *entity.TeacherAttendance) error
	UpdateTeacherAttendance(ctx context.Context, schoolID int64, attendance *entity.TeacherAttendance) error
	GetTeacherAttendance(ctx context.Context, schoolID, teacherID int64, filters map[string]interface{}) ([]*entity.TeacherAttendance, error)
	RecordStaffAttendance(ctx context.Context, schoolID int64, attendance *entity.StaffAttendance) error
	UpdateStaffAttendance(ctx context.Context, schoolID int64, attendance *entity.StaffAttendance) error
	GetStaffAttendance(ctx context.Context, schoolID, employeeID int64, filters map[string]interface{}) ([]*entity.StaffAttendance, error)

	// Notifications
	CreateNotification(ctx context.Context, notification *entity.Notification) error
	GetNotificationsByUser(ctx context.Context, userID int64) ([]*entity.Notification, error)

	// Report Card / Academic History
	GetStudentMarksBySession(ctx context.Context, schoolID, studentID, sessionID int64) ([]*entity.ExamMark, error)
	GetStudentAttendanceBySession(ctx context.Context, schoolID, studentID, sessionID int64) ([]*entity.StudentAttendance, error)

	// Integration Definitions
	GetIntegrationDefinitions(ctx context.Context, schoolID int64) ([]*entity.IntegrationDefinition, error)
}

type operationRepository struct {
	scheduleRepo          *BaseRepository[entity.Schedule]
	examRepo              *BaseRepository[entity.Exam]
	examMarkRepo          *BaseRepository[entity.ExamMark]
	studentAttendanceRepo *BaseRepository[entity.StudentAttendance]
	teacherAttendanceRepo *BaseRepository[entity.TeacherAttendance]
	staffAttendanceRepo   *BaseRepository[entity.StaffAttendance]
	notificationRepo      *BaseRepository[entity.Notification]
	integrationDefRepo    *BaseRepository[entity.IntegrationDefinition]
	db                    *sqlx.DB
}

func NewOperationRepository(db *sqlx.DB) OperationRepository {
	return &operationRepository{
		scheduleRepo:          NewBaseRepository[entity.Schedule](db, "schedules"),
		examRepo:              NewBaseRepository[entity.Exam](db, "exams"),
		examMarkRepo:          NewBaseRepository[entity.ExamMark](db, "exam_marks"),
		studentAttendanceRepo: NewBaseRepository[entity.StudentAttendance](db, "student_attendance"),
		teacherAttendanceRepo: NewBaseRepository[entity.TeacherAttendance](db, "teacher_attendance"),
		staffAttendanceRepo:   NewBaseRepository[entity.StaffAttendance](db, "staff_attendance"),
		notificationRepo:      NewBaseRepository[entity.Notification](db, "notifications"),
		integrationDefRepo:    NewBaseRepository[entity.IntegrationDefinition](db, "integration_definitions"),
		db:                    db,
	}
}

func (r *operationRepository) CreateSchedule(ctx context.Context, schedule *entity.Schedule) error {
	return r.scheduleRepo.Create(ctx, schedule)
}

func (r *operationRepository) UpdateSchedule(ctx context.Context, schedule *entity.Schedule) error {
	return r.scheduleRepo.Update(ctx, schedule, "id = $1 AND deleted_at IS NULL", schedule.ID)
}

func (r *operationRepository) DeleteSchedule(ctx context.Context, schoolID, id int64) error {
	query := `id = $1 AND section_id IN (SELECT id FROM sections WHERE class_id IN (SELECT id FROM classes WHERE school_id = $2)) AND deleted_at IS NULL`
	return r.scheduleRepo.SoftDelete(ctx, query, id, schoolID)
}

func (r *operationRepository) GetAllSchedules(ctx context.Context, schoolID int64, filters map[string]interface{}) ([]*entity.Schedule, error) {
	query := `section_id IN (SELECT id FROM sections WHERE class_id IN (SELECT id FROM classes WHERE school_id = $1)) AND deleted_at IS NULL`
	return r.scheduleRepo.FindAll(ctx, query, schoolID)
}

func (r *operationRepository) FindSchedulesBySection(ctx context.Context, schoolID, sectionID int64) ([]*entity.Schedule, error) {
	query := `section_id = $1 AND section_id IN (SELECT id FROM sections WHERE class_id IN (SELECT id FROM classes WHERE school_id = $2)) AND deleted_at IS NULL`
	return r.scheduleRepo.FindAll(ctx, query, sectionID, schoolID)
}

func (r *operationRepository) GetSchedulesByTeacher(ctx context.Context, schoolID, teacherID int64) ([]*entity.Schedule, error) {
	query := `teacher_id = $1 AND teacher_id IN (SELECT id FROM teachers WHERE school_id = $2) AND deleted_at IS NULL`
	return r.scheduleRepo.FindAll(ctx, query, teacherID, schoolID)
}

func (r *operationRepository) CreateExam(ctx context.Context, exam *entity.Exam) error {
	return r.examRepo.Create(ctx, exam)
}

func (r *operationRepository) UpdateExam(ctx context.Context, exam *entity.Exam) error {
	return r.examRepo.Update(ctx, exam, "id = $1 AND deleted_at IS NULL", exam.ID)
}

func (r *operationRepository) DeleteExam(ctx context.Context, schoolID, id int64) error {
	query := `id = $1 AND section_id IN (SELECT id FROM sections WHERE class_id IN (SELECT id FROM classes WHERE school_id = $2)) AND deleted_at IS NULL`
	return r.examRepo.SoftDelete(ctx, query, id, schoolID)
}

func (r *operationRepository) GetAllExams(ctx context.Context, schoolID int64, filters map[string]interface{}) ([]*entity.Exam, error) {
	query := `section_id IN (SELECT id FROM sections WHERE class_id IN (SELECT id FROM classes WHERE school_id = $1)) AND deleted_at IS NULL`
	return r.examRepo.FindAll(ctx, query, schoolID)
}

func (r *operationRepository) GetExamByID(ctx context.Context, schoolID, id int64) (*entity.Exam, error) {
	query := `id = $1 AND section_id IN (SELECT id FROM sections WHERE class_id IN (SELECT id FROM classes WHERE school_id = $2)) AND deleted_at IS NULL`
	return r.examRepo.FindOne(ctx, query, id, schoolID)
}

func (r *operationRepository) FindExamsBySection(ctx context.Context, schoolID, sectionID int64) ([]*entity.Exam, error) {
	query := `section_id = $1 AND section_id IN (SELECT id FROM sections WHERE class_id IN (SELECT id FROM classes WHERE school_id = $2)) AND deleted_at IS NULL`
	return r.examRepo.FindAll(ctx, query, sectionID, schoolID)
}

func (r *operationRepository) UpdateMark(ctx context.Context, schoolID int64, mark *entity.ExamMark) error {
	// Verify exam belongs to school
	query := `id = $1 AND exam_id IN (SELECT id FROM exams WHERE section_id IN (SELECT id FROM sections WHERE class_id IN (SELECT id FROM classes WHERE school_id = $2))) AND deleted_at IS NULL`
	if mark.ID != 0 {
		return r.examMarkRepo.Update(ctx, mark, query, mark.ID, schoolID)
	}
	// For creation, we should also verify student belongs to school
	return r.examMarkRepo.Create(ctx, mark)
}

func (r *operationRepository) GetMarksByExam(ctx context.Context, schoolID, examID int64) ([]*entity.ExamMark, error) {
	query := `exam_id = $1 AND exam_id IN (SELECT id FROM exams WHERE section_id IN (SELECT id FROM sections WHERE class_id IN (SELECT id FROM classes WHERE school_id = $2))) AND deleted_at IS NULL`
	return r.examMarkRepo.FindAll(ctx, query, examID, schoolID)
}

func (r *operationRepository) GetMarksByStudent(ctx context.Context, schoolID, studentID int64) ([]*entity.ExamMark, error) {
	query := `student_id = $1 AND student_id IN (SELECT id FROM students WHERE school_id = $2) AND deleted_at IS NULL`
	return r.examMarkRepo.FindAll(ctx, query, studentID, schoolID)
}

func (r *operationRepository) RecordStudentAttendance(ctx context.Context, schoolID int64, attendance *entity.StudentAttendance) error {
	// Check if record already exists for this student on this date in this section and subject
	query := "student_id = $1 AND section_id = $2 AND attendance_date::DATE = $3::DATE AND deleted_at IS NULL"
	args := []interface{}{attendance.StudentID, attendance.SectionID, attendance.AttendanceDate}

	if attendance.SubjectID != nil {
		query += " AND subject_id = $4"
		args = append(args, *attendance.SubjectID)
	} else {
		query += " AND subject_id IS NULL"
	}

	existing, _ := r.studentAttendanceRepo.FindOne(ctx, query, args...)

	if existing != nil {
		attendance.ID = existing.ID
		return r.studentAttendanceRepo.Update(ctx, attendance, "id = $1", existing.ID)
	}
	return r.studentAttendanceRepo.Create(ctx, attendance)
}

func (r *operationRepository) UpdateStudentAttendance(ctx context.Context, schoolID int64, attendance *entity.StudentAttendance) error {
	query := `id = $1 AND student_id IN (SELECT id FROM students WHERE school_id = $2) AND deleted_at IS NULL`
	return r.studentAttendanceRepo.Update(ctx, attendance, query, attendance.ID, schoolID)
}

func (r *operationRepository) GetStudentAttendance(ctx context.Context, schoolID, studentID int64, filters map[string]interface{}) ([]*entity.StudentAttendance, error) {
	whereClause := "student_id IN (SELECT id FROM students WHERE school_id = $1) AND deleted_at IS NULL"
	args := []interface{}{schoolID}
	argCount := 2

	if studentID != 0 {
		whereClause += fmt.Sprintf(" AND student_id = $%d", argCount)
		args = append(args, studentID)
		argCount++
	}

	if endDate, ok := filters["end_date"]; ok {
		whereClause += fmt.Sprintf(" AND attendance_date <= $%d", argCount)
		args = append(args, endDate)
		argCount++
	}

	whereClause += " ORDER BY attendance_date DESC, id DESC"

	return r.studentAttendanceRepo.FindAll(ctx, whereClause, args...)
}

func (r *operationRepository) GetSectionAttendance(ctx context.Context, schoolID, sectionID int64, date string, subjectID *int64) ([]*entity.StudentAttendance, error) {
	// Use a temporary struct to handle JOIN results since StudentAttendance
	// display fields are marked with db:"-"
	type attendanceResult struct {
		entity.StudentAttendance
		StudentName   string         `db:"student_name"`
		StudentNumber string         `db:"student_number"`
		SubjectName   sql.NullString `db:"subject_name"`
		ClassName     sql.NullString `db:"class_name"`
		SectionName   sql.NullString `db:"section_name"`
	}

	var results []attendanceResult
	query := `
		SELECT 
			sa.*, 
			s.full_name as student_name, 
			s.student_number,
			sub.name as subject_name,
			c.name as class_name,
			sec.code as section_name
		FROM student_attendance sa
		JOIN students s ON sa.student_id = s.id
		LEFT JOIN subjects sub ON sa.subject_id = sub.id
		LEFT JOIN sections sec ON sa.section_id = sec.id
		LEFT JOIN classes c ON sec.class_id = c.id
		WHERE sa.section_id = $1 AND sa.attendance_date::DATE = $2::DATE AND sa.deleted_at IS NULL AND s.school_id = $3
		ORDER BY sa.attendance_date DESC, s.full_name ASC
	`
	args := []interface{}{sectionID, date, schoolID}
	if subjectID != nil {
		query = strings.Replace(query, "WHERE", "WHERE sa.subject_id = $4 AND", 1)
		args = append(args, *subjectID)
	}

	err := r.db.SelectContext(ctx, &results, query, args...)
	if err != nil {
		return nil, err
	}

	// Map back to entity.StudentAttendance
	finalResults := make([]*entity.StudentAttendance, len(results))
	for i := range results {
		res := &results[i]
		res.StudentAttendance.StudentName = &res.StudentName
		res.StudentAttendance.StudentNumber = &res.StudentNumber

		if res.SubjectName.Valid {
			subjName := res.SubjectName.String
			res.StudentAttendance.SubjectName = &subjName
		}
		if res.ClassName.Valid {
			clsName := res.ClassName.String
			res.StudentAttendance.ClassName = &clsName
		}
		if res.SectionName.Valid {
			secName := res.SectionName.String
			res.StudentAttendance.SectionName = &secName
		}

		finalResults[i] = &res.StudentAttendance
	}

	return finalResults, nil
}

func (r *operationRepository) GetAttendanceWithFilters(ctx context.Context, schoolID int64, filters map[string]interface{}) ([]*entity.StudentAttendance, error) {
	type attendanceResult struct {
		entity.StudentAttendance
		StudentName   string         `db:"student_name"`
		StudentNumber string         `db:"student_number"`
		SubjectName   sql.NullString `db:"subject_name"`
		ClassName     sql.NullString `db:"class_name"`
		SectionName   sql.NullString `db:"section_name"`
	}

	var results []attendanceResult
	query := `
		SELECT 
			sa.*, 
			s.full_name as student_name, 
			s.student_number,
			sub.name as subject_name,
			c.name as class_name,
			sec.code as section_name
		FROM student_attendance sa
		JOIN students s ON sa.student_id = s.id
		LEFT JOIN subjects sub ON sa.subject_id = sub.id
		LEFT JOIN sections sec ON sa.section_id = sec.id
		LEFT JOIN classes c ON sec.class_id = c.id
		WHERE sa.deleted_at IS NULL AND s.school_id = $1
	`

	args := []interface{}{schoolID}
	argCount := 2

	if sectionID, ok := filters["section_id"]; ok {
		query += fmt.Sprintf(" AND sa.section_id = $%d", argCount)
		args = append(args, sectionID)
		argCount++
	}
	if sessionID, ok := filters["academic_session_id"]; ok {
		query += fmt.Sprintf(" AND sa.academic_session_id = $%d", argCount)
		args = append(args, sessionID)
		argCount++
	}
	if date, ok := filters["date"]; ok {
		query += fmt.Sprintf(" AND sa.attendance_date::DATE = $%d::DATE", argCount)
		args = append(args, date)
		argCount++
	}
	if subjectID, ok := filters["subject_id"]; ok {
		query += fmt.Sprintf(" AND sa.subject_id = $%d", argCount)
		args = append(args, subjectID)
		argCount++
	}

	query += " ORDER BY sa.attendance_date DESC, s.full_name ASC"

	err := r.db.SelectContext(ctx, &results, query, args...)
	if err != nil {
		return nil, err
	}

	finalResults := make([]*entity.StudentAttendance, len(results))
	for i := range results {
		res := &results[i]
		res.StudentAttendance.StudentName = &res.StudentName
		res.StudentAttendance.StudentNumber = &res.StudentNumber

		if res.SubjectName.Valid {
			subjName := res.SubjectName.String
			res.StudentAttendance.SubjectName = &subjName
		}
		if res.ClassName.Valid {
			clsName := res.ClassName.String
			res.StudentAttendance.ClassName = &clsName
		}
		if res.SectionName.Valid {
			secName := res.SectionName.String
			res.StudentAttendance.SectionName = &secName
		}

		finalResults[i] = &res.StudentAttendance
	}

	return finalResults, nil
}

func (r *operationRepository) RecordTeacherAttendance(ctx context.Context, schoolID int64, attendance *entity.TeacherAttendance) error {
	// Check if record already exists for this teacher on this date
	existing, _ := r.teacherAttendanceRepo.FindOne(ctx,
		"teacher_id = $1 AND attendance_date::DATE = $2::DATE AND deleted_at IS NULL",
		attendance.TeacherID, attendance.AttendanceDate)

	if existing != nil {
		attendance.ID = existing.ID
		return r.teacherAttendanceRepo.Update(ctx, attendance, "id = $1", existing.ID)
	}
	return r.teacherAttendanceRepo.Create(ctx, attendance)
}

func (r *operationRepository) UpdateTeacherAttendance(ctx context.Context, schoolID int64, attendance *entity.TeacherAttendance) error {
	query := `id = $1 AND teacher_id IN (SELECT id FROM teachers WHERE school_id = $2) AND deleted_at IS NULL`
	return r.teacherAttendanceRepo.Update(ctx, attendance, query, attendance.ID, schoolID)
}

func (r *operationRepository) GetTeacherAttendance(ctx context.Context, schoolID, teacherID int64, filters map[string]interface{}) ([]*entity.TeacherAttendance, error) {
	whereClause := "teacher_id IN (SELECT id FROM teachers WHERE school_id = $1) AND deleted_at IS NULL"
	args := []interface{}{schoolID}
	argCount := 2

	if teacherID != 0 {
		whereClause += fmt.Sprintf(" AND teacher_id = $%d", argCount)
		args = append(args, teacherID)
		argCount++
	}

	if startDate, ok := filters["start_date"]; ok {
		whereClause += fmt.Sprintf(" AND attendance_date >= $%d", argCount)
		args = append(args, startDate)
		argCount++
	}

	if endDate, ok := filters["end_date"]; ok {
		whereClause += fmt.Sprintf(" AND attendance_date <= $%d", argCount)
		args = append(args, endDate)
		argCount++
	}

	whereClause += " ORDER BY attendance_date DESC, id DESC"

	return r.teacherAttendanceRepo.FindAll(ctx, whereClause, args...)
}

func (r *operationRepository) RecordStaffAttendance(ctx context.Context, schoolID int64, attendance *entity.StaffAttendance) error {
	// Check if record already exists for this employee on this date
	existing, _ := r.staffAttendanceRepo.FindOne(ctx,
		"employee_id = $1 AND attendance_date::DATE = $2::DATE AND deleted_at IS NULL",
		attendance.EmployeeID, attendance.AttendanceDate)

	if existing != nil {
		attendance.ID = existing.ID
		return r.staffAttendanceRepo.Update(ctx, attendance, "id = $1", existing.ID)
	}
	return r.staffAttendanceRepo.Create(ctx, attendance)
}

func (r *operationRepository) UpdateStaffAttendance(ctx context.Context, schoolID int64, attendance *entity.StaffAttendance) error {
	// Note: employee_id here usually refers to staff table
	query := `id = $1 AND employee_id IN (SELECT id FROM staff WHERE school_id = $2) AND deleted_at IS NULL`
	return r.staffAttendanceRepo.Update(ctx, attendance, query, attendance.ID, schoolID)
}

func (r *operationRepository) GetStaffAttendance(ctx context.Context, schoolID, employeeID int64, filters map[string]interface{}) ([]*entity.StaffAttendance, error) {
	whereClause := "employee_id IN (SELECT id FROM staff WHERE school_id = $1) AND deleted_at IS NULL"
	args := []interface{}{schoolID}
	argCount := 2

	if employeeID != 0 {
		whereClause += fmt.Sprintf(" AND employee_id = $%d", argCount)
		args = append(args, employeeID)
		argCount++
	}

	if startDate, ok := filters["start_date"]; ok {
		whereClause += fmt.Sprintf(" AND attendance_date >= $%d", argCount)
		args = append(args, startDate)
		argCount++
	}

	if endDate, ok := filters["end_date"]; ok {
		whereClause += fmt.Sprintf(" AND attendance_date <= $%d", argCount)
		args = append(args, endDate)
		argCount++
	}

	whereClause += " ORDER BY attendance_date DESC, id DESC"

	return r.staffAttendanceRepo.FindAll(ctx, whereClause, args...)
}

func (r *operationRepository) CreateNotification(ctx context.Context, notification *entity.Notification) error {
	return r.notificationRepo.Create(ctx, notification)
}

func (r *operationRepository) GetNotificationsByUser(ctx context.Context, userID int64) ([]*entity.Notification, error) {
	return r.notificationRepo.FindAll(ctx, "user_id = $1 AND deleted_at IS NULL", userID)
}

func (r *operationRepository) GetStudentMarksBySession(ctx context.Context, schoolID, studentID, sessionID int64) ([]*entity.ExamMark, error) {
	var marks []*entity.ExamMark
	query := `
		SELECT em.* FROM exam_marks em
		JOIN exams e ON em.exam_id = e.id
		JOIN students s ON em.student_id = s.id
		WHERE em.student_id = $1 AND e.academic_session_id = $2 AND em.deleted_at IS NULL AND e.deleted_at IS NULL AND s.school_id = $3
	`
	err := r.db.SelectContext(ctx, &marks, query, studentID, sessionID, schoolID)
	return marks, err
}

func (r *operationRepository) GetStudentAttendanceBySession(ctx context.Context, schoolID, studentID, sessionID int64) ([]*entity.StudentAttendance, error) {
	query := `student_id = $1 AND academic_session_id = $2 AND student_id IN (SELECT id FROM students WHERE school_id = $3) AND deleted_at IS NULL`
	return r.studentAttendanceRepo.FindAll(ctx, query, studentID, sessionID, schoolID)
}

func (r *operationRepository) GetIntegrationDefinitions(ctx context.Context, schoolID int64) ([]*entity.IntegrationDefinition, error) {
	query := `school_id = $1 AND deleted_at IS NULL ORDER BY category, name`
	return r.integrationDefRepo.FindAll(ctx, query, schoolID)
}

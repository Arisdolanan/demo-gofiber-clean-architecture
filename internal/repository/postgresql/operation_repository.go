package postgresql

import (
	"context"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/jmoiron/sqlx"
)

type OperationRepository interface {
	// Schedules
	CreateSchedule(ctx context.Context, schedule *entity.Schedule) error
	UpdateSchedule(ctx context.Context, schedule *entity.Schedule) error
	DeleteSchedule(ctx context.Context, id int64) error
	GetAllSchedules(ctx context.Context, filters map[string]interface{}) ([]*entity.Schedule, error)
	FindSchedulesBySection(ctx context.Context, sectionID int64) ([]*entity.Schedule, error)
	GetSchedulesByTeacher(ctx context.Context, teacherID int64) ([]*entity.Schedule, error)

	// Exams
	CreateExam(ctx context.Context, exam *entity.Exam) error
	UpdateExam(ctx context.Context, exam *entity.Exam) error
	DeleteExam(ctx context.Context, id int64) error
	GetAllExams(ctx context.Context, filters map[string]interface{}) ([]*entity.Exam, error)
	GetExamByID(ctx context.Context, id int64) (*entity.Exam, error)
	FindExamsBySection(ctx context.Context, sectionID int64) ([]*entity.Exam, error)

	// Marks
	UpdateMark(ctx context.Context, mark *entity.ExamMark) error
	GetMarksByExam(ctx context.Context, examID int64) ([]*entity.ExamMark, error)
	GetMarksByStudent(ctx context.Context, studentID int64) ([]*entity.ExamMark, error)

	// Attendance
	RecordStudentAttendance(ctx context.Context, attendance *entity.StudentAttendance) error
	UpdateStudentAttendance(ctx context.Context, attendance *entity.StudentAttendance) error
	GetStudentAttendance(ctx context.Context, studentID int64, filters map[string]interface{}) ([]*entity.StudentAttendance, error)
	GetSectionAttendance(ctx context.Context, sectionID int64, date string, subjectID *int64) ([]*entity.StudentAttendance, error)
	GetAttendanceWithFilters(ctx context.Context, filters map[string]interface{}) ([]*entity.StudentAttendance, error)
	RecordTeacherAttendance(ctx context.Context, attendance *entity.TeacherAttendance) error
	UpdateTeacherAttendance(ctx context.Context, attendance *entity.TeacherAttendance) error
	GetTeacherAttendance(ctx context.Context, teacherID int64, filters map[string]interface{}) ([]*entity.TeacherAttendance, error)
	RecordStaffAttendance(ctx context.Context, attendance *entity.StaffAttendance) error
	UpdateStaffAttendance(ctx context.Context, attendance *entity.StaffAttendance) error
	GetStaffAttendance(ctx context.Context, employeeID int64, filters map[string]interface{}) ([]*entity.StaffAttendance, error)

	// Notifications
	CreateNotification(ctx context.Context, notification *entity.Notification) error
	GetNotificationsByUser(ctx context.Context, userID int64) ([]*entity.Notification, error)

	// Report Card / Academic History
	GetStudentMarksBySession(ctx context.Context, studentID, sessionID int64) ([]*entity.ExamMark, error)
	GetStudentAttendanceBySession(ctx context.Context, studentID, sessionID int64) ([]*entity.StudentAttendance, error)
}

type operationRepository struct {
	scheduleRepo          *BaseRepository[entity.Schedule]
	examRepo              *BaseRepository[entity.Exam]
	examMarkRepo          *BaseRepository[entity.ExamMark]
	studentAttendanceRepo *BaseRepository[entity.StudentAttendance]
	teacherAttendanceRepo *BaseRepository[entity.TeacherAttendance]
	staffAttendanceRepo   *BaseRepository[entity.StaffAttendance]
	notificationRepo      *BaseRepository[entity.Notification]
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
		db:                    db,
	}
}

func (r *operationRepository) CreateSchedule(ctx context.Context, schedule *entity.Schedule) error {
	return r.scheduleRepo.Create(ctx, schedule)
}

func (r *operationRepository) UpdateSchedule(ctx context.Context, schedule *entity.Schedule) error {
	return r.scheduleRepo.Update(ctx, schedule, "id = $1 AND deleted_at IS NULL", schedule.ID)
}

func (r *operationRepository) DeleteSchedule(ctx context.Context, id int64) error {
	return r.scheduleRepo.SoftDelete(ctx, "id = $1 AND deleted_at IS NULL", id)
}

func (r *operationRepository) GetAllSchedules(ctx context.Context, filters map[string]interface{}) ([]*entity.Schedule, error) {
	return r.scheduleRepo.FindAll(ctx, "deleted_at IS NULL")
}

func (r *operationRepository) FindSchedulesBySection(ctx context.Context, sectionID int64) ([]*entity.Schedule, error) {
	return r.scheduleRepo.FindAll(ctx, "section_id = $1 AND deleted_at IS NULL", sectionID)
}

func (r *operationRepository) GetSchedulesByTeacher(ctx context.Context, teacherID int64) ([]*entity.Schedule, error) {
	return r.scheduleRepo.FindAll(ctx, "teacher_id = $1 AND deleted_at IS NULL", teacherID)
}

func (r *operationRepository) CreateExam(ctx context.Context, exam *entity.Exam) error {
	return r.examRepo.Create(ctx, exam)
}

func (r *operationRepository) UpdateExam(ctx context.Context, exam *entity.Exam) error {
	return r.examRepo.Update(ctx, exam, "id = $1 AND deleted_at IS NULL", exam.ID)
}

func (r *operationRepository) DeleteExam(ctx context.Context, id int64) error {
	return r.examRepo.SoftDelete(ctx, "id = $1 AND deleted_at IS NULL", id)
}

func (r *operationRepository) GetAllExams(ctx context.Context, filters map[string]interface{}) ([]*entity.Exam, error) {
	return r.examRepo.FindAll(ctx, "deleted_at IS NULL")
}

func (r *operationRepository) GetExamByID(ctx context.Context, id int64) (*entity.Exam, error) {
	return r.examRepo.FindByID(ctx, id)
}

func (r *operationRepository) FindExamsBySection(ctx context.Context, sectionID int64) ([]*entity.Exam, error) {
	return r.examRepo.FindAll(ctx, "section_id = $1 AND deleted_at IS NULL", sectionID)
}

func (r *operationRepository) UpdateMark(ctx context.Context, mark *entity.ExamMark) error {
	// If ID exists, update. Otherwise try to find by student_id and exam_id or just use Create (if it's a new record)
	if mark.ID != 0 {
		return r.examMarkRepo.Update(ctx, mark, "id = $1 AND deleted_at IS NULL", mark.ID)
	}
	return r.examMarkRepo.Create(ctx, mark)
}

func (r *operationRepository) GetMarksByExam(ctx context.Context, examID int64) ([]*entity.ExamMark, error) {
	return r.examMarkRepo.FindAll(ctx, "exam_id = $1 AND deleted_at IS NULL", examID)
}

func (r *operationRepository) GetMarksByStudent(ctx context.Context, studentID int64) ([]*entity.ExamMark, error) {
	return r.examMarkRepo.FindAll(ctx, "student_id = $1 AND deleted_at IS NULL", studentID)
}

func (r *operationRepository) RecordStudentAttendance(ctx context.Context, attendance *entity.StudentAttendance) error {
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

func (r *operationRepository) UpdateStudentAttendance(ctx context.Context, attendance *entity.StudentAttendance) error {
	return r.studentAttendanceRepo.Update(ctx, attendance, "id = $1 AND deleted_at IS NULL", attendance.ID)
}

func (r *operationRepository) GetStudentAttendance(ctx context.Context, studentID int64, filters map[string]interface{}) ([]*entity.StudentAttendance, error) {
	return r.studentAttendanceRepo.FindAll(ctx, "student_id = $1 AND deleted_at IS NULL", studentID)
}

func (r *operationRepository) GetSectionAttendance(ctx context.Context, sectionID int64, date string, subjectID *int64) ([]*entity.StudentAttendance, error) {
	// Use a temporary struct to handle JOIN results since StudentAttendance 
	// display fields are marked with db:"-"
	type attendanceResult struct {
		entity.StudentAttendance
		StudentName   string `db:"student_name"`
		StudentNumber string `db:"student_number"`
		SubjectName   string `db:"subject_name"`
	}
	
	var results []attendanceResult
	query := `
		SELECT 
			sa.*, 
			s.full_name as student_name, 
			s.student_number,
			sub.name as subject_name
		FROM student_attendance sa
		JOIN students s ON sa.student_id = s.id
		LEFT JOIN subjects sub ON sa.subject_id = sub.id
		WHERE sa.section_id = $1 AND sa.attendance_date::DATE = $2::DATE AND sa.deleted_at IS NULL
	`
	args := []interface{}{sectionID, date}
	if subjectID != nil {
		query += " AND sa.subject_id = $3"
		args = append(args, *subjectID)
	}

	err := r.db.SelectContext(ctx, &results, query, args...)
	if err != nil {
		return nil, err
	}

	// Map back to entity.StudentAttendance
	finalResults := make([]*entity.StudentAttendance, len(results))
	for i, res := range results {
		res.StudentAttendance.StudentName = &res.StudentName
		res.StudentAttendance.StudentNumber = &res.StudentNumber
		res.StudentAttendance.SubjectName = &res.SubjectName
		finalResults[i] = &res.StudentAttendance
	}
	
	return finalResults, nil
}

func (r *operationRepository) GetAttendanceWithFilters(ctx context.Context, filters map[string]interface{}) ([]*entity.StudentAttendance, error) {
	type attendanceResult struct {
		entity.StudentAttendance
		StudentName   string `db:"student_name"`
		StudentNumber string `db:"student_number"`
		SubjectName   string `db:"subject_name"`
	}
	
	var results []attendanceResult
	query := `
		SELECT 
			sa.*, 
			s.full_name as student_name, 
			s.student_number,
			sub.name as subject_name
		FROM student_attendance sa
		JOIN students s ON sa.student_id = s.id
		LEFT JOIN subjects sub ON sa.subject_id = sub.id
		WHERE sa.deleted_at IS NULL
	`
	err := r.db.SelectContext(ctx, &results, query)
	if err != nil {
		return nil, err
	}

	finalResults := make([]*entity.StudentAttendance, len(results))
	for i, res := range results {
		res.StudentAttendance.StudentName = &res.StudentName
		res.StudentAttendance.StudentNumber = &res.StudentNumber
		res.StudentAttendance.SubjectName = &res.SubjectName
		finalResults[i] = &res.StudentAttendance
	}
	
	return finalResults, nil
}

func (r *operationRepository) RecordTeacherAttendance(ctx context.Context, attendance *entity.TeacherAttendance) error {
	// Check if record already exists for this teacher on this date
	existing, _ := r.teacherAttendanceRepo.FindOne(ctx,
		"teacher_id = $1 AND attendance_date = $2 AND deleted_at IS NULL",
		attendance.TeacherID, attendance.AttendanceDate)

	if existing != nil {
		attendance.ID = existing.ID
		return r.teacherAttendanceRepo.Update(ctx, attendance, "id = $1", existing.ID)
	}
	return r.teacherAttendanceRepo.Create(ctx, attendance)
}

func (r *operationRepository) UpdateTeacherAttendance(ctx context.Context, attendance *entity.TeacherAttendance) error {
	return r.teacherAttendanceRepo.Update(ctx, attendance, "id = $1 AND deleted_at IS NULL", attendance.ID)
}

func (r *operationRepository) GetTeacherAttendance(ctx context.Context, teacherID int64, filters map[string]interface{}) ([]*entity.TeacherAttendance, error) {
	return r.teacherAttendanceRepo.FindAll(ctx, "teacher_id = $1 AND deleted_at IS NULL", teacherID)
}

func (r *operationRepository) RecordStaffAttendance(ctx context.Context, attendance *entity.StaffAttendance) error {
	// Check if record already exists for this employee on this date
	existing, _ := r.staffAttendanceRepo.FindOne(ctx,
		"employee_id = $1 AND attendance_date = $2 AND deleted_at IS NULL",
		attendance.EmployeeID, attendance.AttendanceDate)

	if existing != nil {
		attendance.ID = existing.ID
		return r.staffAttendanceRepo.Update(ctx, attendance, "id = $1", existing.ID)
	}
	return r.staffAttendanceRepo.Create(ctx, attendance)
}

func (r *operationRepository) UpdateStaffAttendance(ctx context.Context, attendance *entity.StaffAttendance) error {
	return r.staffAttendanceRepo.Update(ctx, attendance, "id = $1 AND deleted_at IS NULL", attendance.ID)
}

func (r *operationRepository) GetStaffAttendance(ctx context.Context, employeeID int64, filters map[string]interface{}) ([]*entity.StaffAttendance, error) {
	return r.staffAttendanceRepo.FindAll(ctx, "employee_id = $1 AND deleted_at IS NULL", employeeID)
}

func (r *operationRepository) CreateNotification(ctx context.Context, notification *entity.Notification) error {
	return r.notificationRepo.Create(ctx, notification)
}

func (r *operationRepository) GetNotificationsByUser(ctx context.Context, userID int64) ([]*entity.Notification, error) {
	return r.notificationRepo.FindAll(ctx, "user_id = $1 AND deleted_at IS NULL", userID)
}

func (r *operationRepository) GetStudentMarksBySession(ctx context.Context, studentID, sessionID int64) ([]*entity.ExamMark, error) {
	var marks []*entity.ExamMark
	query := `
		SELECT em.* FROM exam_marks em
		JOIN exams e ON em.exam_id = e.id
		WHERE em.student_id = $1 AND e.academic_session_id = $2 AND em.deleted_at IS NULL AND e.deleted_at IS NULL
	`
	err := r.db.SelectContext(ctx, &marks, query, studentID, sessionID)
	return marks, err
}

func (r *operationRepository) GetStudentAttendanceBySession(ctx context.Context, studentID, sessionID int64) ([]*entity.StudentAttendance, error) {
	return r.studentAttendanceRepo.FindAll(ctx, "student_id = $1 AND academic_session_id = $2 AND deleted_at IS NULL", studentID, sessionID)
}

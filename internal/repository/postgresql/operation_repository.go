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
	GetSectionAttendance(ctx context.Context, sectionID int64, date string) ([]*entity.StudentAttendance, error)
	GetAttendanceWithFilters(ctx context.Context, filters map[string]interface{}) ([]*entity.StudentAttendance, error)
	RecordTeacherAttendance(ctx context.Context, attendance *entity.TeacherAttendance) error

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
	return r.studentAttendanceRepo.Create(ctx, attendance)
}

func (r *operationRepository) UpdateStudentAttendance(ctx context.Context, attendance *entity.StudentAttendance) error {
	return r.studentAttendanceRepo.Update(ctx, attendance, "id = $1 AND deleted_at IS NULL", attendance.ID)
}

func (r *operationRepository) GetStudentAttendance(ctx context.Context, studentID int64, filters map[string]interface{}) ([]*entity.StudentAttendance, error) {
	return r.studentAttendanceRepo.FindAll(ctx, "student_id = $1 AND deleted_at IS NULL", studentID)
}

func (r *operationRepository) GetSectionAttendance(ctx context.Context, sectionID int64, date string) ([]*entity.StudentAttendance, error) {
	return r.studentAttendanceRepo.FindAll(ctx, "section_id = $1 AND attendance_date = $2 AND deleted_at IS NULL", sectionID, date)
}

func (r *operationRepository) GetAttendanceWithFilters(ctx context.Context, filters map[string]interface{}) ([]*entity.StudentAttendance, error) {
	return r.studentAttendanceRepo.FindAll(ctx, "deleted_at IS NULL")
}

func (r *operationRepository) RecordTeacherAttendance(ctx context.Context, attendance *entity.TeacherAttendance) error {
	return r.teacherAttendanceRepo.Create(ctx, attendance)
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

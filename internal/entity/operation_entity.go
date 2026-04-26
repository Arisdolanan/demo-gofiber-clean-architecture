package entity

import "time"

type Schedule struct {
	ID                int64  `json:"id" db:"id"`
	SectionID         int64  `json:"section_id" db:"section_id"`
	SubjectID         int64  `json:"subject_id" db:"subject_id"`
	TeacherID         int64  `json:"teacher_id" db:"teacher_id"`
	AcademicSessionID int64  `json:"academic_session_id" db:"academic_session_id"`
	DayOfWeek         int    `json:"day_of_week" db:"day_of_week"` // 0=Sunday, 1=Monday, etc.
	StartTime         string `json:"start_time" db:"start_time"`
	EndTime           string `json:"end_time" db:"end_time"`
	RoomNumber        string `json:"room_number" db:"room_number"`
	BaseEntity
}

type Exam struct {
	ID                int64      `json:"id" db:"id"`
	SectionID         int64      `json:"section_id" db:"section_id"`
	SubjectID         int64      `json:"subject_id" db:"subject_id"`
	AcademicSessionID int64      `json:"academic_session_id" db:"academic_session_id"`
	Title             string     `json:"title" db:"title" validate:"required"`
	Description       string     `json:"description" db:"description"`
	ExamType          string     `json:"exam_type" db:"exam_type"` // daily, midterm, final, practice
	ExamDate          time.Time  `json:"exam_date" db:"exam_date"`
	DurationMinutes   int        `json:"duration_minutes" db:"duration_minutes"`
	MaxScore          int        `json:"max_score" db:"max_score"`
	BaseEntity
}

type ExamMark struct {
	ID        int64      `json:"id" db:"id"`
	ExamID    int64      `json:"exam_id" db:"exam_id"`
	StudentID int64      `json:"student_id" db:"student_id"`
	Score     float64    `json:"score" db:"score"`
	Notes     string     `json:"notes" db:"notes"`
	EnteredBy *int64     `json:"entered_by,omitempty" db:"entered_by"` // Teacher ID
	EnteredAt *time.Time `json:"entered_at,omitempty" db:"entered_at"`
	BaseEntity
}

type StudentAttendance struct {
	ID                int64      `json:"id" db:"id"`
	StudentID         int64      `json:"student_id" db:"student_id"`
	SectionID         int64      `json:"section_id" db:"section_id"`
	AcademicSessionID int64      `json:"academic_session_id" db:"academic_session_id"`
	AttendanceDate    time.Time  `json:"attendance_date" db:"attendance_date"`
	Status            string     `json:"status" db:"status"` // present, absent, late, sick, permission
	Notes             string     `json:"notes" db:"notes"`
	MarkedBy          *int64     `json:"marked_by,omitempty" db:"marked_by"` // Teacher ID
	BaseEntity
}

type TeacherAttendance struct {
	ID             int64      `json:"id" db:"id"`
	TeacherID      int64      `json:"teacher_id" db:"teacher_id"`
	AttendanceDate time.Time  `json:"attendance_date" db:"attendance_date"`
	CheckInTime    *time.Time `json:"check_in_time,omitempty" db:"check_in_time"`
	CheckOutTime   *time.Time `json:"check_out_time,omitempty" db:"check_out_time"`
	Status         string     `json:"status" db:"status"` // present, absent, late, sick, permission
	Notes          string     `json:"notes" db:"notes"`
	BaseEntity
}

type Message struct {
	ID          int64     `json:"id" db:"id"`
	SchoolID    int64     `json:"school_id" db:"school_id"`
	SenderID    int64     `json:"sender_id" db:"sender_id"`
	RecipientID int64     `json:"recipient_id" db:"recipient_id"`
	Subject     string    `json:"subject" db:"subject"`
	Body        string    `json:"body" db:"body"`
	IsRead      bool      `json:"is_read" db:"is_read"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Notification struct {
	ID            int64     `json:"id" db:"id"`
	UserID        int64     `json:"user_id" db:"user_id"`
	SchoolID      *int64    `json:"school_id" db:"school_id"`
	Title         string    `json:"title" db:"title"`
	Message       string    `json:"message" db:"message"`
	ReferenceType string    `json:"reference_type" db:"reference_type"`
	ReferenceID   *int64    `json:"reference_id" db:"reference_id"`
	IsRead        bool      `json:"is_read" db:"is_read"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type SubjectGrade struct {
	SubjectID   int64     `json:"subject_id"`
	SubjectName string    `json:"subject_name"`
	Grades      []float64 `json:"grades"`
	Average     float64   `json:"average"`
}

type ReportCard struct {
	StudentID         int64          `json:"student_id"`
	StudentName       string         `json:"student_name"`
	AcademicSessionID int64          `json:"academic_session_id"`
	SessionName       string         `json:"session_name"`
	SubjectGrades     []SubjectGrade `json:"subject_grades"`
	TotalAverage      float64        `json:"total_average"`
	AttendanceSummary map[string]int `json:"attendance_summary"` // Key: status, Value: count
}

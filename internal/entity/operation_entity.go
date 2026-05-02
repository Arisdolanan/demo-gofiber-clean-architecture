package entity

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// DateOnly represents a date without time for JSON and database operations
type DateOnly struct {
	time.Time
}

// UnmarshalJSON parses date in YYYY-MM-DD format
func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		d.Time = time.Time{}
		return nil
	}

	// Try parsing as full datetime first
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		d.Time = t
		return nil
	}

	// Try parsing as full datetime without timezone
	t, err = time.Parse("2006-01-02T15:04:05.999999999", s)
	if err == nil {
		d.Time = t
		return nil
	}

	// Try parsing as date only
	t, err = time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// MarshalJSON formats date as YYYY-MM-DD
func (d DateOnly) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.Time.Format("2006-01-02"))), nil
}

// Value implements driver.Valuer for database operations
func (d DateOnly) Value() (driver.Value, error) {
	return d.Time.Format("2006-01-02"), nil
}

// Scan implements sql.Scanner for database operations
func (d *DateOnly) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		d.Time = v
	case []byte:
		t, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		d.Time = t
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		d.Time = t
	default:
		return fmt.Errorf("cannot scan %T into DateOnly", value)
	}
	return nil
}

// TimeOnly represents a time without date for JSON and database operations
type TimeOnly struct {
	time.Time
}

// UnmarshalJSON parses time in HH:MM:SS format
func (t *TimeOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t.Time = time.Time{}
		return nil
	}

	// Try parsing as full datetime first
	parsed, err := time.Parse(time.RFC3339, s)
	if err == nil {
		t.Time = parsed
		return nil
	}

	// Try parsing as full datetime without timezone
	parsed, err = time.Parse("2006-01-02T15:04:05.999999999", s)
	if err == nil {
		t.Time = parsed
		return nil
	}

	// Try parsing as time only (with seconds)
	parsed, err = time.Parse("15:04:05", s)
	if err == nil {
		t.Time = parsed
		return nil
	}

	// Try parsing as time only (without seconds)
	parsed, err = time.Parse("15:04", s)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

// MarshalJSON formats time as HH:MM:SS
func (t TimeOnly) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.Time.Format("15:04:05"))), nil
}

// Value implements driver.Valuer for database operations
func (t TimeOnly) Value() (driver.Value, error) {
	return t.Time.Format("15:04:05"), nil
}

// Scan implements sql.Scanner for database operations
func (t *TimeOnly) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		t.Time = v
	case []byte:
		parsed, err := time.Parse("15:04:05", string(v))
		if err != nil {
			return err
		}
		t.Time = parsed
	case string:
		parsed, err := time.Parse("15:04:05", v)
		if err != nil {
			return err
		}
		t.Time = parsed
	default:
		return fmt.Errorf("cannot scan %T into TimeOnly", value)
	}
	return nil
}

type Schedule struct {
	ID                int64  `json:"id" db:"id"`
	SchoolID          int64  `json:"school_id" db:"school_id"`
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
	ID                int64     `json:"id" db:"id"`
	SchoolID          int64     `json:"school_id" db:"school_id"`
	SectionID         int64     `json:"section_id" db:"section_id"`
	SubjectID         int64     `json:"subject_id" db:"subject_id"`
	AcademicSessionID int64     `json:"academic_session_id" db:"academic_session_id"`
	Title             string    `json:"title" db:"title" validate:"required"`
	Description       string    `json:"description" db:"description"`
	ExamType          string    `json:"exam_type" db:"exam_type"` // daily, midterm, final, practice
	ExamDate          time.Time `json:"exam_date" db:"exam_date"`
	DurationMinutes   int       `json:"duration_minutes" db:"duration_minutes"`
	MaxScore          int       `json:"max_score" db:"max_score"`
	BaseEntity
}

type ExamMark struct {
	ID        int64      `json:"id" db:"id"`
	SchoolID  int64      `json:"school_id" db:"school_id"`
	ExamID    int64      `json:"exam_id" db:"exam_id"`
	StudentID int64      `json:"student_id" db:"student_id"`
	Score     float64    `json:"score" db:"score"`
	Notes     string     `json:"notes" db:"notes"`
	EnteredBy *int64     `json:"entered_by,omitempty" db:"entered_by"` // Teacher ID
	EnteredAt *time.Time `json:"entered_at,omitempty" db:"entered_at"`
	BaseEntity
}

type StudentAttendance struct {
	ID                int64     `json:"id" db:"id"`
	SchoolID          int64     `json:"school_id" db:"school_id"`
	StudentID         int64     `json:"student_id" db:"student_id"`
	SectionID         int64     `json:"section_id" db:"section_id"`
	AcademicSessionID int64     `json:"academic_session_id" db:"academic_session_id"`
	SubjectID         *int64    `json:"subject_id,omitempty" db:"subject_id"`
	AttendanceDate    DateOnly  `json:"attendance_date" db:"attendance_date"`
	CheckInTime       *TimeOnly `json:"check_in_time,omitempty" db:"check_in_time"`
	CheckOutTime      *TimeOnly `json:"check_out_time,omitempty" db:"check_out_time"`
	Status            string    `json:"status" db:"status"` // present, absent, late, sick, permission
	Notes             *string   `json:"notes,omitempty" db:"notes"`
	MarkedBy          *int64    `json:"marked_by,omitempty" db:"marked_by"` // Teacher ID
	CheckInLocation   *string   `json:"check_in_location,omitempty" db:"check_in_location"`
	CheckOutLocation  *string   `json:"check_out_location,omitempty" db:"check_out_location"`
	CheckInIPAddress  *string   `json:"check_in_ip_address,omitempty" db:"check_in_ip_address"`
	CheckOutIPAddress *string   `json:"check_out_ip_address,omitempty" db:"check_out_ip_address"`
	CheckInDevice     *string   `json:"check_in_device,omitempty" db:"check_in_device"`
	CheckOutDevice    *string   `json:"check_out_device,omitempty" db:"check_out_device"`

	// Display fields (populated via JOIN)
	StudentName   *string `json:"student_name,omitempty" db:"-"`
	StudentNumber *string `json:"student_number,omitempty" db:"-"`
	SubjectName   *string `json:"subject_name,omitempty" db:"-"`
	ClassName     *string `json:"class_name,omitempty" db:"-"`
	SectionName   *string `json:"section_name,omitempty" db:"-"`

	BaseEntity
}

// StudentAttendanceEntry represents a single student's attendance entry in a bulk request
type StudentAttendanceEntry struct {
	StudentID int64   `json:"student_id"`
	Status    string  `json:"status"` // present, absent, late, sick, permission
	Notes     *string `json:"notes,omitempty"`
}

// SectionAttendanceRequest is used by teachers to record attendance for a whole section.
// subject_id is required - it is mandatory supporting data, not a filter.
type SectionAttendanceRequest struct {
	SectionID         int64                    `json:"section_id"`
	AcademicSessionID int64                    `json:"academic_session_id"`
	SubjectID         int64                    `json:"subject_id"` // REQUIRED: must be filled by teacher before saving
	AttendanceDate    DateOnly                 `json:"attendance_date"`
	MarkedBy          *int64                   `json:"marked_by,omitempty"` // Teacher ID
	Attendances       []StudentAttendanceEntry `json:"attendances"`
	IPAddress         string                   `json:"ip_address"`
	UserAgent         string                   `json:"user_agent"`
}

type TeacherAttendance struct {
	ID                int64     `json:"id" db:"id"`
	SchoolID          int64     `json:"school_id" db:"school_id"`
	TeacherID         int64     `json:"teacher_id" db:"teacher_id"`
	AttendanceDate    DateOnly  `json:"attendance_date" db:"attendance_date"`
	CheckInTime       *TimeOnly `json:"check_in_time,omitempty" db:"check_in_time"`
	CheckOutTime      *TimeOnly `json:"check_out_time,omitempty" db:"check_out_time"`
	CheckInLocation   *string   `json:"check_in_location,omitempty" db:"check_in_location"`
	CheckOutLocation  *string   `json:"check_out_location,omitempty" db:"check_out_location"`
	CheckInIPAddress  *string   `json:"check_in_ip_address,omitempty" db:"check_in_ip_address"`
	CheckOutIPAddress *string   `json:"check_out_ip_address,omitempty" db:"check_out_ip_address"`
	CheckInDevice     *string   `json:"check_in_device,omitempty" db:"check_in_device"`
	CheckOutDevice    *string   `json:"check_out_device,omitempty" db:"check_out_device"`
	Status            string    `json:"status" db:"status"` // present, absent, late, sick, permission
	Notes             *string   `json:"notes,omitempty" db:"notes"`
	BaseEntity
}

type StaffAttendance struct {
	ID                int64     `json:"id" db:"id"`
	SchoolID          int64     `json:"school_id" db:"school_id"`
	EmployeeID        int64     `json:"employee_id" db:"employee_id"`
	AttendanceDate    DateOnly  `json:"attendance_date" db:"attendance_date"`
	CheckInTime       *TimeOnly `json:"check_in_time,omitempty" db:"check_in_time"`
	CheckOutTime      *TimeOnly `json:"check_out_time,omitempty" db:"check_out_time"`
	CheckInLocation   *string   `json:"check_in_location,omitempty" db:"check_in_location"`
	CheckOutLocation  *string   `json:"check_out_location,omitempty" db:"check_out_location"`
	CheckInIPAddress  *string   `json:"check_in_ip_address,omitempty" db:"check_in_ip_address"`
	CheckOutIPAddress *string   `json:"check_out_ip_address,omitempty" db:"check_out_ip_address"`
	CheckInDevice     *string   `json:"check_in_device,omitempty" db:"check_in_device"`
	CheckOutDevice    *string   `json:"check_out_device,omitempty" db:"check_out_device"`
	Status            string    `json:"status" db:"status"` // present, absent, late, sick, permission
	Notes             *string   `json:"notes,omitempty" db:"notes"`
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
	ID            int64      `json:"id" db:"id"`
	UserID        int64      `json:"user_id" db:"user_id"`
	SchoolID      *int64     `json:"school_id" db:"school_id"`
	Title         string     `json:"title" db:"title"`
	Message       string     `json:"message" db:"message"`
	ReferenceType string     `json:"reference_type" db:"reference_type"`
	ReferenceID   *int64     `json:"reference_id" db:"reference_id"`
	IsRead        bool       `json:"is_read" db:"is_read"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	DeletedAt     *time.Time `json:"deleted_at" db:"deleted_at"`
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

type IntegrationDefinition struct {
	ID          int64   `json:"id" db:"id"`
	SchoolID    int64   `json:"school_id" db:"school_id"`
	Code        string  `json:"code" db:"code"`
	Name        string  `json:"name" db:"name"`
	Provider    *string `json:"provider,omitempty" db:"provider"`
	Category    string  `json:"category" db:"category"`
	Description *string `json:"description,omitempty" db:"description"`
	IsSystem    bool    `json:"is_system" db:"is_system"`
	BaseEntity
}

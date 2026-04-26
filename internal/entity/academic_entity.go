package entity

import "time"

type AcademicSession struct {
	ID        int64     `json:"id" db:"id"`
	SchoolID  int64     `json:"school_id" db:"school_id"`
	Name      string    `json:"name" db:"name" validate:"required"`
	Code      string    `json:"code" db:"code" validate:"required"`
	StartDate time.Time `json:"start_date" db:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" db:"end_date" validate:"required"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	BaseEntity
}

type Class struct {
	ID          int64       `json:"id" db:"id"`
	SchoolID    int64       `json:"school_id" db:"school_id"`
	Name        string      `json:"name" db:"name" validate:"required"`
	Code        string      `json:"code" db:"code" validate:"required"`
	Level       SchoolLevel `json:"level" db:"level" validate:"required"`
	GradeNumber int         `json:"grade_number" db:"grade_number" validate:"required"`
	Description string      `json:"description" db:"description"`
	BaseEntity
}

type Subject struct {
	ID          int64  `json:"id" db:"id"`
	SchoolID    int64  `json:"school_id" db:"school_id"`
	Name        string `json:"name" db:"name" validate:"required"`
	Code        string `json:"code" db:"code" validate:"required"`
	Description string `json:"description" db:"description"`
	CreditHours int    `json:"credit_hours" db:"credit_hours"`
	BaseEntity
}

type Section struct {
	ID                int64   `json:"id" db:"id"`
	ClassID           int64   `json:"class_id" db:"class_id"`
	AcademicSessionID int64   `json:"academic_session_id" db:"academic_session_id"`
	Name              string  `json:"name" db:"name" validate:"required"`
	Code              string  `json:"code" db:"code"`
	RoomNumber        string  `json:"room_number" db:"room_number"`
	Capacity          int     `json:"capacity" db:"capacity"`
	TeacherID         *int64  `json:"teacher_id,omitempty" db:"teacher_id"` // Homeroom teacher
	BaseEntity
}

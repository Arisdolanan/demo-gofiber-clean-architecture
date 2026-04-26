package entity

import "time"

type TeacherStatus string

const (
	TeacherActive   TeacherStatus = "active"
	TeacherInactive TeacherStatus = "inactive"
)

type Teacher struct {
	ID             int64          `json:"id" db:"id"`
	UserID         int64          `json:"user_id" db:"user_id"`
	SchoolID       int64          `json:"school_id" db:"school_id"`
	EmployeeNumber string         `json:"employee_number" db:"employee_number" validate:"required"`
	FullName       string         `json:"full_name" db:"full_name" validate:"required"`
	DateOfBirth    *time.Time     `json:"date_of_birth,omitempty" db:"date_of_birth"`
	Gender         string         `json:"gender" db:"gender"`
	Phone          string         `json:"phone" db:"phone"`
	Email          string         `json:"email" db:"email"`
	Address        string         `json:"address" db:"address"`
	Qualification  string         `json:"qualification" db:"qualification"`
	Specialization string         `json:"specialization" db:"specialization"`
	JoinDate       *time.Time     `json:"join_date,omitempty" db:"join_date"`
	Status         TeacherStatus  `json:"status" db:"status"`
	BaseEntity
}

type StudentStatus string

const (
	StudentActive      StudentStatus = "active"
	StudentInactive    StudentStatus = "inactive"
	StudentGraduated   StudentStatus = "graduated"
	StudentTransferred StudentStatus = "transferred"
)

type Student struct {
	ID            int64          `json:"id" db:"id"`
	UserID        *int64         `json:"user_id,omitempty" db:"user_id"`
	SchoolID      int64          `json:"school_id" db:"school_id"`
	StudentNumber string         `json:"student_number" db:"student_number" validate:"required"`
	FullName      string         `json:"full_name" db:"full_name" validate:"required"`
	DateOfBirth   *time.Time     `json:"date_of_birth,omitempty" db:"date_of_birth"`
	Gender        string         `json:"gender" db:"gender"`
	BloodType     string         `json:"blood_type" db:"blood_type"`
	Phone         string         `json:"phone" db:"phone"`
	Email         string         `json:"email" db:"email"`
	Address       string         `json:"address" db:"address"`
	AdmissionDate *time.Time     `json:"admission_date,omitempty" db:"admission_date"`
	Status        StudentStatus  `json:"status" db:"status"`
	BaseEntity
}

type Parent struct {
	ID         int64  `json:"id" db:"id"`
	UserID     int64  `json:"user_id" db:"user_id"`
	SchoolID   int64  `json:"school_id" db:"school_id"`
	FullName   string `json:"full_name" db:"full_name" validate:"required"`
	Phone      string `json:"phone" db:"phone"`
	Email      string `json:"email" db:"email"`
	Address    string `json:"address" db:"address"`
	Occupation string `json:"occupation" db:"occupation"`
	BaseEntity
}

type StudentParent struct {
	ID           int64  `json:"id" db:"id"`
	StudentID    int64  `json:"student_id" db:"student_id"`
	ParentID     int64  `json:"parent_id" db:"parent_id"`
	Relationship string `json:"relationship" db:"relationship"` // father, mother, guardian
	IsPrimary    bool   `json:"is_primary" db:"is_primary"`
	BaseEntity
}

type StudentSection struct {
	ID                int64      `json:"id" db:"id"`
	StudentID         int64      `json:"student_id" db:"student_id"`
	SectionID         int64      `json:"section_id" db:"section_id"`
	AcademicSessionID int64      `json:"academic_session_id" db:"academic_session_id"`
	RollNumber        string     `json:"roll_number" db:"roll_number"`
	EnrollmentDate    *time.Time `json:"enrollment_date,omitempty" db:"enrollment_date"`
	Status            string     `json:"status" db:"status"` // active, promoted, transferred
	BaseEntity
}

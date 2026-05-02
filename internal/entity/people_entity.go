package entity

import "time"

type TeacherStatus string

const (
	TeacherActive   TeacherStatus = "active"
	TeacherInactive TeacherStatus = "inactive"
)

type Teacher struct {
	ID             int64         `json:"id" db:"id"`
	UserID         int64         `json:"user_id" db:"user_id"`
	SchoolID       int64         `json:"school_id" db:"school_id"`
	EmployeeNumber string        `json:"employee_number" db:"employee_number" validate:"required"`
	FullName       *string       `json:"full_name,omitempty" db:"full_name" validate:"required"`
	DateOfBirth    *time.Time    `json:"date_of_birth,omitempty" db:"date_of_birth"`
	Gender         *string       `json:"gender,omitempty" db:"gender"`
	Phone          *string       `json:"phone,omitempty" db:"phone"`
	Email          *string       `json:"email,omitempty" db:"email"`
	Address        *string       `json:"address,omitempty" db:"address"`
	Qualification  *string       `json:"qualification,omitempty" db:"qualification"`
	Specialization *string       `json:"specialization,omitempty" db:"specialization"`
	JoinDate       *time.Time    `json:"join_date,omitempty" db:"join_date"`
	Status         TeacherStatus `json:"status" db:"status"`
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
	ID            int64         `json:"id" db:"id"`
	UserID        *int64        `json:"user_id,omitempty" db:"user_id"`
	SchoolID      int64         `json:"school_id" db:"school_id"`
	StudentNumber string        `json:"student_number" db:"student_number" validate:"required"`
	NIS           *string       `json:"nis,omitempty" db:"nis"`
	NISN          *string       `json:"nisn,omitempty" db:"nisn"`
	FullName      *string       `json:"full_name,omitempty" db:"full_name" validate:"required"`
	DateOfBirth   *time.Time    `json:"date_of_birth,omitempty" db:"date_of_birth"`
	Gender        *string       `json:"gender,omitempty" db:"gender"`
	BloodType     *string       `json:"blood_type,omitempty" db:"blood_type"`
	Phone         *string       `json:"phone,omitempty" db:"phone"`
	Email         *string       `json:"email,omitempty" db:"email"`
	Address       *string       `json:"address,omitempty" db:"address"`
	AdmissionDate *time.Time    `json:"admission_date,omitempty" db:"admission_date"`
	Status            StudentStatus `json:"status" db:"status"`
	SectionID         *int64        `json:"section_id,omitempty" db:"-"`
	AcademicSessionID *int64        `json:"academic_session_id,omitempty" db:"-"`
	Parents           []StudentParentRequest `json:"parents,omitempty" db:"-"`
	BaseEntity
}

type StudentParentRequest struct {
	ParentID     int64   `json:"parent_id,omitempty" db:"parent_id"`
	Relationship string  `json:"relationship" db:"relationship"`
	IsPrimary    bool    `json:"is_primary" db:"is_primary"`
	// Fields for new parent creation if ParentID is 0
	FullName     *string `json:"full_name,omitempty" db:"full_name"`
	Phone        *string `json:"phone,omitempty" db:"phone"`
	Email        *string `json:"email,omitempty" db:"email"`
	Address      *string `json:"address,omitempty" db:"address"`
	Occupation   *string `json:"occupation,omitempty" db:"occupation"`
}

type Parent struct {
	ID         int64         `json:"id" db:"id"`
	UserID     *int64        `json:"user_id" db:"user_id"`
	SchoolID   int64         `json:"school_id" db:"school_id"`
	FullName   *string       `json:"full_name,omitempty" db:"full_name" validate:"required"`
	Phone      *string       `json:"phone,omitempty" db:"phone"`
	Email      *string       `json:"email,omitempty" db:"email"`
	Address    *string       `json:"address,omitempty" db:"address"`
	Occupation *string       `json:"occupation,omitempty" db:"occupation"`
	Children   []ParentChild `json:"children,omitempty" db:"-"`
	BaseEntity
}

type ParentChild struct {
	ID            int64  `json:"id" db:"id"`
	FullName      string `json:"full_name" db:"full_name"`
	StudentNumber string `json:"student_number" db:"student_number"`
	Relationship  string `json:"relationship" db:"relationship"`
	IsPrimary     bool   `json:"is_primary" db:"is_primary"`
}

type StudentParent struct {
	ID           int64  `json:"id" db:"id"`
	SchoolID     int64  `json:"school_id" db:"school_id"`
	StudentID    int64  `json:"student_id" db:"student_id"`
	ParentID     int64  `json:"parent_id" db:"parent_id"`
	Relationship string `json:"relationship" db:"relationship"` // father, mother, guardian
	IsPrimary    bool   `json:"is_primary" db:"is_primary"`
}

type StudentSection struct {
	ID                int64      `json:"id" db:"id"`
	SchoolID          int64      `json:"school_id" db:"school_id"`
	StudentID         int64      `json:"student_id" db:"student_id"`
	SectionID         int64      `json:"section_id" db:"section_id"`
	AcademicSessionID int64      `json:"academic_session_id" db:"academic_session_id"`
	RollNumber        *string    `json:"roll_number,omitempty" db:"roll_number"`
	EnrollmentDate    *time.Time `json:"enrollment_date,omitempty" db:"enrollment_date"`
	Status            *string    `json:"status,omitempty" db:"status"` // active, promoted, transferred
}

type StaffStatus string

const (
	StaffActive   StaffStatus = "active"
	StaffInactive StaffStatus = "inactive"
)

type Staff struct {
	ID             int64       `json:"id" db:"id"`
	UserID         *int64      `json:"user_id,omitempty" db:"user_id"`
	SchoolID       int64       `json:"school_id" db:"school_id"`
	EmployeeNumber string      `json:"employee_number" db:"employee_number" validate:"required"`
	FullName       *string     `json:"full_name,omitempty" db:"full_name" validate:"required"`
	DateOfBirth    *time.Time  `json:"date_of_birth,omitempty" db:"date_of_birth"`
	Gender         *string     `json:"gender,omitempty" db:"gender"`
	Phone          *string     `json:"phone,omitempty" db:"phone"`
	Email          *string     `json:"email,omitempty" db:"email"`
	Address        *string     `json:"address,omitempty" db:"address"`
	Position       *string     `json:"position,omitempty" db:"position"`
	Department     *string     `json:"department,omitempty" db:"department"`
	JoinDate       *time.Time  `json:"join_date,omitempty" db:"join_date"`
	Status         StaffStatus `json:"status" db:"status"`
	BaseEntity
}

type EmployeeUserResponse struct {
	EmployeeID int64  `json:"employee_id,omitempty"`
	UserID     int64  `json:"user_id"`
	Username   string `json:"username"`
	Password   string `json:"password,omitempty"`
	UserType   string `json:"user_type,omitempty"` // "teacher", "staff", "student", "parent"
}

type EnrollmentCredentials struct {
	Student *EmployeeUserResponse   `json:"student"`
	Parents []*EmployeeUserResponse `json:"parents"`
}

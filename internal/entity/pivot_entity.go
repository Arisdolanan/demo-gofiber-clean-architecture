package entity

// ClassSubject represents the many-to-many relationship between classes and subjects
type ClassSubject struct {
	ID                int64 `json:"id" db:"id"`
	SchoolID          int64 `json:"school_id" db:"school_id"`
	ClassID           int64 `json:"class_id" db:"class_id"`
	SubjectID         int64 `json:"subject_id" db:"subject_id"`
	AcademicSessionID int64 `json:"academic_session_id" db:"academic_session_id"`
	BaseEntity
}

// TeacherSubject represents the assignment of teachers to subjects in sections
type TeacherSubject struct {
	ID                int64 `json:"id" db:"id"`
	SchoolID          int64 `json:"school_id" db:"school_id"`
	TeacherID         int64 `json:"teacher_id" db:"teacher_id"`
	SectionID         int64 `json:"section_id" db:"section_id"`
	SubjectID         int64 `json:"subject_id" db:"subject_id"`
	AcademicSessionID int64 `json:"academic_session_id" db:"academic_session_id"`
	BaseEntity
}

package postgresql

import (
	"context"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/jmoiron/sqlx"
)

type AcademicRepository interface {
	// Sessions
	CreateSession(ctx context.Context, session *entity.AcademicSession) error
	UpdateSession(ctx context.Context, session *entity.AcademicSession) error
	DeleteSession(ctx context.Context, schoolID, id int64) error
	FindSessionByID(ctx context.Context, schoolID, id int64) (*entity.AcademicSession, error)
	FindSessionsBySchool(ctx context.Context, schoolID int64) ([]*entity.AcademicSession, error)
	FindActiveSession(ctx context.Context, schoolID int64) (*entity.AcademicSession, error)

	// Classes
	CreateClass(ctx context.Context, class *entity.Class) error
	UpdateClass(ctx context.Context, class *entity.Class) error
	DeleteClass(ctx context.Context, schoolID, id int64) error
	FindClassByID(ctx context.Context, schoolID, id int64) (*entity.Class, error)
	FindClassesBySchool(ctx context.Context, schoolID int64) ([]*entity.Class, error)

	// Sections
	CreateSection(ctx context.Context, section *entity.Section) error
	UpdateSection(ctx context.Context, section *entity.Section) error
	DeleteSection(ctx context.Context, schoolID, id int64) error
	FindSectionByID(ctx context.Context, schoolID, id int64) (*entity.Section, error)
	FindSectionsByClass(ctx context.Context, schoolID, classID int64) ([]*entity.Section, error)
	FindSectionsBySession(ctx context.Context, schoolID, sessionID int64) ([]*entity.Section, error)
	FindAllSections(ctx context.Context, schoolID int64) ([]*entity.Section, error)

	// Subjects
	CreateSubject(ctx context.Context, subject *entity.Subject) error
	UpdateSubject(ctx context.Context, subject *entity.Subject) error
	DeleteSubject(ctx context.Context, schoolID, id int64) error
	FindSubjectByID(ctx context.Context, schoolID, id int64) (*entity.Subject, error)
	FindSubjectsBySchool(ctx context.Context, schoolID int64) ([]*entity.Subject, error)
}

type academicRepository struct {
	sessionRepo *BaseRepository[entity.AcademicSession]
	classRepo   *BaseRepository[entity.Class]
	sectionRepo *BaseRepository[entity.Section]
	subjectRepo *BaseRepository[entity.Subject]
}

func NewAcademicRepository(db *sqlx.DB) AcademicRepository {
	return &academicRepository{
		sessionRepo: NewBaseRepository[entity.AcademicSession](db, "academic_sessions"),
		classRepo:   NewBaseRepository[entity.Class](db, "classes"),
		sectionRepo: NewBaseRepository[entity.Section](db, "sections"),
		subjectRepo: NewBaseRepository[entity.Subject](db, "subjects"),
	}
}

func (r *academicRepository) CreateSession(ctx context.Context, session *entity.AcademicSession) error {
	return r.sessionRepo.Create(ctx, session)
}

func (r *academicRepository) UpdateSession(ctx context.Context, session *entity.AcademicSession) error {
	return r.sessionRepo.Update(ctx, session, "id = $1 AND deleted_at IS NULL", session.ID)
}

func (r *academicRepository) DeleteSession(ctx context.Context, schoolID, id int64) error {
	return r.sessionRepo.SoftDelete(ctx, "id = $1 AND school_id = $2 AND deleted_at IS NULL", id, schoolID)
}

func (r *academicRepository) FindSessionByID(ctx context.Context, schoolID, id int64) (*entity.AcademicSession, error) {
	return r.sessionRepo.FindOne(ctx, "id = $1 AND school_id = $2 AND deleted_at IS NULL", id, schoolID)
}

func (r *academicRepository) FindSessionsBySchool(ctx context.Context, schoolID int64) ([]*entity.AcademicSession, error) {
	return r.sessionRepo.FindAll(ctx, "school_id = $1 AND deleted_at IS NULL", schoolID)
}

func (r *academicRepository) FindActiveSession(ctx context.Context, schoolID int64) (*entity.AcademicSession, error) {
	return r.sessionRepo.FindOne(ctx, "school_id = $1 AND is_active = true AND deleted_at IS NULL", schoolID)
}

func (r *academicRepository) CreateClass(ctx context.Context, class *entity.Class) error {
	return r.classRepo.Create(ctx, class)
}

func (r *academicRepository) UpdateClass(ctx context.Context, class *entity.Class) error {
	return r.classRepo.Update(ctx, class, "id = $1 AND school_id = $2 AND deleted_at IS NULL", class.ID, class.SchoolID)
}

func (r *academicRepository) DeleteClass(ctx context.Context, schoolID, id int64) error {
	return r.classRepo.SoftDelete(ctx, "id = $1 AND school_id = $2 AND deleted_at IS NULL", id, schoolID)
}

func (r *academicRepository) FindClassByID(ctx context.Context, schoolID, id int64) (*entity.Class, error) {
	return r.classRepo.FindOne(ctx, "id = $1 AND school_id = $2 AND deleted_at IS NULL", id, schoolID)
}

func (r *academicRepository) FindClassesBySchool(ctx context.Context, schoolID int64) ([]*entity.Class, error) {
	return r.classRepo.FindAll(ctx, "school_id = $1 AND deleted_at IS NULL", schoolID)
}

func (r *academicRepository) CreateSection(ctx context.Context, section *entity.Section) error {
	return r.sectionRepo.Create(ctx, section)
}

func (r *academicRepository) UpdateSection(ctx context.Context, section *entity.Section) error {
	// Sections don't have school_id directly, but we should verify via class
	return r.sectionRepo.Update(ctx, section, "id = $1 AND deleted_at IS NULL", section.ID)
}

func (r *academicRepository) DeleteSection(ctx context.Context, schoolID, id int64) error {
	// Verify section belongs to school via class
	query := `id = $1 AND class_id IN (SELECT id FROM classes WHERE school_id = $2) AND deleted_at IS NULL`
	return r.sectionRepo.SoftDelete(ctx, query, id, schoolID)
}

func (r *academicRepository) FindSectionByID(ctx context.Context, schoolID, id int64) (*entity.Section, error) {
	query := `id = $1 AND class_id IN (SELECT id FROM classes WHERE school_id = $2) AND deleted_at IS NULL`
	return r.sectionRepo.FindOne(ctx, query, id, schoolID)
}

func (r *academicRepository) FindSectionsByClass(ctx context.Context, schoolID, classID int64) ([]*entity.Section, error) {
	query := `class_id = $1 AND class_id IN (SELECT id FROM classes WHERE school_id = $2) AND deleted_at IS NULL`
	return r.sectionRepo.FindAll(ctx, query, classID, schoolID)
}

func (r *academicRepository) FindSectionsBySession(ctx context.Context, schoolID, sessionID int64) ([]*entity.Section, error) {
	query := `academic_session_id = $1 AND academic_session_id IN (SELECT id FROM academic_sessions WHERE school_id = $2) AND deleted_at IS NULL`
	return r.sectionRepo.FindAll(ctx, query, sessionID, schoolID)
}

func (r *academicRepository) FindAllSections(ctx context.Context, schoolID int64) ([]*entity.Section, error) {
	query := `class_id IN (SELECT id FROM classes WHERE school_id = $1) AND deleted_at IS NULL`
	return r.sectionRepo.FindAll(ctx, query, schoolID)
}

func (r *academicRepository) CreateSubject(ctx context.Context, subject *entity.Subject) error {
	return r.subjectRepo.Create(ctx, subject)
}

func (r *academicRepository) UpdateSubject(ctx context.Context, subject *entity.Subject) error {
	return r.subjectRepo.Update(ctx, subject, "id = $1 AND school_id = $2 AND deleted_at IS NULL", subject.ID, subject.SchoolID)
}

func (r *academicRepository) DeleteSubject(ctx context.Context, schoolID, id int64) error {
	return r.subjectRepo.SoftDelete(ctx, "id = $1 AND school_id = $2 AND deleted_at IS NULL", id, schoolID)
}

func (r *academicRepository) FindSubjectByID(ctx context.Context, schoolID, id int64) (*entity.Subject, error) {
	return r.subjectRepo.FindOne(ctx, "id = $1 AND school_id = $2 AND deleted_at IS NULL", id, schoolID)
}

func (r *academicRepository) FindSubjectsBySchool(ctx context.Context, schoolID int64) ([]*entity.Subject, error) {
	return r.subjectRepo.FindAll(ctx, "school_id = $1 AND deleted_at IS NULL", schoolID)
}

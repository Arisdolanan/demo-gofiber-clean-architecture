package usecase

import (
	"context"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/repository/postgresql"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type AcademicUsecase interface {
	// Sessions
	CreateSession(ctx context.Context, schoolID int64, session *entity.AcademicSession) error
	UpdateSession(ctx context.Context, schoolID int64, session *entity.AcademicSession) error
	DeleteSession(ctx context.Context, schoolID, id int64) error
	GetSessionByID(ctx context.Context, schoolID, id int64) (*entity.AcademicSession, error)
	GetSessionsBySchool(ctx context.Context, schoolID int64) ([]*entity.AcademicSession, error)
	GetActiveSession(ctx context.Context, schoolID int64) (*entity.AcademicSession, error)

	// Classes
	CreateClass(ctx context.Context, schoolID int64, class *entity.Class) error
	UpdateClass(ctx context.Context, schoolID int64, class *entity.Class) error
	DeleteClass(ctx context.Context, schoolID, id int64) error
	GetClassByID(ctx context.Context, schoolID, id int64) (*entity.Class, error)
	GetClassesBySchool(ctx context.Context, schoolID int64) ([]*entity.Class, error)

	// Sections
	CreateSection(ctx context.Context, schoolID int64, section *entity.Section) error
	UpdateSection(ctx context.Context, schoolID int64, section *entity.Section) error
	DeleteSection(ctx context.Context, schoolID, id int64) error
	GetSectionByID(ctx context.Context, schoolID, id int64) (*entity.Section, error)
	GetSectionsByClass(ctx context.Context, schoolID, classID int64) ([]*entity.Section, error)
	GetSectionsBySession(ctx context.Context, schoolID, sessionID int64) ([]*entity.Section, error)
	GetAllSections(ctx context.Context, schoolID int64) ([]*entity.Section, error)

	// Subjects
	CreateSubject(ctx context.Context, schoolID int64, subject *entity.Subject) error
	UpdateSubject(ctx context.Context, schoolID int64, subject *entity.Subject) error
	DeleteSubject(ctx context.Context, schoolID, id int64) error
	GetSubjectByID(ctx context.Context, schoolID, id int64) (*entity.Subject, error)
	GetSubjectsBySchool(ctx context.Context, schoolID int64) ([]*entity.Subject, error)
}

type academicUsecase struct {
	repo     postgresql.AcademicRepository
	validate *validator.Validate
	log      *logrus.Logger
}

func NewAcademicUsecase(repo postgresql.AcademicRepository, validate *validator.Validate, log *logrus.Logger) AcademicUsecase {
	return &academicUsecase{
		repo:     repo,
		validate: validate,
		log:      log,
	}
}

func (uc *academicUsecase) CreateSession(ctx context.Context, schoolID int64, session *entity.AcademicSession) error {
	session.SchoolID = schoolID
	if err := uc.validate.Struct(session); err != nil {
		return err
	}
	return uc.repo.CreateSession(ctx, session)
}

func (uc *academicUsecase) UpdateSession(ctx context.Context, schoolID int64, session *entity.AcademicSession) error {
	session.SchoolID = schoolID
	if err := uc.validate.Struct(session); err != nil {
		return err
	}
	return uc.repo.UpdateSession(ctx, session)
}

func (uc *academicUsecase) DeleteSession(ctx context.Context, schoolID, id int64) error {
	return uc.repo.DeleteSession(ctx, schoolID, id)
}

func (uc *academicUsecase) GetSessionByID(ctx context.Context, schoolID, id int64) (*entity.AcademicSession, error) {
	return uc.repo.FindSessionByID(ctx, schoolID, id)
}

func (uc *academicUsecase) GetSessionsBySchool(ctx context.Context, schoolID int64) ([]*entity.AcademicSession, error) {
	return uc.repo.FindSessionsBySchool(ctx, schoolID)
}

func (uc *academicUsecase) GetActiveSession(ctx context.Context, schoolID int64) (*entity.AcademicSession, error) {
	return uc.repo.FindActiveSession(ctx, schoolID)
}

func (uc *academicUsecase) CreateClass(ctx context.Context, schoolID int64, class *entity.Class) error {
	class.SchoolID = schoolID
	if err := uc.validate.Struct(class); err != nil {
		return err
	}
	return uc.repo.CreateClass(ctx, class)
}

func (uc *academicUsecase) UpdateClass(ctx context.Context, schoolID int64, class *entity.Class) error {
	class.SchoolID = schoolID
	if err := uc.validate.Struct(class); err != nil {
		return err
	}
	return uc.repo.UpdateClass(ctx, class)
}

func (uc *academicUsecase) DeleteClass(ctx context.Context, schoolID, id int64) error {
	return uc.repo.DeleteClass(ctx, schoolID, id)
}

func (uc *academicUsecase) GetClassByID(ctx context.Context, schoolID, id int64) (*entity.Class, error) {
	return uc.repo.FindClassByID(ctx, schoolID, id)
}

func (uc *academicUsecase) GetClassesBySchool(ctx context.Context, schoolID int64) ([]*entity.Class, error) {
	return uc.repo.FindClassesBySchool(ctx, schoolID)
}

func (uc *academicUsecase) CreateSection(ctx context.Context, schoolID int64, section *entity.Section) error {
	section.SchoolID = schoolID
	if err := uc.validate.Struct(section); err != nil {
		return err
	}
	return uc.repo.CreateSection(ctx, section)
}

func (uc *academicUsecase) UpdateSection(ctx context.Context, schoolID int64, section *entity.Section) error {
	section.SchoolID = schoolID
	if err := uc.validate.Struct(section); err != nil {
		return err
	}
	return uc.repo.UpdateSection(ctx, section)
}

func (uc *academicUsecase) DeleteSection(ctx context.Context, schoolID, id int64) error {
	return uc.repo.DeleteSection(ctx, schoolID, id)
}

func (uc *academicUsecase) GetSectionByID(ctx context.Context, schoolID, id int64) (*entity.Section, error) {
	return uc.repo.FindSectionByID(ctx, schoolID, id)
}

func (uc *academicUsecase) GetSectionsByClass(ctx context.Context, schoolID, classID int64) ([]*entity.Section, error) {
	return uc.repo.FindSectionsByClass(ctx, schoolID, classID)
}

func (uc *academicUsecase) GetSectionsBySession(ctx context.Context, schoolID, sessionID int64) ([]*entity.Section, error) {
	return uc.repo.FindSectionsBySession(ctx, schoolID, sessionID)
}

func (uc *academicUsecase) GetAllSections(ctx context.Context, schoolID int64) ([]*entity.Section, error) {
	return uc.repo.FindAllSections(ctx, schoolID)
}

func (uc *academicUsecase) CreateSubject(ctx context.Context, schoolID int64, subject *entity.Subject) error {
	subject.SchoolID = schoolID
	if err := uc.validate.Struct(subject); err != nil {
		return err
	}
	return uc.repo.CreateSubject(ctx, subject)
}

func (uc *academicUsecase) UpdateSubject(ctx context.Context, schoolID int64, subject *entity.Subject) error {
	subject.SchoolID = schoolID
	if err := uc.validate.Struct(subject); err != nil {
		return err
	}
	return uc.repo.UpdateSubject(ctx, subject)
}

func (uc *academicUsecase) DeleteSubject(ctx context.Context, schoolID, id int64) error {
	return uc.repo.DeleteSubject(ctx, schoolID, id)
}

func (uc *academicUsecase) GetSubjectByID(ctx context.Context, schoolID, id int64) (*entity.Subject, error) {
	return uc.repo.FindSubjectByID(ctx, schoolID, id)
}

func (uc *academicUsecase) GetSubjectsBySchool(ctx context.Context, schoolID int64) ([]*entity.Subject, error) {
	return uc.repo.FindSubjectsBySchool(ctx, schoolID)
}

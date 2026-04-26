package usecase

import (
	"context"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/repository/postgresql"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type PeopleUsecase interface {
	// Teachers
	CreateTeacher(ctx context.Context, teacher *entity.Teacher) error
	UpdateTeacher(ctx context.Context, teacher *entity.Teacher) error
	GetTeachersBySchool(ctx context.Context, schoolID int64) ([]*entity.Teacher, error)
	GetTeacherByUserID(ctx context.Context, userID int64) (*entity.Teacher, error)
	GetTeacherByID(ctx context.Context, id int64) (*entity.Teacher, error)

	// Students
	CreateStudent(ctx context.Context, student *entity.Student) error
	UpdateStudent(ctx context.Context, student *entity.Student) error
	DeleteStudent(ctx context.Context, id int64) error
	GetStudentByID(ctx context.Context, id int64) (*entity.Student, error)
	GetStudentsBySchool(ctx context.Context, schoolID int64) ([]*entity.Student, error)
	GetStudentByUserID(ctx context.Context, userID int64) (*entity.Student, error)
	GetAllStudents(ctx context.Context, schoolID int64, limit, offset int, filters map[string]interface{}) ([]*entity.Student, int64, error)
	GetStudentsBySection(ctx context.Context, sectionID int64) ([]*entity.Student, error)

	// Enrollment
	EnrollStudent(ctx context.Context, enrollment *entity.StudentSection) error
	GetStudentSections(ctx context.Context, studentID int64) ([]*entity.Section, error)

	// Parents
	CreateParent(ctx context.Context, parent *entity.Parent) error
	UpdateParent(ctx context.Context, parent *entity.Parent) error
	GetParentByID(ctx context.Context, id int64) (*entity.Parent, error)
	LinkParentToStudent(ctx context.Context, link *entity.StudentParent) error
	GetStudentParents(ctx context.Context, studentID int64) ([]*entity.Parent, error)
}

type peopleUsecase struct {
	repo     postgresql.PeopleRepository
	validate *validator.Validate
	log      *logrus.Logger
}

func NewPeopleUsecase(repo postgresql.PeopleRepository, validate *validator.Validate, log *logrus.Logger) PeopleUsecase {
	return &peopleUsecase{
		repo:     repo,
		validate: validate,
		log:      log,
	}
}

func (uc *peopleUsecase) CreateTeacher(ctx context.Context, teacher *entity.Teacher) error {
	if err := uc.validate.Struct(teacher); err != nil {
		return err
	}
	return uc.repo.CreateTeacher(ctx, teacher)
}

func (uc *peopleUsecase) UpdateTeacher(ctx context.Context, teacher *entity.Teacher) error {
	if err := uc.validate.Struct(teacher); err != nil {
		return err
	}
	return uc.repo.UpdateTeacher(ctx, teacher)
}

func (uc *peopleUsecase) GetTeachersBySchool(ctx context.Context, schoolID int64) ([]*entity.Teacher, error) {
	return uc.repo.FindTeachersBySchool(ctx, schoolID)
}

func (uc *peopleUsecase) GetTeacherByUserID(ctx context.Context, userID int64) (*entity.Teacher, error) {
	return uc.repo.FindTeacherByUserID(ctx, userID)
}

func (uc *peopleUsecase) GetTeacherByID(ctx context.Context, id int64) (*entity.Teacher, error) {
	return uc.repo.FindTeacherByID(ctx, id)
}

func (uc *peopleUsecase) CreateStudent(ctx context.Context, student *entity.Student) error {
	if err := uc.validate.Struct(student); err != nil {
		return err
	}
	return uc.repo.CreateStudent(ctx, student)
}

func (uc *peopleUsecase) UpdateStudent(ctx context.Context, student *entity.Student) error {
	if err := uc.validate.Struct(student); err != nil {
		return err
	}
	return uc.repo.UpdateStudent(ctx, student)
}

func (uc *peopleUsecase) DeleteStudent(ctx context.Context, id int64) error {
	return uc.repo.DeleteStudent(ctx, id)
}

func (uc *peopleUsecase) GetStudentByID(ctx context.Context, id int64) (*entity.Student, error) {
	return uc.repo.FindStudentByID(ctx, id)
}

func (uc *peopleUsecase) GetStudentsBySchool(ctx context.Context, schoolID int64) ([]*entity.Student, error) {
	return uc.repo.FindStudentsBySchool(ctx, schoolID)
}

func (uc *peopleUsecase) GetStudentByUserID(ctx context.Context, userID int64) (*entity.Student, error) {
	return uc.repo.FindStudentByUserID(ctx, userID)
}

func (uc *peopleUsecase) GetAllStudents(ctx context.Context, schoolID int64, limit, offset int, filters map[string]interface{}) ([]*entity.Student, int64, error) {
	return uc.repo.GetAllStudents(ctx, schoolID, limit, offset, filters)
}

func (uc *peopleUsecase) GetStudentsBySection(ctx context.Context, sectionID int64) ([]*entity.Student, error) {
	return uc.repo.GetStudentsBySection(ctx, sectionID)
}

func (uc *peopleUsecase) EnrollStudent(ctx context.Context, enrollment *entity.StudentSection) error {
	return uc.repo.EnrollStudentInSection(ctx, enrollment)
}

func (uc *peopleUsecase) GetStudentSections(ctx context.Context, studentID int64) ([]*entity.Section, error) {
	return uc.repo.GetStudentSections(ctx, studentID)
}

func (uc *peopleUsecase) CreateParent(ctx context.Context, parent *entity.Parent) error {
	if err := uc.validate.Struct(parent); err != nil {
		return err
	}
	return uc.repo.CreateParent(ctx, parent)
}

func (uc *peopleUsecase) UpdateParent(ctx context.Context, parent *entity.Parent) error {
	if err := uc.validate.Struct(parent); err != nil {
		return err
	}
	return uc.repo.UpdateParent(ctx, parent)
}

func (uc *peopleUsecase) GetParentByID(ctx context.Context, id int64) (*entity.Parent, error) {
	return uc.repo.FindParentByID(ctx, id)
}

func (uc *peopleUsecase) LinkParentToStudent(ctx context.Context, link *entity.StudentParent) error {
	return uc.repo.LinkParentToStudent(ctx, link)
}

func (uc *peopleUsecase) GetStudentParents(ctx context.Context, studentID int64) ([]*entity.Parent, error) {
	return uc.repo.GetStudentParents(ctx, studentID)
}

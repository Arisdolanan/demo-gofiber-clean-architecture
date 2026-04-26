package postgresql

import (
	"context"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/jmoiron/sqlx"
)

type PeopleRepository interface {
	// Teachers
	CreateTeacher(ctx context.Context, teacher *entity.Teacher) error
	UpdateTeacher(ctx context.Context, teacher *entity.Teacher) error
	FindTeachersBySchool(ctx context.Context, schoolID int64) ([]*entity.Teacher, error)
	FindTeacherByUserID(ctx context.Context, userID int64) (*entity.Teacher, error)
	FindTeacherByID(ctx context.Context, id int64) (*entity.Teacher, error)

	// Students
	CreateStudent(ctx context.Context, student *entity.Student) error
	UpdateStudent(ctx context.Context, student *entity.Student) error
	DeleteStudent(ctx context.Context, id int64) error
	FindStudentByID(ctx context.Context, id int64) (*entity.Student, error)
	FindStudentsBySchool(ctx context.Context, schoolID int64) ([]*entity.Student, error)
	FindStudentByUserID(ctx context.Context, userID int64) (*entity.Student, error)
	GetAllStudents(ctx context.Context, schoolID int64, limit, offset int, filters map[string]interface{}) ([]*entity.Student, int64, error)
	GetStudentsBySection(ctx context.Context, sectionID int64) ([]*entity.Student, error)

	// Student-Section enrollment
	EnrollStudentInSection(ctx context.Context, enrollment *entity.StudentSection) error
	UpdateStudentSection(ctx context.Context, enrollment *entity.StudentSection) error
	GetStudentSections(ctx context.Context, studentID int64) ([]*entity.Section, error)

	// Parents
	CreateParent(ctx context.Context, parent *entity.Parent) error
	UpdateParent(ctx context.Context, parent *entity.Parent) error
	FindParentByID(ctx context.Context, id int64) (*entity.Parent, error)
	LinkParentToStudent(ctx context.Context, link *entity.StudentParent) error
	IsStudentParent(ctx context.Context, studentID, parentUserID int64) (bool, error)
	GetStudentParents(ctx context.Context, studentID int64) ([]*entity.Parent, error)
}

type peopleRepository struct {
	teacherRepo        *BaseRepository[entity.Teacher]
	studentRepo        *BaseRepository[entity.Student]
	parentRepo         *BaseRepository[entity.Parent]
	studentSectionRepo *BaseRepository[entity.StudentSection]
	studentParentRepo  *BaseRepository[entity.StudentParent]
	db                 *sqlx.DB
}

func NewPeopleRepository(db *sqlx.DB) PeopleRepository {
	return &peopleRepository{
		teacherRepo:        NewBaseRepository[entity.Teacher](db, "teachers"),
		studentRepo:        NewBaseRepository[entity.Student](db, "students"),
		parentRepo:         NewBaseRepository[entity.Parent](db, "parents"),
		studentSectionRepo: NewBaseRepository[entity.StudentSection](db, "student_sections"),
		studentParentRepo:  NewBaseRepository[entity.StudentParent](db, "student_parents"),
		db:                 db,
	}
}

func (r *peopleRepository) CreateTeacher(ctx context.Context, teacher *entity.Teacher) error {
	return r.teacherRepo.Create(ctx, teacher)
}

func (r *peopleRepository) UpdateTeacher(ctx context.Context, teacher *entity.Teacher) error {
	return r.teacherRepo.Update(ctx, teacher, "id = $1 AND deleted_at IS NULL", teacher.ID)
}

func (r *peopleRepository) FindTeachersBySchool(ctx context.Context, schoolID int64) ([]*entity.Teacher, error) {
	return r.teacherRepo.FindAll(ctx, "school_id = $1 AND deleted_at IS NULL", schoolID)
}

func (r *peopleRepository) FindTeacherByUserID(ctx context.Context, userID int64) (*entity.Teacher, error) {
	return r.teacherRepo.FindOne(ctx, "user_id = $1 AND deleted_at IS NULL", userID)
}

func (r *peopleRepository) FindTeacherByID(ctx context.Context, id int64) (*entity.Teacher, error) {
	return r.teacherRepo.FindByID(ctx, id)
}

func (r *peopleRepository) CreateStudent(ctx context.Context, student *entity.Student) error {
	return r.studentRepo.Create(ctx, student)
}

func (r *peopleRepository) UpdateStudent(ctx context.Context, student *entity.Student) error {
	return r.studentRepo.Update(ctx, student, "id = $1 AND deleted_at IS NULL", student.ID)
}

func (r *peopleRepository) DeleteStudent(ctx context.Context, id int64) error {
	return r.studentRepo.SoftDelete(ctx, "id = $1 AND deleted_at IS NULL", id)
}

func (r *peopleRepository) FindStudentByID(ctx context.Context, id int64) (*entity.Student, error) {
	return r.studentRepo.FindByID(ctx, id)
}

func (r *peopleRepository) FindStudentsBySchool(ctx context.Context, schoolID int64) ([]*entity.Student, error) {
	return r.studentRepo.FindAll(ctx, "school_id = $1 AND deleted_at IS NULL", schoolID)
}

func (r *peopleRepository) FindStudentByUserID(ctx context.Context, userID int64) (*entity.Student, error) {
	return r.studentRepo.FindOne(ctx, "user_id = $1 AND deleted_at IS NULL", userID)
}

func (r *peopleRepository) GetAllStudents(ctx context.Context, schoolID int64, limit, offset int, filters map[string]interface{}) ([]*entity.Student, int64, error) {
	where := "school_id = $1 AND deleted_at IS NULL"
	args := []interface{}{schoolID}

	// Dynamic filters can be added here
	count, err := r.studentRepo.Count(ctx, where, args...)
	if err != nil {
		return nil, 0, err
	}

	students, err := r.studentRepo.FindAllWithPagination(ctx, limit, offset, where, args...)
	return students, count, err
}

func (r *peopleRepository) GetStudentsBySection(ctx context.Context, sectionID int64) ([]*entity.Student, error) {
	var students []*entity.Student
	query := `
		SELECT s.* FROM students s
		JOIN student_sections ss ON s.id = ss.student_id
		WHERE ss.section_id = $1 AND s.deleted_at IS NULL AND ss.deleted_at IS NULL
	`
	err := r.db.SelectContext(ctx, &students, query, sectionID)
	return students, err
}

func (r *peopleRepository) EnrollStudentInSection(ctx context.Context, enrollment *entity.StudentSection) error {
	return r.studentSectionRepo.Create(ctx, enrollment)
}

func (r *peopleRepository) UpdateStudentSection(ctx context.Context, enrollment *entity.StudentSection) error {
	return r.studentSectionRepo.Update(ctx, enrollment, "id = $1 AND deleted_at IS NULL", enrollment.ID)
}

func (r *peopleRepository) GetStudentSections(ctx context.Context, studentID int64) ([]*entity.Section, error) {
	var sections []*entity.Section
	query := `
		SELECT s.* FROM sections s
		JOIN student_sections ss ON s.id = ss.section_id
		WHERE ss.student_id = $1 AND ss.deleted_at IS NULL
	`
	err := r.db.SelectContext(ctx, &sections, query, studentID)
	return sections, err
}

func (r *peopleRepository) CreateParent(ctx context.Context, parent *entity.Parent) error {
	return r.parentRepo.Create(ctx, parent)
}

func (r *peopleRepository) UpdateParent(ctx context.Context, parent *entity.Parent) error {
	return r.parentRepo.Update(ctx, parent, "id = $1 AND deleted_at IS NULL", parent.ID)
}

func (r *peopleRepository) FindParentByID(ctx context.Context, id int64) (*entity.Parent, error) {
	return r.parentRepo.FindByID(ctx, id)
}

func (r *peopleRepository) LinkParentToStudent(ctx context.Context, link *entity.StudentParent) error {
	return r.studentParentRepo.Create(ctx, link)
}

func (r *peopleRepository) IsStudentParent(ctx context.Context, studentID, parentUserID int64) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM student_parents sp
		JOIN parents p ON sp.parent_id = p.id
		WHERE sp.student_id = $1 AND p.user_id = $2 AND sp.deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &count, query, studentID, parentUserID)
	return count > 0, err
}

func (r *peopleRepository) GetStudentParents(ctx context.Context, studentID int64) ([]*entity.Parent, error) {
	var parents []*entity.Parent
	query := `
		SELECT p.* FROM parents p
		JOIN student_parents sp ON p.id = sp.parent_id
		WHERE sp.student_id = $1 AND sp.deleted_at IS NULL AND p.deleted_at IS NULL
	`
	err := r.db.SelectContext(ctx, &parents, query, studentID)
	return parents, err
}

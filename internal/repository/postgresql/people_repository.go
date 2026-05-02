package postgresql

import (
	"context"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/jmoiron/sqlx"
)

type PeopleRepository interface {
	// Teachers
	CreateTeacher(ctx context.Context, teacher *entity.Teacher) error
	UpdateTeacher(ctx context.Context, schoolID int64, teacher *entity.Teacher) error
	FindTeachersBySchool(ctx context.Context, schoolID int64) ([]*entity.Teacher, error)
	FindTeacherByUserID(ctx context.Context, schoolID, userID int64) (*entity.Teacher, error)
	FindTeacherByID(ctx context.Context, schoolID, id int64) (*entity.Teacher, error)
	DeleteTeacher(ctx context.Context, schoolID, id int64) error

	// Students
	CreateStudent(ctx context.Context, student *entity.Student) error
	UpdateStudent(ctx context.Context, schoolID int64, student *entity.Student) error
	DeleteStudent(ctx context.Context, schoolID, id int64) error
	FindStudentByID(ctx context.Context, schoolID, id int64) (*entity.Student, error)
	FindStudentsBySchool(ctx context.Context, schoolID int64) ([]*entity.Student, error)
	FindStudentByUserID(ctx context.Context, schoolID, userID int64) (*entity.Student, error)
	GetAllStudents(ctx context.Context, schoolID int64, limit, offset int, filters map[string]interface{}) ([]*entity.Student, int64, error)
	GetStudentsBySection(ctx context.Context, schoolID, sectionID int64) ([]*entity.Student, error)

	// Student-Section enrollment
	EnrollStudentInSection(ctx context.Context, enrollment *entity.StudentSection) error
	UpdateStudentSection(ctx context.Context, enrollment *entity.StudentSection) error
	GetStudentSections(ctx context.Context, schoolID, studentID int64) ([]*entity.Section, error)

	// Parents
	CreateParent(ctx context.Context, parent *entity.Parent) error
	UpdateParent(ctx context.Context, schoolID int64, parent *entity.Parent) error
	DeleteParent(ctx context.Context, schoolID, id int64) error
	FindParentByID(ctx context.Context, schoolID, id int64) (*entity.Parent, error)
	FindParentsBySchool(ctx context.Context, schoolID int64) ([]*entity.Parent, error)
	LinkParentToStudent(ctx context.Context, link *entity.StudentParent) error
	DeleteStudentParentsLinks(ctx context.Context, schoolID int64, studentID int64) error
	IsStudentParent(ctx context.Context, studentID, parentUserID int64) (bool, error)
	GetStudentParents(ctx context.Context, schoolID, studentID int64) ([]*entity.Parent, error)
	SearchParents(ctx context.Context, schoolID int64, query string) ([]*entity.Parent, error)

	// Staff
	CreateStaff(ctx context.Context, staff *entity.Staff) error
	UpdateStaff(ctx context.Context, schoolID int64, staff *entity.Staff) error
	FindStaffBySchool(ctx context.Context, schoolID int64) ([]*entity.Staff, error)
	FindStaffByUserID(ctx context.Context, schoolID, userID int64) (*entity.Staff, error)
	FindStaffByID(ctx context.Context, schoolID, id int64) (*entity.Staff, error)
	DeleteStaff(ctx context.Context, schoolID, id int64) error
}

type peopleRepository struct {
	teacherRepo        *BaseRepository[entity.Teacher]
	studentRepo        *BaseRepository[entity.Student]
	parentRepo         *BaseRepository[entity.Parent]
	staffRepo          *BaseRepository[entity.Staff]
	studentSectionRepo *BaseRepository[entity.StudentSection]
	studentParentRepo  *BaseRepository[entity.StudentParent]
	db                 *sqlx.DB
}

func NewPeopleRepository(db *sqlx.DB) PeopleRepository {
	return &peopleRepository{
		teacherRepo:        NewBaseRepository[entity.Teacher](db, "teachers"),
		studentRepo:        NewBaseRepository[entity.Student](db, "students"),
		parentRepo:         NewBaseRepository[entity.Parent](db, "parents"),
		staffRepo:          NewBaseRepository[entity.Staff](db, "staff"),
		studentSectionRepo: NewBaseRepository[entity.StudentSection](db, "student_sections"),
		studentParentRepo:  NewBaseRepository[entity.StudentParent](db, "student_parents"),
		db:                 db,
	}
}

func (r *peopleRepository) CreateTeacher(ctx context.Context, teacher *entity.Teacher) error {
	return r.teacherRepo.Create(ctx, teacher)
}

func (r *peopleRepository) UpdateTeacher(ctx context.Context, schoolID int64, teacher *entity.Teacher) error {
	return r.teacherRepo.Update(ctx, teacher, "id = $1 AND school_id = $2 AND deleted_at IS NULL", teacher.ID, schoolID)
}

func (r *peopleRepository) FindTeachersBySchool(ctx context.Context, schoolID int64) ([]*entity.Teacher, error) {
	return r.teacherRepo.FindAll(ctx, "school_id = $1 AND deleted_at IS NULL", schoolID)
}

func (r *peopleRepository) FindTeacherByUserID(ctx context.Context, schoolID, userID int64) (*entity.Teacher, error) {
	return r.teacherRepo.FindOne(ctx, "user_id = $1 AND school_id = $2 AND deleted_at IS NULL", userID, schoolID)
}

func (r *peopleRepository) FindTeacherByID(ctx context.Context, schoolID, id int64) (*entity.Teacher, error) {
	return r.teacherRepo.FindOne(ctx, "id = $1 AND school_id = $2 AND deleted_at IS NULL", id, schoolID)
}

func (r *peopleRepository) DeleteTeacher(ctx context.Context, schoolID, id int64) error {
	return r.teacherRepo.SoftDelete(ctx, "id = $1 AND school_id = $2 AND deleted_at IS NULL", id, schoolID)
}

func (r *peopleRepository) CreateStudent(ctx context.Context, student *entity.Student) error {
	return r.studentRepo.Create(ctx, student)
}

func (r *peopleRepository) UpdateStudent(ctx context.Context, schoolID int64, student *entity.Student) error {
	return r.studentRepo.Update(ctx, student, "id = $1::BIGINT AND school_id = $2::BIGINT AND deleted_at IS NULL", student.ID, schoolID)
}

func (r *peopleRepository) DeleteStudent(ctx context.Context, schoolID, id int64) error {
	return r.studentRepo.SoftDelete(ctx, "id = $1 AND school_id = $2 AND deleted_at IS NULL", id, schoolID)
}

func (r *peopleRepository) FindStudentByID(ctx context.Context, schoolID, id int64) (*entity.Student, error) {
	student, err := r.studentRepo.FindOne(ctx, "id = $1 AND school_id = $2 AND deleted_at IS NULL", id, schoolID)
	if err != nil {
		return nil, err
	}

	// Fetch parents for this student
	query := `
		SELECT 
			p.id as parent_id, 
			p.full_name, 
			p.phone, 
			p.email, 
			p.address, 
			p.occupation,
			sp.relationship, 
			sp.is_primary
		FROM parents p
		JOIN student_parents sp ON p.id = sp.parent_id
		WHERE sp.student_id = $1 AND sp.school_id = $2 AND p.deleted_at IS NULL
	`
	var parents []entity.StudentParentRequest
	err = r.db.SelectContext(ctx, &parents, query, id, schoolID)
	if err == nil {
		student.Parents = parents
	}

	return student, nil
}

func (r *peopleRepository) FindStudentsBySchool(ctx context.Context, schoolID int64) ([]*entity.Student, error) {
	return r.studentRepo.FindAll(ctx, "school_id = $1 AND deleted_at IS NULL", schoolID)
}

func (r *peopleRepository) FindStudentByUserID(ctx context.Context, schoolID, userID int64) (*entity.Student, error) {
	return r.studentRepo.FindOne(ctx, "user_id = $1 AND school_id = $2 AND deleted_at IS NULL", userID, schoolID)
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
	if err != nil {
		return nil, 0, err
	}

	// Hydrate parents for each student in the list
	for _, s := range students {
		query := `
			SELECT 
				p.id as parent_id, 
				p.full_name, 
				p.phone, 
				p.email, 
				p.address, 
				p.occupation,
				sp.relationship, 
				sp.is_primary
			FROM parents p
			JOIN student_parents sp ON p.id = sp.parent_id
			WHERE sp.student_id = $1 AND sp.school_id = $2 AND p.deleted_at IS NULL
		`
		var parents []entity.StudentParentRequest
		if err := r.db.SelectContext(ctx, &parents, query, s.ID, schoolID); err == nil {
			s.Parents = parents
		}
	}

	return students, count, nil
}

func (r *peopleRepository) GetStudentsBySection(ctx context.Context, schoolID, sectionID int64) ([]*entity.Student, error) {
	var students []*entity.Student
	query := `
		SELECT s.* FROM students s
		JOIN student_sections ss ON s.id = ss.student_id
		WHERE ss.section_id = $1 AND s.school_id = $2 AND s.deleted_at IS NULL AND ss.deleted_at IS NULL
	`
	err := r.db.SelectContext(ctx, &students, query, sectionID, schoolID)
	return students, err
}

func (r *peopleRepository) EnrollStudentInSection(ctx context.Context, enrollment *entity.StudentSection) error {
	return r.studentSectionRepo.Create(ctx, enrollment)
}

func (r *peopleRepository) UpdateStudentSection(ctx context.Context, enrollment *entity.StudentSection) error {
	return r.studentSectionRepo.Update(ctx, enrollment, "id = $1 AND deleted_at IS NULL", enrollment.ID)
}

func (r *peopleRepository) GetStudentSections(ctx context.Context, schoolID, studentID int64) ([]*entity.Section, error) {
	var sections []*entity.Section
	query := `
		SELECT s.* FROM sections s
		JOIN student_sections ss ON s.id = ss.section_id
		JOIN students st ON ss.student_id = st.id
		WHERE ss.student_id = $1 AND st.school_id = $2 AND ss.deleted_at IS NULL
	`
	err := r.db.SelectContext(ctx, &sections, query, studentID, schoolID)
	return sections, err
}

func (r *peopleRepository) CreateParent(ctx context.Context, parent *entity.Parent) error {
	err := r.parentRepo.Create(ctx, parent)
	if err != nil {
		return err
	}

	// Link children if any
	if len(parent.Children) > 0 {
		for _, child := range parent.Children {
			link := &entity.StudentParent{
				SchoolID:     parent.SchoolID,
				StudentID:    child.ID,
				ParentID:     parent.ID,
				Relationship: child.Relationship,
				IsPrimary:    child.IsPrimary,
			}
			err = r.studentParentRepo.Create(ctx, link)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *peopleRepository) UpdateParent(ctx context.Context, schoolID int64, parent *entity.Parent) error {
	err := r.parentRepo.Update(ctx, parent, "id = $1::BIGINT AND school_id = $2::BIGINT AND deleted_at IS NULL", parent.ID, schoolID)
	if err != nil {
		return err
	}

	// Update children - remove existing and add new
	if len(parent.Children) > 0 {
		// Remove existing links
		_, err = r.db.ExecContext(ctx, "DELETE FROM student_parents WHERE parent_id = $1 AND school_id = $2", parent.ID, schoolID)
		if err != nil {
			return err
		}

		// Add new links
		for _, child := range parent.Children {
			link := &entity.StudentParent{
				SchoolID:     schoolID,
				StudentID:    child.ID,
				ParentID:     parent.ID,
				Relationship: child.Relationship,
				IsPrimary:    child.IsPrimary,
			}
			err = r.studentParentRepo.Create(ctx, link)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *peopleRepository) FindParentByID(ctx context.Context, schoolID, id int64) (*entity.Parent, error) {
	parent, err := r.parentRepo.FindOne(ctx, "id = $1 AND school_id = $2 AND deleted_at IS NULL", id, schoolID)
	if err != nil {
		return nil, err
	}

	// Get children
	children, err := r.getParentChildren(ctx, id)
	if err == nil {
		parent.Children = children
	}

	return parent, nil
}

func (r *peopleRepository) FindParentsBySchool(ctx context.Context, schoolID int64) ([]*entity.Parent, error) {
	var parents []*entity.Parent
	query := `SELECT * FROM parents WHERE school_id = $1 AND deleted_at IS NULL ORDER BY full_name`
	err := r.db.SelectContext(ctx, &parents, query, schoolID)
	if err != nil {
		return nil, err
	}

	// Get children for each parent
	for i := range parents {
		children, err := r.getParentChildren(ctx, parents[i].ID)
		if err == nil {
			parents[i].Children = children
		}
	}

	return parents, nil
}

func (r *peopleRepository) getParentChildren(ctx context.Context, parentID int64) ([]entity.ParentChild, error) {
	var children []entity.ParentChild
	query := `
		SELECT s.id, s.full_name, s.student_number, sp.relationship, sp.is_primary
		FROM student_parents sp
		JOIN students s ON sp.student_id = s.id
		WHERE sp.parent_id = $1 AND sp.deleted_at IS NULL AND s.deleted_at IS NULL
		ORDER BY sp.is_primary DESC, s.full_name
	`
	err := r.db.SelectContext(ctx, &children, query, parentID)
	return children, err
}

func (r *peopleRepository) DeleteParent(ctx context.Context, schoolID, id int64) error {
	return r.parentRepo.SoftDelete(ctx, "id = $1 AND school_id = $2 AND deleted_at IS NULL", id, schoolID)
}

func (r *peopleRepository) LinkParentToStudent(ctx context.Context, link *entity.StudentParent) error {
	return r.studentParentRepo.Create(ctx, link)
}

func (r *peopleRepository) DeleteStudentParentsLinks(ctx context.Context, schoolID int64, studentID int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM student_parents WHERE student_id = $1 AND school_id = $2", studentID, schoolID)
	return err
}

func (r *peopleRepository) IsStudentParent(ctx context.Context, studentID, parentUserID int64) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM student_parents sp
		JOIN parents p ON sp.parent_id = p.id
		WHERE sp.student_id = $1 AND p.user_id = $2 AND sp.deleted_at IS NULL AND sp.school_id = p.school_id
	`
	err := r.db.GetContext(ctx, &count, query, studentID, parentUserID)
	return count > 0, err
}

func (r *peopleRepository) GetStudentParents(ctx context.Context, schoolID, studentID int64) ([]*entity.Parent, error) {
	var parents []*entity.Parent
	query := `
		SELECT p.* FROM parents p
		JOIN student_parents sp ON p.id = sp.parent_id
		JOIN students s ON sp.student_id = s.id
		WHERE sp.student_id = $1 AND s.school_id = $2 AND sp.school_id = $2 AND sp.deleted_at IS NULL AND p.deleted_at IS NULL
	`
	err := r.db.SelectContext(ctx, &parents, query, studentID, schoolID)
	return parents, err
}

func (r *peopleRepository) SearchParents(ctx context.Context, schoolID int64, query string) ([]*entity.Parent, error) {
	var parents []*entity.Parent
	sqlQuery := `
		SELECT * FROM parents 
		WHERE school_id = $1 AND deleted_at IS NULL 
		AND (full_name ILIKE $2 OR phone ILIKE $2 OR email ILIKE $2)
		LIMIT 10
	`
	err := r.db.SelectContext(ctx, &parents, sqlQuery, schoolID, "%"+query+"%")
	return parents, err
}

// Staff methods
func (r *peopleRepository) CreateStaff(ctx context.Context, staff *entity.Staff) error {
	return r.staffRepo.Create(ctx, staff)
}

func (r *peopleRepository) UpdateStaff(ctx context.Context, schoolID int64, staff *entity.Staff) error {
	return r.staffRepo.Update(ctx, staff, "id = $1 AND school_id = $2 AND deleted_at IS NULL", staff.ID, schoolID)
}

func (r *peopleRepository) FindStaffBySchool(ctx context.Context, schoolID int64) ([]*entity.Staff, error) {
	var staff []*entity.Staff
	query := `SELECT * FROM staff WHERE school_id = $1 AND deleted_at IS NULL ORDER BY full_name`
	err := r.db.SelectContext(ctx, &staff, query, schoolID)
	return staff, err
}

func (r *peopleRepository) FindStaffByUserID(ctx context.Context, schoolID, userID int64) (*entity.Staff, error) {
	return r.staffRepo.FindOne(ctx, "user_id = $1 AND school_id = $2 AND deleted_at IS NULL", userID, schoolID)
}

func (r *peopleRepository) FindStaffByID(ctx context.Context, schoolID, id int64) (*entity.Staff, error) {
	return r.staffRepo.FindOne(ctx, "id = $1 AND school_id = $2 AND deleted_at IS NULL", id, schoolID)
}

func (r *peopleRepository) DeleteStaff(ctx context.Context, schoolID, id int64) error {
	return r.staffRepo.SoftDelete(ctx, "id = $1 AND school_id = $2 AND deleted_at IS NULL", id, schoolID)
}

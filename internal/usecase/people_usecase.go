package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/repository/postgresql"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type PeopleUsecase interface {
	// Teachers
	CreateTeacher(ctx context.Context, schoolID int64, teacher *entity.Teacher) (*entity.EmployeeUserResponse, error)
	UpdateTeacher(ctx context.Context, schoolID int64, teacher *entity.Teacher) error
	GetTeachersBySchool(ctx context.Context, schoolID int64) ([]*entity.Teacher, error)
	GetTeacherByUserID(ctx context.Context, schoolID, userID int64) (*entity.Teacher, error)
	GetTeacherByID(ctx context.Context, schoolID, id int64) (*entity.Teacher, error)
	DeleteTeacher(ctx context.Context, schoolID, id int64) error

	// Students
	CreateStudent(ctx context.Context, schoolID int64, student *entity.Student) (*entity.EnrollmentCredentials, error)
	UpdateStudent(ctx context.Context, schoolID int64, student *entity.Student) error
	DeleteStudent(ctx context.Context, schoolID, id int64) error
	GetStudentByID(ctx context.Context, schoolID, id int64) (*entity.Student, error)
	GetStudentsBySchool(ctx context.Context, schoolID int64) ([]*entity.Student, error)
	GetStudentByUserID(ctx context.Context, schoolID, userID int64) (*entity.Student, error)
	GetAllStudents(ctx context.Context, schoolID int64, limit, offset int, filters map[string]interface{}) ([]*entity.Student, int64, error)
	GetStudentsBySection(ctx context.Context, schoolID, sectionID int64) ([]*entity.Student, error)

	// Enrollment
	EnrollStudent(ctx context.Context, schoolID int64, enrollment *entity.StudentSection) error
	GetStudentSections(ctx context.Context, schoolID, studentID int64) ([]*entity.Section, error)

	// Parents
	CreateParent(ctx context.Context, schoolID int64, parent *entity.Parent) (*entity.EmployeeUserResponse, error)
	UpdateParent(ctx context.Context, schoolID int64, parent *entity.Parent) error
	DeleteParent(ctx context.Context, schoolID, id int64) error
	GetParentByID(ctx context.Context, schoolID, id int64) (*entity.Parent, error)
	GetParentsBySchool(ctx context.Context, schoolID int64) ([]*entity.Parent, error)
	LinkParentToStudent(ctx context.Context, schoolID int64, link *entity.StudentParent) error
	GetStudentParents(ctx context.Context, schoolID, studentID int64) ([]*entity.Parent, error)
	SearchParents(ctx context.Context, schoolID int64, query string) ([]*entity.Parent, error)

	// Staff
	CreateStaff(ctx context.Context, schoolID int64, staff *entity.Staff) (*entity.EmployeeUserResponse, error)
	UpdateStaff(ctx context.Context, schoolID int64, staff *entity.Staff) error
	GetStaffBySchool(ctx context.Context, schoolID int64) ([]*entity.Staff, error)
	GetStaffByUserID(ctx context.Context, schoolID, userID int64) (*entity.Staff, error)
	DeleteStaff(ctx context.Context, schoolID, id int64) error
}

type peopleUsecase struct {
	repo        postgresql.PeopleRepository
	userUsecase UserUseCase
	validate    *validator.Validate
	log         *logrus.Logger
}

func NewPeopleUsecase(repo postgresql.PeopleRepository, userUsecase UserUseCase, validate *validator.Validate, log *logrus.Logger) PeopleUsecase {
	return &peopleUsecase{
		repo:        repo,
		userUsecase: userUsecase,
		validate:    validate,
		log:         log,
	}
}


func (uc *peopleUsecase) CreateTeacher(ctx context.Context, schoolID int64, teacher *entity.Teacher) (*entity.EmployeeUserResponse, error) {
	teacher.SchoolID = schoolID
	if err := uc.validate.Struct(teacher); err != nil {
		return nil, err
	}

	userResp, err := uc.ensureUserForEmployee(ctx, schoolID, teacher.EmployeeNumber, entity.UserTeacher, teacher.Email, teacher.FullName, teacher.Phone)
	if err != nil {
		return nil, err
	}
	teacher.UserID = userResp.UserID

	if err := uc.repo.CreateTeacher(ctx, teacher); err != nil {
		return nil, err
	}

	userResp.EmployeeID = teacher.ID
	return userResp, nil
}

func (uc *peopleUsecase) UpdateTeacher(ctx context.Context, schoolID int64, teacher *entity.Teacher) error {
	teacher.SchoolID = schoolID
	if err := uc.validate.Struct(teacher); err != nil {
		return err
	}

	// Sync user data
	_, _ = uc.ensureUserForEmployee(ctx, schoolID, teacher.EmployeeNumber, entity.UserTeacher, teacher.Email, teacher.FullName, teacher.Phone)

	return uc.repo.UpdateTeacher(ctx, schoolID, teacher)
}

func (uc *peopleUsecase) GetTeachersBySchool(ctx context.Context, schoolID int64) ([]*entity.Teacher, error) {
	return uc.repo.FindTeachersBySchool(ctx, schoolID)
}

func (uc *peopleUsecase) GetTeacherByUserID(ctx context.Context, schoolID, userID int64) (*entity.Teacher, error) {
	return uc.repo.FindTeacherByUserID(ctx, schoolID, userID)
}

func (uc *peopleUsecase) GetTeacherByID(ctx context.Context, schoolID, id int64) (*entity.Teacher, error) {
	return uc.repo.FindTeacherByID(ctx, schoolID, id)
}

func (uc *peopleUsecase) DeleteTeacher(ctx context.Context, schoolID, id int64) error {
	// Find teacher to get UserID
	teacher, err := uc.repo.FindTeacherByID(ctx, schoolID, id)
	if err == nil && teacher != nil {
		// Deactivate associated user
		if err := uc.userUsecase.DeactivateUser(ctx, schoolID, teacher.UserID); err != nil {
			uc.log.Warnf("Failed to deactivate user %d for teacher %d: %v", teacher.UserID, id, err)
		}
	}
	return uc.repo.DeleteTeacher(ctx, schoolID, id)
}

func (uc *peopleUsecase) CreateStudent(ctx context.Context, schoolID int64, student *entity.Student) (*entity.EnrollmentCredentials, error) {
	student.SchoolID = schoolID
	if err := uc.validate.Struct(student); err != nil {
		return nil, err
	}

	credentials := &entity.EnrollmentCredentials{
		Parents: []*entity.EmployeeUserResponse{},
	}

	// Create user for student
	userResp, err := uc.ensureUserForStudent(ctx, schoolID, student)
	if err == nil && userResp != nil {
		student.UserID = &userResp.UserID
		credentials.Student = userResp
	}

	// Create student
	if err := uc.repo.CreateStudent(ctx, student); err != nil {
		return nil, err
	}

	// Link parents if provided
	if len(student.Parents) > 0 {
		for _, pReq := range student.Parents {
			var parentID int64 = pReq.ParentID

			// If no ParentID, create new parent
			if parentID == 0 {
				newParent := &entity.Parent{
					SchoolID:   schoolID,
					FullName:   pReq.FullName,
					Phone:      pReq.Phone,
					Email:      pReq.Email,
					Address:    pReq.Address,
					Occupation: pReq.Occupation,
				}

				// Create user for parent
				pUserResp, _ := uc.ensureUserForParent(ctx, schoolID, newParent)
				if pUserResp != nil {
					newParent.UserID = &pUserResp.UserID
					credentials.Parents = append(credentials.Parents, pUserResp)
				}

				if err := uc.repo.CreateParent(ctx, newParent); err != nil {
					continue
				}
				parentID = newParent.ID
			}

			// Link parent to student
			link := &entity.StudentParent{
				SchoolID:     schoolID,
				StudentID:    student.ID,
				ParentID:     parentID,
				Relationship: pReq.Relationship,
				IsPrimary:    pReq.IsPrimary,
			}
			if err := uc.repo.LinkParentToStudent(ctx, link); err != nil {
				uc.log.Errorf("Failed to link parent %d to student %d: %v", parentID, student.ID, err)
			}
		}
	}

	// Handle automatic enrollment if section_id is provided
	if student.SectionID != nil {
		secID := *student.SectionID
		if secID > 0 {
			var sessID int64 = 0
			if student.AcademicSessionID != nil {
				sessID = *student.AcademicSessionID
			}
			
			enrollment := &entity.StudentSection{
				SchoolID:          schoolID,
				StudentID:         student.ID,
				SectionID:         secID,
				AcademicSessionID: sessID,
				Status:            utils.StringPtr("active"),
			}
			if err := uc.repo.EnrollStudentInSection(ctx, enrollment); err != nil {
				uc.log.Errorf("Failed to enroll student %d in section %d: %v", student.ID, secID, err)
			}
		}
	}

	return credentials, nil
}

func (uc *peopleUsecase) UpdateStudent(ctx context.Context, schoolID int64, student *entity.Student) error {
	student.SchoolID = schoolID
	if err := uc.validate.Struct(student); err != nil {
		return err
	}
	if err := uc.repo.UpdateStudent(ctx, schoolID, student); err != nil {
		return err
	}

	// Sync user data
	_, _ = uc.ensureUserForStudent(ctx, schoolID, student)

	// Handle parent updates/linking from student form
	if student.Parents != nil {
		// First, delete existing links for this student to ensure a clean state
		if err := uc.repo.DeleteStudentParentsLinks(ctx, schoolID, student.ID); err != nil {
			uc.log.Errorf("Failed to clear existing parent links for student %d: %v", student.ID, err)
		}

		for _, pReq := range student.Parents {
			var parentID int64
			
			if pReq.ParentID > 0 {
				parentID = pReq.ParentID
				// Update existing parent info if provided
				parent := &entity.Parent{
					ID:         pReq.ParentID,
					SchoolID:   schoolID,
					FullName:   pReq.FullName,
					Phone:      pReq.Phone,
					Email:      pReq.Email,
					Address:    pReq.Address,
					Occupation: pReq.Occupation,
				}
				if err := uc.repo.UpdateParent(ctx, schoolID, parent); err != nil {
					uc.log.Errorf("Failed to update parent %d: %v", pReq.ParentID, err)
				}
			} else {
				// Create new parent
				newParent := &entity.Parent{
					SchoolID:   schoolID,
					FullName:   pReq.FullName,
					Phone:      pReq.Phone,
					Email:      pReq.Email,
					Address:    pReq.Address,
					Occupation: pReq.Occupation,
				}
				if err := uc.repo.CreateParent(ctx, newParent); err == nil {
					parentID = newParent.ID
				} else {
					uc.log.Errorf("Failed to create new parent for student %d: %v", student.ID, err)
				}
			}

			// Link the parent to the student
			if parentID > 0 {
				link := &entity.StudentParent{
					SchoolID:     schoolID,
					StudentID:    student.ID,
					ParentID:     parentID,
					Relationship: pReq.Relationship,
					IsPrimary:    pReq.IsPrimary,
				}
				if err := uc.repo.LinkParentToStudent(ctx, link); err != nil {
					uc.log.Errorf("Failed to link parent %d to student %d: %v", parentID, student.ID, err)
				}
			}
		}
	}

	return nil
}

func (uc *peopleUsecase) DeleteStudent(ctx context.Context, schoolID, id int64) error {
	// Find student to get UserID
	student, err := uc.repo.FindStudentByID(ctx, schoolID, id)
	if err == nil && student != nil && student.UserID != nil {
		// Deactivate associated user
		if err := uc.userUsecase.DeactivateUser(ctx, schoolID, *student.UserID); err != nil {
			uc.log.Warnf("Failed to deactivate user %d for student %d: %v", *student.UserID, id, err)
		}
	}
	return uc.repo.DeleteStudent(ctx, schoolID, id)
}

func (uc *peopleUsecase) GetStudentByID(ctx context.Context, schoolID, id int64) (*entity.Student, error) {
	return uc.repo.FindStudentByID(ctx, schoolID, id)
}

func (uc *peopleUsecase) GetStudentsBySchool(ctx context.Context, schoolID int64) ([]*entity.Student, error) {
	return uc.repo.FindStudentsBySchool(ctx, schoolID)
}

func (uc *peopleUsecase) GetStudentByUserID(ctx context.Context, schoolID, userID int64) (*entity.Student, error) {
	return uc.repo.FindStudentByUserID(ctx, schoolID, userID)
}

func (uc *peopleUsecase) GetAllStudents(ctx context.Context, schoolID int64, limit, offset int, filters map[string]interface{}) ([]*entity.Student, int64, error) {
	return uc.repo.GetAllStudents(ctx, schoolID, limit, offset, filters)
}

func (uc *peopleUsecase) GetStudentsBySection(ctx context.Context, schoolID, sectionID int64) ([]*entity.Student, error) {
	return uc.repo.GetStudentsBySection(ctx, schoolID, sectionID)
}

func (uc *peopleUsecase) EnrollStudent(ctx context.Context, schoolID int64, enrollment *entity.StudentSection) error {
	// Optional: verify that the section belongs to the school
	return uc.repo.EnrollStudentInSection(ctx, enrollment)
}

func (uc *peopleUsecase) GetStudentSections(ctx context.Context, schoolID, studentID int64) ([]*entity.Section, error) {
	return uc.repo.GetStudentSections(ctx, schoolID, studentID)
}

func (uc *peopleUsecase) CreateParent(ctx context.Context, schoolID int64, parent *entity.Parent) (*entity.EmployeeUserResponse, error) {
	parent.SchoolID = schoolID
	// Skip validation for children field since it's handled separately
	if err := uc.validate.StructExcept(parent, "Children"); err != nil {
		return nil, err
	}

	// Create user for parent
	userResp, _ := uc.ensureUserForParent(ctx, schoolID, parent)
	if userResp != nil {
		parent.UserID = &userResp.UserID
	}

	if err := uc.repo.CreateParent(ctx, parent); err != nil {
		return nil, err
	}

	return userResp, nil
}

func (uc *peopleUsecase) UpdateParent(ctx context.Context, schoolID int64, parent *entity.Parent) error {
	parent.SchoolID = schoolID
	// Skip validation for children field since it's handled separately
	if err := uc.validate.StructExcept(parent, "Children"); err != nil {
		return err
	}
	if err := uc.repo.UpdateParent(ctx, schoolID, parent); err != nil {
		return err
	}

	// Sync user data
	_, _ = uc.ensureUserForParent(ctx, schoolID, parent)

	return nil
}

func (uc *peopleUsecase) GetParentByID(ctx context.Context, schoolID, id int64) (*entity.Parent, error) {
	return uc.repo.FindParentByID(ctx, schoolID, id)
}

func (uc *peopleUsecase) GetParentsBySchool(ctx context.Context, schoolID int64) ([]*entity.Parent, error) {
	return uc.repo.FindParentsBySchool(ctx, schoolID)
}

func (uc *peopleUsecase) DeleteParent(ctx context.Context, schoolID, id int64) error {
	// Find parent to get UserID
	parent, err := uc.repo.FindParentByID(ctx, schoolID, id)
	if err == nil && parent != nil && parent.UserID != nil {
		// Deactivate associated user
		if err := uc.userUsecase.DeactivateUser(ctx, schoolID, *parent.UserID); err != nil {
			uc.log.Warnf("Failed to deactivate user %d for parent %d: %v", *parent.UserID, id, err)
		}
	}
	return uc.repo.DeleteParent(ctx, schoolID, id)
}

func (uc *peopleUsecase) LinkParentToStudent(ctx context.Context, schoolID int64, link *entity.StudentParent) error {
	link.SchoolID = schoolID
	return uc.repo.LinkParentToStudent(ctx, link)
}

func (uc *peopleUsecase) GetStudentParents(ctx context.Context, schoolID, studentID int64) ([]*entity.Parent, error) {
	return uc.repo.GetStudentParents(ctx, schoolID, studentID)
}

func (uc *peopleUsecase) SearchParents(ctx context.Context, schoolID int64, query string) ([]*entity.Parent, error) {
	return uc.repo.SearchParents(ctx, schoolID, query)
}

func (uc *peopleUsecase) CreateStaff(ctx context.Context, schoolID int64, staff *entity.Staff) (*entity.EmployeeUserResponse, error) {
	staff.SchoolID = schoolID
	if err := uc.validate.Struct(staff); err != nil {
		return nil, err
	}

	userResp, err := uc.ensureUserForEmployee(ctx, schoolID, staff.EmployeeNumber, entity.UserStaff, staff.Email, staff.FullName, staff.Phone)
	if err != nil {
		return nil, err
	}
	staff.UserID = &userResp.UserID

	if err := uc.repo.CreateStaff(ctx, staff); err != nil {
		return nil, err
	}

	userResp.EmployeeID = staff.ID
	return userResp, nil
}

func (uc *peopleUsecase) UpdateStaff(ctx context.Context, schoolID int64, staff *entity.Staff) error {
	staff.SchoolID = schoolID
	if err := uc.validate.Struct(staff); err != nil {
		return err
	}

	// Sync user data
	_, _ = uc.ensureUserForEmployee(ctx, schoolID, staff.EmployeeNumber, entity.UserStaff, staff.Email, staff.FullName, staff.Phone)

	return uc.repo.UpdateStaff(ctx, schoolID, staff)
}

func (uc *peopleUsecase) GetStaffBySchool(ctx context.Context, schoolID int64) ([]*entity.Staff, error) {
	return uc.repo.FindStaffBySchool(ctx, schoolID)
}

func (uc *peopleUsecase) GetStaffByUserID(ctx context.Context, schoolID, userID int64) (*entity.Staff, error) {
	return uc.repo.FindStaffByUserID(ctx, schoolID, userID)
}

func (uc *peopleUsecase) DeleteStaff(ctx context.Context, schoolID, id int64) error {
	// Find staff to get UserID
	staff, err := uc.repo.FindStaffByID(ctx, schoolID, id)
	if err == nil && staff != nil && staff.UserID != nil {
		// Deactivate associated user
		if err := uc.userUsecase.DeactivateUser(ctx, schoolID, *staff.UserID); err != nil {
			uc.log.Warnf("Failed to deactivate user %d for staff %d: %v", *staff.UserID, id, err)
		}
	}
	return uc.repo.DeleteStaff(ctx, schoolID, id)
}
func (uc *peopleUsecase) ensureUserForEmployee(ctx context.Context, schoolID int64, employeeNumber string, userType entity.UserType, email *string, fullName *string, phone *string) (*entity.EmployeeUserResponse, error) {
	// 1. Try to find existing user by username (employeeNumber)
	existingUser, err := uc.userUsecase.GetUserByUsername(employeeNumber)
	if err == nil && existingUser != nil {
		// User exists, update FullName, Email, and Phone if they changed
		updated := false
		if fullName != nil && (existingUser.FullName == nil || *existingUser.FullName != *fullName) {
			existingUser.FullName = fullName
			updated = true
		}
		if email != nil && existingUser.Email != *email {
			existingUser.Email = *email
			updated = true
		}
		if phone != nil && (existingUser.PhoneNumber == nil || *existingUser.PhoneNumber != *phone) {
			existingUser.PhoneNumber = phone
			updated = true
		}

		if updated {
			_ = uc.userUsecase.UpdateUser(ctx, schoolID, existingUser)
		}

		username := ""
		if existingUser.Username != nil {
			username = *existingUser.Username
		}
		return &entity.EmployeeUserResponse{
			UserID:   existingUser.ID,
			Username: username,
			UserType: string(userType),
		}, nil
	}

	// 2. User doesn't exist, create one
	password := utils.GenerateRandomPassword(10)
	hashedPassword, _ := utils.HashPassword(password)
	
	user := &entity.User{
		Username:           &employeeNumber,
		Password:           hashedPassword,
		UserType:           userType,
		IsActive:           true,
		MustChangePassword: true,
		SchoolID:           []int64{schoolID},
	}

	if fullName != nil {
		user.FullName = fullName
	}
	if email != nil {
		user.Email = *email
	}
	if phone != nil {
		user.PhoneNumber = phone
	}

	if err := uc.userUsecase.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	username := ""
	if user.Username != nil {
		username = *user.Username
	}

	return &entity.EmployeeUserResponse{
		UserID:   user.ID,
		Username: username,
		Password: password,
		UserType: string(userType),
	}, nil
}

func (uc *peopleUsecase) ensureUserForStudent(ctx context.Context, schoolID int64, student *entity.Student) (*entity.EmployeeUserResponse, error) {
	if student.StudentNumber == "" {
		return nil, errors.New("student number is required for user creation")
	}

	// 1. Try to find existing user by username (student number)
	existingUser, err := uc.userUsecase.GetUserByUsername(student.StudentNumber)
	if err == nil && existingUser != nil {
		// Update profile
		updated := false
		if student.FullName != nil && (existingUser.FullName == nil || *existingUser.FullName != *student.FullName) {
			existingUser.FullName = student.FullName
			updated = true
		}
		if student.Email != nil && existingUser.Email != *student.Email {
			existingUser.Email = *student.Email
			updated = true
		}
		if student.Phone != nil && (existingUser.PhoneNumber == nil || *existingUser.PhoneNumber != *student.Phone) {
			existingUser.PhoneNumber = student.Phone
			updated = true
		}

		if updated {
			_ = uc.userUsecase.UpdateUser(ctx, schoolID, existingUser)
		}

		username := ""
		if existingUser.Username != nil {
			username = *existingUser.Username
		}
		return &entity.EmployeeUserResponse{
			UserID:   existingUser.ID,
			Username: username,
			UserType: string(entity.UserStudent),
		}, nil
	}

	// 2. Create new user
	password := utils.GenerateRandomPassword(10)
	hashedPassword, _ := utils.HashPassword(password)
	
	user := &entity.User{
		Username:           &student.StudentNumber,
		Password:           hashedPassword,
		UserType:           entity.UserStudent,
		IsActive:           true,
		MustChangePassword: true,
		SchoolID:           []int64{schoolID},
	}

	if student.FullName != nil {
		user.FullName = student.FullName
	}
	if student.Email != nil {
		user.Email = *student.Email
	}
	if student.Phone != nil {
		user.PhoneNumber = student.Phone
	}

	if err := uc.userUsecase.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	usernameResp := ""
	if user.Username != nil {
		usernameResp = *user.Username
	}

	return &entity.EmployeeUserResponse{
		UserID:   user.ID,
		Username: usernameResp,
		Password: password,
		UserType: string(entity.UserStudent),
	}, nil
}

func (uc *peopleUsecase) ensureUserForParent(ctx context.Context, schoolID int64, parent *entity.Parent) (*entity.EmployeeUserResponse, error) {
	var username string
	if parent.Email != nil && *parent.Email != "" {
		username = *parent.Email
	} else if parent.Phone != nil && *parent.Phone != "" {
		username = *parent.Phone
	} else if parent.FullName != nil {
		username = strings.ReplaceAll(strings.ToLower(*parent.FullName), " ", ".") + ".parent"
	} else {
		username = "parent." + utils.GenerateRandomPassword(5)
	}

	// 1. Try to find existing user
	existingUser, err := uc.userUsecase.GetUserByUsername(username)
	if err == nil && existingUser != nil {
		// Update profile
		updated := false
		if parent.FullName != nil && (existingUser.FullName == nil || *existingUser.FullName != *parent.FullName) {
			existingUser.FullName = parent.FullName
			updated = true
		}
		if parent.Email != nil && existingUser.Email != *parent.Email {
			existingUser.Email = *parent.Email
			updated = true
		}
		if parent.Phone != nil && (existingUser.PhoneNumber == nil || *existingUser.PhoneNumber != *parent.Phone) {
			existingUser.PhoneNumber = parent.Phone
			updated = true
		}

		if updated {
			_ = uc.userUsecase.UpdateUser(ctx, schoolID, existingUser)
		}

		usernameResp := ""
		if existingUser.Username != nil {
			usernameResp = *existingUser.Username
		}
		return &entity.EmployeeUserResponse{
			UserID:   existingUser.ID,
			Username: usernameResp,
			UserType: string(entity.UserParent),
		}, nil
	}

	// 2. Create new user
	password := utils.GenerateRandomPassword(10)
	hashedPassword, _ := utils.HashPassword(password)
	
	user := &entity.User{
		Username:           &username,
		Password:           hashedPassword,
		UserType:           entity.UserParent,
		IsActive:           true,
		MustChangePassword: true,
		SchoolID:           []int64{schoolID},
	}

	if parent.FullName != nil {
		user.FullName = parent.FullName
	}
	if parent.Email != nil {
		user.Email = *parent.Email
	}
	if parent.Phone != nil {
		user.PhoneNumber = parent.Phone
	}

	if err := uc.userUsecase.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	usernameResp := ""
	if user.Username != nil {
		usernameResp = *user.Username
	}

	return &entity.EmployeeUserResponse{
		UserID:   user.ID,
		Username: usernameResp,
		Password: password,
		UserType: string(entity.UserParent),
	}, nil
}

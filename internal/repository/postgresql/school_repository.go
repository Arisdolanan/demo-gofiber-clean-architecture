package postgresql

import (
	"context"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/jmoiron/sqlx"
)

type SchoolRepository interface {
	Create(ctx context.Context, school *entity.School) error
	Update(ctx context.Context, school *entity.School) error
	FindByID(ctx context.Context, id int64) (*entity.School, error)
	FindByCode(ctx context.Context, code string) (*entity.School, error)
	FindAll(ctx context.Context) ([]*entity.School, error)
	
	// App Package related
	CreatePackage(ctx context.Context, pkg *entity.AppPackage) error
	FindPackageByCode(ctx context.Context, code string) (*entity.AppPackage, error)
	FindAllPackages(ctx context.Context, activeOnly bool) ([]*entity.AppPackage, error)

	// License related
	CreateLicense(ctx context.Context, license *entity.SchoolLicense) error
	FindLicenseByKey(ctx context.Context, key string) (*entity.SchoolLicense, error)
	FindLicenseBySchoolID(ctx context.Context, schoolID int64) (*entity.SchoolLicense, error)
}

type schoolRepository struct {
	*BaseRepository[entity.School]
	packageRepo *BaseRepository[entity.AppPackage]
	licenseRepo *BaseRepository[entity.SchoolLicense]
}

func NewSchoolRepository(db *sqlx.DB) SchoolRepository {
	return &schoolRepository{
		BaseRepository: NewBaseRepository[entity.School](db, "schools"),
		packageRepo:    NewBaseRepository[entity.AppPackage](db, "app_packages"),
		licenseRepo:    NewBaseRepository[entity.SchoolLicense](db, "school_licenses"),
	}
}

func (r *schoolRepository) Create(ctx context.Context, school *entity.School) error {
	return r.BaseRepository.Create(ctx, school)
}

func (r *schoolRepository) Update(ctx context.Context, school *entity.School) error {
	return r.BaseRepository.Update(ctx, school, "id = $1", school.ID)
}

func (r *schoolRepository) FindByID(ctx context.Context, id int64) (*entity.School, error) {
	return r.BaseRepository.FindByID(ctx, id)
}

func (r *schoolRepository) FindByCode(ctx context.Context, code string) (*entity.School, error) {
	return r.BaseRepository.FindOne(ctx, "code = $1", code)
}

func (r *schoolRepository) FindAll(ctx context.Context) ([]*entity.School, error) {
	return r.BaseRepository.FindAll(ctx, "")
}

// App Package Implementation
func (r *schoolRepository) CreatePackage(ctx context.Context, pkg *entity.AppPackage) error {
	return r.packageRepo.Create(ctx, pkg)
}

func (r *schoolRepository) FindPackageByCode(ctx context.Context, code string) (*entity.AppPackage, error) {
	return r.packageRepo.FindOne(ctx, "code = $1", code)
}

func (r *schoolRepository) FindAllPackages(ctx context.Context, activeOnly bool) ([]*entity.AppPackage, error) {
	if activeOnly {
		return r.packageRepo.FindAll(ctx, "is_active = true")
	}
	return r.packageRepo.FindAll(ctx, "")
}

// License Implementation
func (r *schoolRepository) CreateLicense(ctx context.Context, license *entity.SchoolLicense) error {
	return r.licenseRepo.Create(ctx, license)
}

func (r *schoolRepository) FindLicenseByKey(ctx context.Context, key string) (*entity.SchoolLicense, error) {
	return r.licenseRepo.FindOne(ctx, "license_key = $1", key)
}

func (r *schoolRepository) FindLicenseBySchoolID(ctx context.Context, schoolID int64) (*entity.SchoolLicense, error) {
	return r.licenseRepo.FindOne(ctx, "school_id = $1 AND status = 'active'", schoolID)
}

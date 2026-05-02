package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/repository/postgresql"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type SchoolUsecase interface {
	RegisterSchool(ctx context.Context, school *entity.School) error
	GetSchoolByID(ctx context.Context, id int64) (*entity.School, error)
	GetAllSchools(ctx context.Context) ([]*entity.School, error)
	UpdateSchool(ctx context.Context, school *entity.School) error
	
	// Package management
	CreatePackage(ctx context.Context, pkg *entity.AppPackage) error
	GetActivePackages(ctx context.Context) ([]*entity.AppPackage, error)

	// Licensing
	AssignLicense(ctx context.Context, schoolID int64, packageCode string) (*entity.SchoolLicense, error)
}

type schoolUsecase struct {
	repo     postgresql.SchoolRepository
	validate *validator.Validate
	log      *logrus.Logger
}

func NewSchoolUsecase(repo postgresql.SchoolRepository, validate *validator.Validate, log *logrus.Logger) SchoolUsecase {
	return &schoolUsecase{
		repo:     repo,
		validate: validate,
		log:      log,
	}
}

func (uc *schoolUsecase) RegisterSchool(ctx context.Context, school *entity.School) error {
	if err := uc.validate.Struct(school); err != nil {
		return err
	}

	// Check if code already exists
	existing, err := uc.repo.FindByCode(ctx, school.Code)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("school code already exists")
	}

	school.Status = entity.SchoolActive
	return uc.repo.Create(ctx, school)
}

func (uc *schoolUsecase) GetSchoolByID(ctx context.Context, id int64) (*entity.School, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *schoolUsecase) GetAllSchools(ctx context.Context) ([]*entity.School, error) {
	return uc.repo.FindAll(ctx)
}

func (uc *schoolUsecase) CreatePackage(ctx context.Context, pkg *entity.AppPackage) error {
	if err := uc.validate.Struct(pkg); err != nil {
		return err
	}
	return uc.repo.CreatePackage(ctx, pkg)
}

func (uc *schoolUsecase) GetActivePackages(ctx context.Context) ([]*entity.AppPackage, error) {
	return uc.repo.FindAllPackages(ctx, true)
}

func (uc *schoolUsecase) AssignLicense(ctx context.Context, schoolID int64, packageCode string) (*entity.SchoolLicense, error) {
	// Find package
	pkg, err := uc.repo.FindPackageByCode(ctx, packageCode)
	if err != nil {
		return nil, err
	}
	if pkg == nil {
		return nil, errors.New("package not found")
	}

	// Create license
	license := &entity.SchoolLicense{
		SchoolID:     schoolID,
		AppPackageID: pkg.ID,
		LicenseKey:   uc.generateLicenseKey(),
		StartDate:    time.Now(),
		EndDate:      time.Now().AddDate(1, 0, 0), // Default 1 year
		Status:       entity.LicenseActive,
	}

	if err := uc.repo.CreateLicense(ctx, license); err != nil {
		return nil, err
	}

	return license, nil
}

func (uc *schoolUsecase) UpdateSchool(ctx context.Context, school *entity.School) error {
	if err := uc.validate.Struct(school); err != nil {
		return err
	}
	return uc.repo.Update(ctx, school)
}

func (uc *schoolUsecase) generateLicenseKey() string {
	// Simple license key generation for now
	return "LIC-" + time.Now().Format("20060102150405")
}

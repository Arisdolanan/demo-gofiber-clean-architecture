package database

import (
	"context"
	"fmt"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/jmoiron/sqlx"
)

type Seeder struct {
	db *sqlx.DB
}

func NewSeeder(db *sqlx.DB) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) SeedAll(ctx context.Context) error {
	if err := s.SeedAppPackages(ctx); err != nil {
		return err
	}
	if err := s.SeedPermissions(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Seeder) SeedAppPackages(ctx context.Context) error {
	packages := []entity.AppPackage{
		{Code: "BASIC", Name: "Basic Package", PriceMonthly: 50, PriceYearly: 500, MaxStudents: 500, MaxTeachers: 50, IsActive: true},
		{Code: "PRO", Name: "Professional Package", PriceMonthly: 150, PriceYearly: 1500, MaxStudents: 2000, MaxTeachers: 200, IsActive: true},
		{Code: "ULTIMATE", Name: "Ultimate Package", PriceMonthly: 500, PriceYearly: 5000, MaxStudents: 10000, MaxTeachers: 1000, IsActive: true},
	}

	for _, pkg := range packages {
		query := `
			INSERT INTO app_packages (code, name, price_monthly, price_yearly, max_students, max_teachers, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (code) DO UPDATE SET 
				name = EXCLUDED.name, 
				price_monthly = EXCLUDED.price_monthly, 
				price_yearly = EXCLUDED.price_yearly
		`
		_, err := s.db.ExecContext(ctx, query, pkg.Code, pkg.Name, pkg.PriceMonthly, pkg.PriceYearly, pkg.MaxStudents, pkg.MaxTeachers, pkg.IsActive)
		if err != nil {
			return fmt.Errorf("failed to seed package %s: %w", pkg.Code, err)
		}
	}
	return nil
}

func (s *Seeder) SeedPermissions(ctx context.Context) error {
	permissions := []entity.Permission{
		{PermissionCode: "SCHOOL_VIEW", ModuleName: "SaaS", PermissionName: "View School Details"},
		{PermissionCode: "SCHOOL_EDIT", ModuleName: "SaaS", PermissionName: "Edit School Details"},
		{PermissionCode: "ACADEMIC_MANAGE", ModuleName: "Academic", PermissionName: "Manage Academic Structure"},
		{PermissionCode: "PEOPLE_MANAGE", ModuleName: "People", PermissionName: "Manage Teachers and Students"},
		{PermissionCode: "EXAM_MANAGE", ModuleName: "Assessment", PermissionName: "Manage Exams and Marks"},
		{PermissionCode: "ATTENDANCE_MANAGE", ModuleName: "Attendance", PermissionName: "Manage Attendance"},
	}

	for _, p := range permissions {
		query := `
			INSERT INTO permissions (permission_code, module_name, permission_name)
			VALUES ($1, $2, $3)
			ON CONFLICT (permission_code) DO NOTHING
		`
		_, err := s.db.ExecContext(ctx, query, p.PermissionCode, p.ModuleName, p.PermissionName)
		if err != nil {
			return fmt.Errorf("failed to seed permission %s: %w", p.PermissionCode, err)
		}
	}
	return nil
}

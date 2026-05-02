package entity

import "time"

type SchoolLevel string

const (
	LevelSD  SchoolLevel = "SD"
	LevelSMP SchoolLevel = "SMP"
	LevelSMA SchoolLevel = "SMA"
)

type SchoolStatus string

const (
	SchoolActive    SchoolStatus = "active"
	SchoolInactive  SchoolStatus = "inactive"
	SchoolSuspended SchoolStatus = "suspended"
)

type School struct {
	ID                    int64        `json:"id" db:"id"`
	Code                  string       `json:"code" db:"code" validate:"required"`
	Name                  string       `json:"name" db:"name" validate:"required"`
	Email                 *string      `json:"email" db:"email" validate:"omitempty,email"`
	Phone                 *string      `json:"phone" db:"phone"`
	Address               *string      `json:"address" db:"address"`
	City                  *string      `json:"city" db:"city"`
	Province              *string      `json:"province" db:"province"`
	Country               *string      `json:"country" db:"country"`
	SchoolLevel           *SchoolLevel `json:"school_level" db:"school_level"`
	Npsn                  *string      `json:"npsn" db:"npsn"`
	PrincipalName         *string      `json:"principal_name" db:"principal_name"`
	Accreditation         *string      `json:"accreditation" db:"accreditation"`
	Logo                  *string      `json:"logo" db:"logo"`
	Status                SchoolStatus `json:"status" db:"status"`
	SubscriptionStartDate *time.Time   `json:"subscription_start_date,omitempty" db:"subscription_start_date"`
	SubscriptionEndDate   *time.Time   `json:"subscription_end_date,omitempty" db:"subscription_end_date"`
	BaseEntity
}

type AppPackage struct {
	ID           int64   `json:"id" db:"id"`
	Code         string  `json:"code" db:"code" validate:"required"`
	Name         string   `json:"name" db:"name" validate:"required"`
	Description  *string  `json:"description" db:"description"`
	PriceMonthly float64 `json:"price_monthly" db:"price_monthly"`
	PriceYearly  float64 `json:"price_yearly" db:"price_yearly"`
	MaxStudents  int     `json:"max_students" db:"max_students"`
	MaxTeachers  int     `json:"max_teachers" db:"max_teachers"`
	IsActive     bool    `json:"is_active" db:"is_active"`
	BaseEntity
}

type LicenseStatus string

const (
	LicenseActive    LicenseStatus = "active"
	LicenseExpired   LicenseStatus = "expired"
	LicenseSuspended LicenseStatus = "suspended"
)

type SchoolLicense struct {
	ID           int64         `json:"id" db:"id"`
	SchoolID     int64         `json:"school_id" db:"school_id"`
	AppPackageID int64         `json:"app_package_id" db:"app_package_id"`
	LicenseKey   string        `json:"license_key" db:"license_key"`
	StartDate    time.Time     `json:"start_date" db:"start_date"`
	EndDate      time.Time     `json:"end_date" db:"end_date"`
	Status       LicenseStatus `json:"status" db:"status"`
	BaseEntity
}

type PaymentStatus string

const (
	PaymentPending  PaymentStatus = "pending"
	PaymentVerified PaymentStatus = "verified"
	PaymentRejected PaymentStatus = "rejected"
)

type Payment struct {
	ID              int64         `json:"id" db:"id"`
	SchoolID        int64         `json:"school_id" db:"school_id"`
	SchoolLicenseID *int64        `json:"school_license_id,omitempty" db:"school_license_id"`
	Amount          float64       `json:"amount" db:"amount" validate:"required"`
	PaymentMethod   *string       `json:"payment_method" db:"payment_method"`
	PaymentDate     *time.Time    `json:"payment_date,omitempty" db:"payment_date"`
	ReferenceNumber *string       `json:"reference_number" db:"reference_number"`
	ProofFileURL    *string       `json:"proof_file_url" db:"proof_file_url"`
	Status          PaymentStatus `json:"status" db:"status"`
	Notes           string        `json:"notes" db:"notes"`
	BaseEntity
}

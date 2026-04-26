package entity

import "time"

type Role struct {
	ID           int64   `json:"id" db:"id"`
	SchoolID     *int64  `json:"school_id,omitempty" db:"school_id"` // NULL for system roles
	Code         string  `json:"code" db:"code" validate:"required"`
	Name         string  `json:"name" db:"name" validate:"required"`
	Description  string  `json:"description" db:"description"`
	IsSystemRole bool    `json:"is_system_role" db:"is_system_role"`
	BaseEntity
}

type Permission struct {
	ID             int64  `json:"id" db:"id"`
	PermissionCode string `json:"permission_code" db:"permission_code" validate:"required"`
	ModuleName     string `json:"module_name" db:"module_name" validate:"required"`
	PermissionName string `json:"permission_name" db:"permission_name" validate:"required"`
	Description    string `json:"description" db:"description"`
	BaseEntity
}

type RolePermission struct {
	ID           int64     `json:"id" db:"id"`
	RoleID       int64     `json:"role_id" db:"role_id"`
	PermissionID int64     `json:"permission_id" db:"permission_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type UserRole struct {
	ID         int64      `json:"id" db:"id"`
	UserID     int64      `json:"user_id" db:"user_id"`
	RoleID     int64      `json:"role_id" db:"role_id"`
	AssignedAt time.Time  `json:"assigned_at" db:"assigned_at"`
	AssignedBy *int64     `json:"assigned_by,omitempty" db:"assigned_by"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}

type ActivityLog struct {
	ID          int64      `json:"id" db:"id"`
	UserID      *int64     `json:"user_id,omitempty" db:"user_id"`
	SchoolID    *int64     `json:"school_id,omitempty" db:"school_id"`
	Action      string     `json:"action" db:"action"`
	Module      string     `json:"module" db:"module"`
	Description string     `json:"description" db:"description"`
	IPAddress   string     `json:"ip_address" db:"ip_address"`
	UserAgent   string     `json:"user_agent" db:"user_agent"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

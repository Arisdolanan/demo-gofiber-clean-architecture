package entity

import (
	"time"

	"github.com/lib/pq"
)

// RoleResponse is the safe public representation of a Role (no sensitive fields).
type RoleResponse struct {
	ID           int64   `json:"id"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Description  *string `json:"description,omitempty"`
	IsSystemRole bool    `json:"is_system_role"`
}

// PermissionResponse is the safe public representation of a Permission.
type PermissionResponse struct {
	ID             int64  `json:"id"`
	PermissionCode string `json:"permission_code"`
	ModuleName     string `json:"module_name"`
	PermissionName string `json:"permission_name"`
	Description    string `json:"description"`
}

// SchoolResponse is the safe public representation of a School (minimal info for login context).
type SchoolResponse struct {
	ID     int64        `json:"id"`
	Code   string       `json:"code"`
	Name   string       `json:"name"`
	Status SchoolStatus `json:"status"`
}

// UserResponse is the safe public representation of a User.
// Sensitive fields (password, tokens, audit metadata) are intentionally excluded.
type UserResponse struct {
	ID                int64                `json:"id"`
	Username          string               `json:"username"`
	Email             string               `json:"email"`
	UserType          UserType             `json:"user_type"`
	IsActive          bool                 `json:"is_active"`
	SchoolID          *int64               `json:"school_id,omitempty"`
	LastLoginAt       *time.Time           `json:"last_login_at,omitempty"`
	EmailVerifiedAt   *time.Time           `json:"email_verified_at,omitempty"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
	Role              *RoleResponse        `json:"role,omitempty"`
	Permissions       []PermissionResponse `json:"permissions,omitempty"`
	School            *SchoolResponse      `json:"school,omitempty"`
	AccessibleSchools []SchoolResponse     `json:"accessible_schools,omitempty"`
}

// NewUserResponse maps a User entity to a safe UserResponse DTO.
func NewUserResponse(u *User) *UserResponse {
	if u == nil {
		return nil
	}

	username := ""
	if u.Username != nil {
		username = *u.Username
	}

	resp := &UserResponse{
		ID:              u.ID,
		Username:        username,
		Email:           u.Email,
		UserType:        u.UserType,
		IsActive:        u.IsActive,
		LastLoginAt:     u.LastLoginAt,
		EmailVerifiedAt: u.EmailVerifiedAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}

	if u.School != nil {
		resp.SchoolID = &u.School.ID
	}

	if u.Role != nil {
		resp.Role = &RoleResponse{
			ID:           u.Role.ID,
			Code:         u.Role.Code,
			Name:         u.Role.Name,
			Description:  u.Role.Description,
			IsSystemRole: u.Role.IsSystemRole,
		}
	}

	if len(u.Permissions) > 0 {
		resp.Permissions = make([]PermissionResponse, len(u.Permissions))
		for i, p := range u.Permissions {
			resp.Permissions[i] = PermissionResponse{
				ID:             p.ID,
				PermissionCode: p.PermissionCode,
				ModuleName:     p.ModuleName,
				PermissionName: p.PermissionName,
				Description:    p.Description,
			}
		}
	}

	if u.School != nil {
		resp.School = &SchoolResponse{
			ID:     u.School.ID,
			Code:   u.School.Code,
			Name:   u.School.Name,
			Status: u.School.Status,
		}
	}

	if len(u.AccessibleSchools) > 0 {
		resp.AccessibleSchools = make([]SchoolResponse, len(u.AccessibleSchools))
		for i, s := range u.AccessibleSchools {
			resp.AccessibleSchools[i] = SchoolResponse{
				ID:     s.ID,
				Code:   s.Code,
				Name:   s.Name,
				Status: s.Status,
			}
		}
	}

	return resp
}

// AuthToken is the response returned after a successful login or token refresh.
type AuthToken struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresAt    time.Time     `json:"expires_at"`
	User         *UserResponse `json:"user,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Email    string        `json:"email" validate:"required,email"`
	Password string        `json:"password" validate:"required,min=8"`
	SchoolID pq.Int64Array `json:"school_id" swaggertype:"array,number"`
	UserType UserType      `json:"user_type" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type SwitchSchoolRequest struct {
	SchoolID int64 `json:"school_id" validate:"required"`
}

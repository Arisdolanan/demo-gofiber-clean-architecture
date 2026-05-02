package entity

import "time"

type ActivityLog struct {
	ID          int64     `json:"id" db:"id"`
	UserID      *int64    `json:"user_id" db:"user_id"`
	UserName    string    `json:"user_name,omitempty" db:"user_name"`
	SchoolID    *int64    `json:"school_id" db:"school_id"`
	Action      string    `json:"action" db:"action"`
	Module      string    `json:"module" db:"module"`
	Description string    `json:"description" db:"description"`
	IPAddress   *string   `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent   *string   `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

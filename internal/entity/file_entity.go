package entity

import "time"

// File represents a file in the system
type File struct {
	ID           int64      `json:"id" db:"id"`
	UserID       int64      `json:"user_id" db:"user_id"`
	Filename     string     `json:"filename" db:"filename" validate:"required,min=1,max=255"`
	OriginalName string     `json:"original_name" db:"original_name" validate:"required,min=1,max=255"`
	MimeType     string     `json:"mime_type" db:"mime_type" validate:"required,min=1,max=100"`
	Size         int64      `json:"size" db:"size" validate:"required,min=1"`
	Path         string     `json:"path" db:"path" validate:"required,min=1,max=500"`
	Description  string     `json:"description" db:"description" validate:"max=1000"`
	Category     string     `json:"category" db:"category" validate:"max=50"`
	IsPublic     bool       `json:"is_public" db:"is_public"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy    *int64     `json:"created_by,omitempty" db:"created_by"`
	UpdatedBy    *int64     `json:"updated_by,omitempty" db:"updated_by"`
	DeletedBy    *int64     `json:"deleted_by,omitempty" db:"deleted_by"`
}

// FileUploadRequest represents the request for file upload
type FileUploadRequest struct {
	Description string `json:"description" validate:"max=1000"`
	Category    string `json:"category" validate:"max=50"`
	IsPublic    bool   `json:"is_public"`
}

// FileUpdateRequest represents the request for file update
type FileUpdateRequest struct {
	Description string `json:"description" validate:"max=1000"`
	Category    string `json:"category" validate:"max=50"`
	IsPublic    bool   `json:"is_public"`
}

// FileListResponse represents the response for file list
type FileListResponse struct {
	Files      []File `json:"files"`
	TotalCount int64  `json:"total_count"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

// FileDownloadResponse represents the response for file download
type FileDownloadResponse struct {
	Filename string `json:"filename"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
	Data     []byte `json:"-"`
}

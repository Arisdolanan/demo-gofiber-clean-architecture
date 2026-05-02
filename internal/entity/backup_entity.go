package entity

import "time"

// BackupRecord represents a record in the backups table
type BackupRecord struct {
	ID          int64     `json:"id"`
	Filename    string    `json:"filename"`
	StoragePath string    `json:"storage_path"`
	SizeBytes   int64     `json:"size_bytes"`
	Status      string    `json:"status"`
	SchoolID    int64     `json:"school_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// RestoreRequest represents a request to restore from a backup
type RestoreRequest struct {
	Filename string `json:"filename" validate:"required"`
}

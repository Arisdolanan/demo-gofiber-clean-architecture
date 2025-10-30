package postgresql

import (
	"context"
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/jmoiron/sqlx"
)

// FileRepositoryInterface defines the interface for file repository operations
type FileRepositoryInterface interface {
	// Basic CRUD operations
	Create(ctx context.Context, file *entity.File) error
	GetByID(ctx context.Context, id int64) (*entity.File, error)
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.File, error)
	GetByUserIDAndCategory(ctx context.Context, userID int64, category string, limit, offset int) ([]*entity.File, error)
	GetPublicFiles(ctx context.Context, limit, offset int) ([]*entity.File, error)
	GetPrivateFilesByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.File, error)
	Update(ctx context.Context, file *entity.File) error
	Delete(ctx context.Context, id int64) error
	HardDelete(ctx context.Context, id int64) error

	// Advanced queries
	GetByFilename(ctx context.Context, filename string) (*entity.File, error)
	CountByUserID(ctx context.Context, userID int64) (int64, error)
	CountByUserIDAndCategory(ctx context.Context, userID int64, category string) (int64, error)
	CountPublicFiles(ctx context.Context) (int64, error)
	CountPrivateFilesByUserID(ctx context.Context, userID int64) (int64, error)
	GetFilesByMimeType(ctx context.Context, mimeType string, limit, offset int) ([]*entity.File, error)
	GetUserStorageUsage(ctx context.Context, userID int64) (int64, error)

	// Search operations
	SearchFiles(ctx context.Context, userID int64, query string, limit, offset int) ([]*entity.File, error)
	SearchPublicFiles(ctx context.Context, query string, limit, offset int) ([]*entity.File, error)
}

// FileRepository implements the FileRepositoryInterface
type FileRepository struct {
	*BaseRepository[entity.File]
}

// NewFileRepository creates a new file repository
func NewFileRepository(db *sqlx.DB) FileRepositoryInterface {
	return &FileRepository{
		BaseRepository: NewBaseRepository[entity.File](db, "files"),
	}
}

// Create creates a new file record
func (r *FileRepository) Create(ctx context.Context, file *entity.File) error {
	return r.BaseRepository.Create(ctx, file)
}

// GetByID retrieves a file by its ID
func (r *FileRepository) GetByID(ctx context.Context, id int64) (*entity.File, error) {
	var file entity.File
	query := `SELECT id, user_id, filename, original_name, mime_type, size, path, description, category, is_public, created_at, updated_at, deleted_at FROM files WHERE id = $1 AND deleted_at IS NULL`
	err := r.GetContext(ctx, &file, query, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}

// GetByUserID retrieves files by user ID with pagination
func (r *FileRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.File, error) {
	var files []*entity.File
	query := `
		SELECT id, user_id, filename, original_name, mime_type, size, path, description, category, is_public, created_at, updated_at, deleted_at
		FROM files 
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`

	err := r.SelectContext(ctx, &files, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// GetByUserIDAndCategory retrieves files by user ID and category with pagination
func (r *FileRepository) GetByUserIDAndCategory(ctx context.Context, userID int64, category string, limit, offset int) ([]*entity.File, error) {
	var files []*entity.File
	query := `
		SELECT id, user_id, filename, original_name, mime_type, size, path, description, category, is_public, created_at, updated_at, deleted_at
		FROM files 
		WHERE user_id = $1 AND category = $2 AND deleted_at IS NULL
		ORDER BY created_at DESC 
		LIMIT $3 OFFSET $4
	`

	err := r.SelectContext(ctx, &files, query, userID, category, limit, offset)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// GetPublicFiles retrieves public files with pagination
func (r *FileRepository) GetPublicFiles(ctx context.Context, limit, offset int) ([]*entity.File, error) {
	var files []*entity.File
	query := `
		SELECT id, user_id, filename, original_name, mime_type, size, path, description, category, is_public, created_at, updated_at, deleted_at
		FROM files 
		WHERE is_public = true AND deleted_at IS NULL
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`

	err := r.SelectContext(ctx, &files, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// GetPrivateFilesByUserID retrieves private files for a specific user with pagination
func (r *FileRepository) GetPrivateFilesByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.File, error) {
	var files []*entity.File
	query := `
		SELECT id, user_id, filename, original_name, mime_type, size, path, description, category, is_public, created_at, updated_at, deleted_at
		FROM files 
		WHERE user_id = $1 AND is_public = false AND deleted_at IS NULL
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`

	err := r.SelectContext(ctx, &files, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// Update updates an existing file record
func (r *FileRepository) Update(ctx context.Context, file *entity.File) error {
	query := `
		UPDATE files 
		SET user_id = $2, filename = $3, original_name = $4, mime_type = $5, 
		    size = $6, path = $7, description = $8, category = $9, 
		    is_public = $10, updated_at = $11
		WHERE id = $1
	`
	_, err := r.ExecContext(ctx, query,
		file.ID, file.UserID, file.Filename, file.OriginalName, file.MimeType,
		file.Size, file.Path, file.Description, file.Category, file.IsPublic, file.UpdatedAt)
	return err
}

// Delete soft deletes a file record (sets deleted_at timestamp)
func (r *FileRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE files SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.ExecContext(ctx, query, time.Now(), id)
	return err
}

// HardDelete permanently removes a file record from database
func (r *FileRepository) HardDelete(ctx context.Context, id int64) error {
	query := `DELETE FROM files WHERE id = $1`
	_, err := r.ExecContext(ctx, query, id)
	return err
}

// GetByFilename retrieves a file by its filename
func (r *FileRepository) GetByFilename(ctx context.Context, filename string) (*entity.File, error) {
	var file entity.File
	query := `SELECT id, user_id, filename, original_name, mime_type, size, path, description, category, is_public, created_at, updated_at, deleted_at FROM files WHERE filename = $1 AND deleted_at IS NULL`
	err := r.GetContext(ctx, &file, query, filename)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}

// CountByUserID counts files by user ID
func (r *FileRepository) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	return r.BaseRepository.Count(ctx, "user_id = $1 AND deleted_at IS NULL", userID)
}

// CountByUserIDAndCategory counts files by user ID and category
func (r *FileRepository) CountByUserIDAndCategory(ctx context.Context, userID int64, category string) (int64, error) {
	return r.BaseRepository.Count(ctx, "user_id = $1 AND category = $2 AND deleted_at IS NULL", userID, category)
}

// CountPublicFiles counts public files
func (r *FileRepository) CountPublicFiles(ctx context.Context) (int64, error) {
	return r.BaseRepository.Count(ctx, "is_public = true AND deleted_at IS NULL")
}

// CountPrivateFilesByUserID counts private files for a specific user
func (r *FileRepository) CountPrivateFilesByUserID(ctx context.Context, userID int64) (int64, error) {
	return r.BaseRepository.Count(ctx, "user_id = $1 AND is_public = false AND deleted_at IS NULL", userID)
}

// GetFilesByMimeType retrieves files by MIME type with pagination
func (r *FileRepository) GetFilesByMimeType(ctx context.Context, mimeType string, limit, offset int) ([]*entity.File, error) {
	var files []*entity.File
	query := `
		SELECT id, user_id, filename, original_name, mime_type, size, path, description, category, is_public, created_at, updated_at, deleted_at
		FROM files 
		WHERE mime_type = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`

	err := r.SelectContext(ctx, &files, query, mimeType, limit, offset)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// GetUserStorageUsage calculates total storage usage for a user
func (r *FileRepository) GetUserStorageUsage(ctx context.Context, userID int64) (int64, error) {
	var totalSize int64
	query := `SELECT COALESCE(SUM(size), 0) FROM files WHERE user_id = $1 AND deleted_at IS NULL`

	err := r.GetContext(ctx, &totalSize, query, userID)
	if err != nil {
		return 0, err
	}

	return totalSize, nil
}

// SearchFiles searches files by filename or description for a specific user
func (r *FileRepository) SearchFiles(ctx context.Context, userID int64, query string, limit, offset int) ([]*entity.File, error) {
	var files []*entity.File
	sqlQuery := `
		SELECT id, user_id, filename, original_name, mime_type, size, path, description, category, is_public, created_at, updated_at, deleted_at
		FROM files 
		WHERE user_id = $1 AND deleted_at IS NULL AND (
			LOWER(original_name) LIKE LOWER($2) OR 
			LOWER(description) LIKE LOWER($2) OR 
			LOWER(category) LIKE LOWER($2)
		)
		ORDER BY created_at DESC 
		LIMIT $3 OFFSET $4
	`

	searchQuery := "%" + query + "%"
	err := r.SelectContext(ctx, &files, sqlQuery, userID, searchQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// SearchPublicFiles searches public files by filename or description
func (r *FileRepository) SearchPublicFiles(ctx context.Context, query string, limit, offset int) ([]*entity.File, error) {
	var files []*entity.File
	sqlQuery := `
		SELECT id, user_id, filename, original_name, mime_type, size, path, description, category, is_public, created_at, updated_at, deleted_at
		FROM files 
		WHERE is_public = true AND deleted_at IS NULL AND (
			LOWER(original_name) LIKE LOWER($1) OR 
			LOWER(description) LIKE LOWER($1) OR 
			LOWER(category) LIKE LOWER($1)
		)
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`

	searchQuery := "%" + query + "%"
	err := r.SelectContext(ctx, &files, sqlQuery, searchQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	return files, nil
}

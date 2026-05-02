package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/repository/postgresql"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

// FileUseCase defines the interface for file operations
type FileUseCase interface {
	// File management operations
	UploadFile(ctx context.Context, schoolID, userID int64, file *multipart.FileHeader, req *entity.FileUploadRequest) (*entity.File, error)
	GetFileByID(ctx context.Context, schoolID, id int64) (*entity.File, error)
	GetUserFiles(ctx context.Context, schoolID, userID int64, page, pageSize int) (*entity.FileListResponse, error)
	GetUserFilesByCategory(ctx context.Context, schoolID, userID int64, category string, page, pageSize int) (*entity.FileListResponse, error)
	GetPublicFiles(ctx context.Context, schoolID int64, page, pageSize int) (*entity.FileListResponse, error)
	GetPrivateFilesByUserID(ctx context.Context, schoolID, userID int64, page, pageSize int) (*entity.FileListResponse, error)
	UpdateFile(ctx context.Context, schoolID, id, userID int64, req *entity.FileUpdateRequest) (*entity.File, error)
	DeleteFile(ctx context.Context, schoolID, id, userID int64) error

	// Advanced operations
	DownloadFile(ctx context.Context, schoolID, id, userID int64) (*entity.FileDownloadResponse, error)
	DownloadPublicFile(ctx context.Context, schoolID, id int64) (*entity.FileDownloadResponse, error)
	GetUserStorageUsage(ctx context.Context, schoolID, userID int64) (int64, error)
	SearchFiles(ctx context.Context, schoolID, userID int64, query string, page, pageSize int) (*entity.FileListResponse, error)
	SearchPublicFiles(ctx context.Context, schoolID int64, query string, page, pageSize int) (*entity.FileListResponse, error)

	// File validation and utilities
	ValidateFileType(filename, mimeType string) error
	ValidateFileSize(size int64) error
}

type fileUseCase struct {
	fileRepo     postgresql.FileRepositoryInterface
	log          *logrus.Logger
	validate     *validator.Validate
	uploadPath   string
	maxFileSize  int64
	allowedTypes []string
}

// NewFileUseCase creates a new file use case
func NewFileUseCase(
	fileRepo postgresql.FileRepositoryInterface,
	log *logrus.Logger,
	validate *validator.Validate,
) FileUseCase {
	// Create upload directory if it doesn't exist
	uploadPath := "./storage/uploads"
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		log.Errorf("Failed to create upload directory: %v", err)
	}

	return &fileUseCase{
		fileRepo:    fileRepo,
		log:         log,
		validate:    validate,
		uploadPath:  uploadPath,
		maxFileSize: 50 * 1024 * 1024, // 50MB max file size
		allowedTypes: []string{
			// Images
			"image/jpeg", "image/jpg", "image/png", "image/gif", "image/bmp", "image/webp",
			// Documents
			"application/pdf", "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			"application/vnd.ms-powerpoint", "application/vnd.openxmlformats-officedocument.presentationml.presentation",
			"text/plain", "text/csv",
			// Archives
			"application/zip", "application/x-rar-compressed", "application/x-7z-compressed",
			// Videos
			"video/mp4", "video/avi", "video/mov", "video/wmv", "video/flv",
			// Audio
			"audio/mp3", "audio/wav", "audio/ogg", "audio/m4a",
		},
	}
}

// UploadFile handles file upload with validation and storage
func (uc *fileUseCase) UploadFile(ctx context.Context, schoolID, userID int64, fileHeader *multipart.FileHeader, req *entity.FileUploadRequest) (*entity.File, error) {
	// Validate request
	if err := uc.validate.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Validate file size
	if err := uc.ValidateFileSize(fileHeader.Size); err != nil {
		return nil, err
	}

	// Get MIME type from header
	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		// Fallback to detecting from filename extension
		mimeType = uc.getMimeTypeFromExtension(fileHeader.Filename)
	}

	// Validate file type
	if err := uc.ValidateFileType(fileHeader.Filename, mimeType); err != nil {
		return nil, err
	}

	// Generate unique filename
	uniqueFilename := uc.generateUniqueFilename(fileHeader.Filename)
	filePath := filepath.Join(uc.uploadPath, uniqueFilename)

	// Open uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		// Clean up on error
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Create file entity
	file := &entity.File{
		UserID:       userID,
		SchoolID:     schoolID,
		Filename:     uniqueFilename,
		OriginalName: fileHeader.Filename,
		MimeType:     mimeType,
		Size:         fileHeader.Size,
		Path:         filePath,
		Description:  req.Description,
		Category:     req.Category,
		IsPublic:     req.IsPublic,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save to database
	if err := uc.fileRepo.Create(ctx, file); err != nil {
		// Clean up file on database error
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file to database: %w", err)
	}

	uc.log.Infof("File uploaded successfully: %s by user %d", file.OriginalName, userID)
	return file, nil
}

// GetFileByID retrieves a file by its ID
func (uc *fileUseCase) GetFileByID(ctx context.Context, schoolID, id int64) (*entity.File, error) {
	file, err := uc.fileRepo.GetByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, errors.New("file not found")
	}
	return file, nil
}

// GetUserFiles retrieves files for a specific user with pagination
func (uc *fileUseCase) GetUserFiles(ctx context.Context, schoolID, userID int64, page, pageSize int) (*entity.FileListResponse, error) {
	// Normalize pagination parameters
	pagination := utils.CalculatePagination(page, pageSize, 0) // totalCount will be calculated after query

	files, err := uc.fileRepo.GetByUserID(ctx, schoolID, userID, pagination.PageSize, pagination.Offset)
	if err != nil {
		return nil, err
	}

	totalCount, err := uc.fileRepo.CountByUserID(ctx, schoolID, userID)
	if err != nil {
		return nil, err
	}

	// Recalculate pagination with actual total count
	pagination = utils.CalculatePagination(page, pageSize, totalCount)

	return &entity.FileListResponse{
		Files:      uc.convertFilePointers(files),
		TotalCount: totalCount,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: pagination.TotalPages,
	}, nil
}

// GetUserFilesByCategory retrieves files for a specific user and category with pagination
func (uc *fileUseCase) GetUserFilesByCategory(ctx context.Context, schoolID, userID int64, category string, page, pageSize int) (*entity.FileListResponse, error) {
	// Normalize pagination parameters
	pagination := utils.CalculatePagination(page, pageSize, 0) // totalCount will be calculated after query

	files, err := uc.fileRepo.GetByUserIDAndCategory(ctx, schoolID, userID, category, pagination.PageSize, pagination.Offset)
	if err != nil {
		return nil, err
	}

	totalCount, err := uc.fileRepo.CountByUserIDAndCategory(ctx, schoolID, userID, category)
	if err != nil {
		return nil, err
	}

	// Recalculate pagination with actual total count
	pagination = utils.CalculatePagination(page, pageSize, totalCount)

	return &entity.FileListResponse{
		Files:      uc.convertFilePointers(files),
		TotalCount: totalCount,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: pagination.TotalPages,
	}, nil
}

// GetPublicFiles retrieves public files with pagination
func (uc *fileUseCase) GetPublicFiles(ctx context.Context, schoolID int64, page, pageSize int) (*entity.FileListResponse, error) {
	// Normalize pagination parameters
	pagination := utils.CalculatePagination(page, pageSize, 0) // totalCount will be calculated after query

	files, err := uc.fileRepo.GetPublicFiles(ctx, schoolID, pagination.PageSize, pagination.Offset)
	if err != nil {
		return nil, err
	}

	totalCount, err := uc.fileRepo.CountPublicFiles(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	// Recalculate pagination with actual total count
	pagination = utils.CalculatePagination(page, pageSize, totalCount)

	return &entity.FileListResponse{
		Files:      uc.convertFilePointers(files),
		TotalCount: totalCount,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: pagination.TotalPages,
	}, nil
}

// GetPrivateFilesByUserID retrieves private files for a specific user with pagination
func (uc *fileUseCase) GetPrivateFilesByUserID(ctx context.Context, schoolID, userID int64, page, pageSize int) (*entity.FileListResponse, error) {
	// Normalize pagination parameters
	pagination := utils.CalculatePagination(page, pageSize, 0) // totalCount will be calculated after query

	files, err := uc.fileRepo.GetPrivateFilesByUserID(ctx, schoolID, userID, pagination.PageSize, pagination.Offset)
	if err != nil {
		return nil, err
	}

	totalCount, err := uc.fileRepo.CountPrivateFilesByUserID(ctx, schoolID, userID)
	if err != nil {
		return nil, err
	}

	// Recalculate pagination with actual total count
	pagination = utils.CalculatePagination(page, pageSize, totalCount)

	return &entity.FileListResponse{
		Files:      uc.convertFilePointers(files),
		TotalCount: totalCount,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: pagination.TotalPages,
	}, nil
}

// UpdateFile updates file metadata
func (uc *fileUseCase) UpdateFile(ctx context.Context, schoolID, id, userID int64, req *entity.FileUpdateRequest) (*entity.File, error) {
	// Validate request
	if err := uc.validate.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Get existing file
	file, err := uc.fileRepo.GetByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, errors.New("file not found")
	}

	// Check ownership
	if file.UserID != userID {
		return nil, errors.New("access denied: you can only update your own files")
	}

	// Update fields
	file.Description = req.Description
	file.Category = req.Category
	file.IsPublic = req.IsPublic
	file.UpdatedAt = time.Now()

	// Save changes
	if err := uc.fileRepo.Update(ctx, schoolID, file); err != nil {
		return nil, err
	}

	uc.log.Infof("File updated successfully: %d by user %d", id, userID)
	return file, nil
}

// DeleteFile deletes a file and its physical file
func (uc *fileUseCase) DeleteFile(ctx context.Context, schoolID, id, userID int64) error {
	// Get file
	file, err := uc.fileRepo.GetByID(ctx, schoolID, id)
	if err != nil {
		return err
	}
	if file == nil {
		return errors.New("file not found")
	}

	// Check ownership
	if file.UserID != userID {
		return errors.New("access denied: you can only delete your own files")
	}

	// Delete from database first
	if err := uc.fileRepo.Delete(ctx, schoolID, id); err != nil {
		return err
	}

	// Delete physical file
	if err := os.Remove(file.Path); err != nil {
		uc.log.Warnf("Failed to delete physical file %s: %v", file.Path, err)
		// Don't return error as database deletion was successful
	}

	uc.log.Infof("File deleted successfully: %d by user %d", id, userID)
	return nil
}

// DownloadFile prepares file for download (user files)
func (uc *fileUseCase) DownloadFile(ctx context.Context, schoolID, id, userID int64) (*entity.FileDownloadResponse, error) {
	// Get file
	file, err := uc.fileRepo.GetByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, errors.New("file not found")
	}

	// Check ownership or public access
	if file.UserID != userID && !file.IsPublic {
		return nil, errors.New("access denied")
	}

	// Read file content
	data, err := os.ReadFile(file.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return &entity.FileDownloadResponse{
		Filename: file.OriginalName,
		MimeType: file.MimeType,
		Size:     file.Size,
		Data:     data,
	}, nil
}

// DownloadPublicFile prepares public file for download
func (uc *fileUseCase) DownloadPublicFile(ctx context.Context, schoolID, id int64) (*entity.FileDownloadResponse, error) {
	// Get file
	file, err := uc.fileRepo.GetByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, errors.New("file not found")
	}

	// Check if file is public
	if !file.IsPublic {
		return nil, errors.New("file is not public")
	}

	// Read file content
	data, err := os.ReadFile(file.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return &entity.FileDownloadResponse{
		Filename: file.OriginalName,
		MimeType: file.MimeType,
		Size:     file.Size,
		Data:     data,
	}, nil
}

// GetUserStorageUsage calculates total storage usage for a user
func (uc *fileUseCase) GetUserStorageUsage(ctx context.Context, schoolID, userID int64) (int64, error) {
	return uc.fileRepo.GetUserStorageUsage(ctx, schoolID, userID)
}

// SearchFiles searches user files
func (uc *fileUseCase) SearchFiles(ctx context.Context, schoolID, userID int64, query string, page, pageSize int) (*entity.FileListResponse, error) {
	// Normalize pagination parameters
	pagination := utils.CalculatePagination(page, pageSize, 0)

	if query == "" {
		return uc.GetUserFiles(ctx, schoolID, userID, pagination.Page, pagination.PageSize)
	}

	files, err := uc.fileRepo.SearchFiles(ctx, schoolID, userID, query, pagination.PageSize, pagination.Offset)
	if err != nil {
		return nil, err
	}

	// For search, we'll use the file count as approximation
	totalCount := int64(len(files))
	pagination = utils.CalculatePagination(page, pageSize, totalCount)

	return &entity.FileListResponse{
		Files:      uc.convertFilePointers(files),
		TotalCount: totalCount,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: pagination.TotalPages,
	}, nil
}

// SearchPublicFiles searches public files
func (uc *fileUseCase) SearchPublicFiles(ctx context.Context, schoolID int64, query string, page, pageSize int) (*entity.FileListResponse, error) {
	// Normalize pagination parameters
	pagination := utils.CalculatePagination(page, pageSize, 0)

	if query == "" {
		return uc.GetPublicFiles(ctx, schoolID, pagination.Page, pagination.PageSize)
	}

	files, err := uc.fileRepo.SearchPublicFiles(ctx, schoolID, query, pagination.PageSize, pagination.Offset)
	if err != nil {
		return nil, err
	}

	// For search, we'll use the file count as approximation
	totalCount := int64(len(files))
	pagination = utils.CalculatePagination(page, pageSize, totalCount)

	return &entity.FileListResponse{
		Files:      uc.convertFilePointers(files),
		TotalCount: totalCount,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: pagination.TotalPages,
	}, nil
}

// ValidateFileType validates if the file type is allowed
func (uc *fileUseCase) ValidateFileType(filename, mimeType string) error {
	// Check MIME type
	for _, allowedType := range uc.allowedTypes {
		if strings.EqualFold(mimeType, allowedType) {
			return nil
		}
	}

	return fmt.Errorf("file type not allowed: %s", mimeType)
}

// ValidateFileSize validates if the file size is within limits
func (uc *fileUseCase) ValidateFileSize(size int64) error {
	if size <= 0 {
		return errors.New("file size must be greater than 0")
	}
	if size > uc.maxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", uc.maxFileSize)
	}
	return nil
}

// Helper methods

// generateUniqueFilename generates a unique filename to prevent conflicts
func (uc *fileUseCase) generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	nameWithoutExt := strings.TrimSuffix(originalFilename, ext)

	// Generate random string
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	randomStr := hex.EncodeToString(randomBytes)

	// Add timestamp for extra uniqueness
	timestamp := time.Now().Unix()

	return fmt.Sprintf("%s_%d_%s%s", nameWithoutExt, timestamp, randomStr, ext)
}

// getMimeTypeFromExtension gets MIME type from file extension as fallback
func (uc *fileUseCase) getMimeTypeFromExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".txt":  "text/plain",
		".zip":  "application/zip",
		".mp4":  "video/mp4",
		".mp3":  "audio/mp3",
	}

	if mimeType, exists := mimeTypes[ext]; exists {
		return mimeType
	}

	return "application/octet-stream"
}

// convertFilePointers converts []*entity.File to []entity.File
func (uc *fileUseCase) convertFilePointers(files []*entity.File) []entity.File {
	result := make([]entity.File, len(files))
	for i, file := range files {
		if file != nil {
			result[i] = *file
		}
	}
	return result
}

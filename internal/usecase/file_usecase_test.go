package usecase

import (
	"context"
	"os"
	"testing"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFileRepository is a mock for FileRepositoryInterface
type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) Create(ctx context.Context, file *entity.File) error {
	args := m.Called(ctx, file)
	return args.Error(0)
}

func (m *MockFileRepository) GetByID(ctx context.Context, id int64) (*entity.File, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.File), args.Error(1)
}

func (m *MockFileRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.File, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*entity.File), args.Error(1)
}

func (m *MockFileRepository) GetByUserIDAndCategory(ctx context.Context, userID int64, category string, limit, offset int) ([]*entity.File, error) {
	args := m.Called(ctx, userID, category, limit, offset)
	return args.Get(0).([]*entity.File), args.Error(1)
}

func (m *MockFileRepository) GetPublicFiles(ctx context.Context, limit, offset int) ([]*entity.File, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*entity.File), args.Error(1)
}

func (m *MockFileRepository) GetPrivateFilesByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.File, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*entity.File), args.Error(1)
}

func (m *MockFileRepository) Update(ctx context.Context, file *entity.File) error {
	args := m.Called(ctx, file)
	return args.Error(0)
}

func (m *MockFileRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFileRepository) HardDelete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFileRepository) GetByFilename(ctx context.Context, filename string) (*entity.File, error) {
	args := m.Called(ctx, filename)
	return args.Get(0).(*entity.File), args.Error(1)
}

func (m *MockFileRepository) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockFileRepository) CountByUserIDAndCategory(ctx context.Context, userID int64, category string) (int64, error) {
	args := m.Called(ctx, userID, category)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockFileRepository) CountPublicFiles(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockFileRepository) CountPrivateFilesByUserID(ctx context.Context, userID int64) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockFileRepository) GetFilesByMimeType(ctx context.Context, mimeType string, limit, offset int) ([]*entity.File, error) {
	args := m.Called(ctx, mimeType, limit, offset)
	return args.Get(0).([]*entity.File), args.Error(1)
}

func (m *MockFileRepository) GetUserStorageUsage(ctx context.Context, userID int64) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockFileRepository) SearchFiles(ctx context.Context, userID int64, query string, limit, offset int) ([]*entity.File, error) {
	args := m.Called(ctx, userID, query, limit, offset)
	return args.Get(0).([]*entity.File), args.Error(1)
}

func (m *MockFileRepository) SearchPublicFiles(ctx context.Context, query string, limit, offset int) ([]*entity.File, error) {
	args := m.Called(ctx, query, limit, offset)
	return args.Get(0).([]*entity.File), args.Error(1)
}

func TestFileUseCase_ValidateFileType(t *testing.T) {
	// Setup
	mockRepo := new(MockFileRepository)
	logger := logrus.New()
	validate := validator.New()

	// Create temp directory for testing
	tempDir := "./test_uploads"
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	uc := &fileUseCase{
		fileRepo:    mockRepo,
		log:         logger,
		validate:    validate,
		uploadPath:  tempDir,
		maxFileSize: 50 * 1024 * 1024,
		allowedTypes: []string{
			"image/jpeg", "image/png", "application/pdf", "text/plain",
		},
	}

	tests := []struct {
		name        string
		filename    string
		mimeType    string
		expectError bool
	}{
		{
			name:        "Valid JPEG file",
			filename:    "test.jpg",
			mimeType:    "image/jpeg",
			expectError: false,
		},
		{
			name:        "Valid PNG file",
			filename:    "test.png",
			mimeType:    "image/png",
			expectError: false,
		},
		{
			name:        "Invalid file type",
			filename:    "test.exe",
			mimeType:    "application/x-executable",
			expectError: true,
		},
		{
			name:        "Valid PDF file",
			filename:    "document.pdf",
			mimeType:    "application/pdf",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.ValidateFileType(tt.filename, tt.mimeType)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileUseCase_ValidateFileSize(t *testing.T) {
	// Setup
	mockRepo := new(MockFileRepository)
	logger := logrus.New()
	validate := validator.New()

	uc := &fileUseCase{
		fileRepo:     mockRepo,
		log:          logger,
		validate:     validate,
		uploadPath:   "./test_uploads",
		maxFileSize:  10 * 1024 * 1024, // 10MB
		allowedTypes: []string{"image/jpeg"},
	}

	tests := []struct {
		name        string
		size        int64
		expectError bool
	}{
		{
			name:        "Valid file size",
			size:        5 * 1024 * 1024, // 5MB
			expectError: false,
		},
		{
			name:        "File too large",
			size:        15 * 1024 * 1024, // 15MB
			expectError: true,
		},
		{
			name:        "Zero file size",
			size:        0,
			expectError: true,
		},
		{
			name:        "Negative file size",
			size:        -1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.ValidateFileSize(tt.size)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileUseCase_GetUserStorageUsage(t *testing.T) {
	// Setup
	mockRepo := new(MockFileRepository)
	logger := logrus.New()
	validate := validator.New()

	uc := &fileUseCase{
		fileRepo:     mockRepo,
		log:          logger,
		validate:     validate,
		uploadPath:   "./test_uploads",
		maxFileSize:  50 * 1024 * 1024,
		allowedTypes: []string{"image/jpeg"},
	}

	ctx := context.Background()
	userID := int64(1)
	expectedUsage := int64(1024 * 1024) // 1MB

	// Mock expectations
	mockRepo.On("GetUserStorageUsage", ctx, userID).Return(expectedUsage, nil)

	// Execute
	usage, err := uc.GetUserStorageUsage(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUsage, usage)
	mockRepo.AssertExpectations(t)
}

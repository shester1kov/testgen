package document

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

// Mock repository
type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) Create(ctx context.Context, doc *entity.Document) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *MockDocumentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Document, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Document), args.Error(1)
}

func (m *MockDocumentRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Document, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Document), args.Error(1)
}

func (m *MockDocumentRepository) Update(ctx context.Context, doc *entity.Document) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *MockDocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDocumentRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.Document, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Document), args.Error(1)
}

func (m *MockDocumentRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDocumentRepository) CountAll(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func TestUploadUseCase_Execute(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	tempDir := filepath.Join(os.TempDir(), "test-uploads")
	defer os.RemoveAll(tempDir)

	uc := NewUploadUseCase(mockRepo, tempDir, 50*1024*1024) // 50MB max

	ctx := context.Background()
	userID := uuid.New()
	fileContent := strings.NewReader("test content")

	params := UploadParams{
		UserID:      userID,
		Title:       "Test Document",
		FileName:    "test.pdf",
		FileSize:    12, // len("test content")
		FileType:    "pdf",
		FileContent: fileContent,
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Document")).Return(nil)

	result, err := uc.Execute(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Document", result.Title)
	assert.Equal(t, entity.StatusUploaded, result.Status)
	assert.Equal(t, userID, result.UserID)
	mockRepo.AssertExpectations(t)

	// Verify file was created
	assert.FileExists(t, result.FilePath)
}

func TestUploadUseCase_Execute_InvalidFileType(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	tempDir := filepath.Join(os.TempDir(), "test-uploads")
	defer os.RemoveAll(tempDir)

	uc := NewUploadUseCase(mockRepo, tempDir, 50*1024*1024)

	ctx := context.Background()
	userID := uuid.New()
	fileContent := strings.NewReader("test content")

	params := UploadParams{
		UserID:      userID,
		Title:       "Test Document",
		FileName:    "test.exe",
		FileSize:    12,
		FileType:    "exe",
		FileContent: fileContent,
	}

	_, err := uc.Execute(ctx, params)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported file type")
}

func TestUploadUseCase_Execute_FileTooLarge(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	tempDir := filepath.Join(os.TempDir(), "test-uploads")
	defer os.RemoveAll(tempDir)

	uc := NewUploadUseCase(mockRepo, tempDir, 10*1024*1024) // 10MB max

	ctx := context.Background()
	userID := uuid.New()
	fileContent := strings.NewReader("test content")

	params := UploadParams{
		UserID:      userID,
		Title:       "Test Document",
		FileName:    "test.pdf",
		FileSize:    100 * 1024 * 1024, // 100MB
		FileType:    "pdf",
		FileContent: fileContent,
	}

	_, err := uc.Execute(ctx, params)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file size exceeds maximum")
}

func TestUploadUseCase_Execute_DBError(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	tempDir := filepath.Join(os.TempDir(), "test-uploads")
	defer os.RemoveAll(tempDir)

	uc := NewUploadUseCase(mockRepo, tempDir, 50*1024*1024)

	ctx := context.Background()
	userID := uuid.New()
	fileContent := strings.NewReader("test content")

	params := UploadParams{
		UserID:      userID,
		Title:       "Test Document",
		FileName:    "test.pdf",
		FileSize:    12,
		FileType:    "pdf",
		FileContent: fileContent,
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Document")).Return(assert.AnError)

	_, err := uc.Execute(ctx, params)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save document metadata")
	mockRepo.AssertExpectations(t)
}

func TestUploadUseCase_Execute_AllFileTypes(t *testing.T) {
	validTypes := []string{"pdf", "docx", "pptx", "txt"}

	for _, fileType := range validTypes {
		t.Run(fileType, func(t *testing.T) {
			mockRepo := new(MockDocumentRepository)
			tempDir := filepath.Join(os.TempDir(), "test-uploads-"+fileType)
			defer os.RemoveAll(tempDir)

			uc := NewUploadUseCase(mockRepo, tempDir, 50*1024*1024)

			ctx := context.Background()
			userID := uuid.New()
			fileContent := strings.NewReader("test content")

			params := UploadParams{
				UserID:      userID,
				Title:       "Test Document",
				FileName:    "test." + fileType,
				FileSize:    12,
				FileType:    fileType,
				FileContent: fileContent,
			}

			mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Document")).Return(nil)

			result, err := uc.Execute(ctx, params)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, entity.FileType(fileType), result.FileType)
			mockRepo.AssertExpectations(t)
		})
	}
}

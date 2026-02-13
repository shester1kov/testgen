package document

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

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

// Mock storage
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	args := m.Called(ctx, key, reader, size, contentType)
	return args.Error(0)
}

func (m *MockStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockStorage) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockStorage) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	args := m.Called(ctx, key, expiry)
	return args.String(0), args.Error(1)
}

func TestUploadUseCase_Execute(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	mockStorage := new(MockStorage)
	logger := zap.NewNop()

	uc := NewUploadUseCase(mockRepo, mockStorage, 50*1024*1024, logger)

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

	// Expect storage.Upload to be called
	mockStorage.On("Upload", ctx, mock.MatchedBy(func(key string) bool {
		return strings.HasPrefix(key, "documents/"+userID.String()+"/") && strings.HasSuffix(key, ".pdf")
	}), mock.Anything, int64(12), "application/pdf").Return(nil)

	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Document")).Return(nil)

	result, err := uc.Execute(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Document", result.Title)
	assert.Equal(t, entity.StatusUploaded, result.Status)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, entity.FileTypePDF, result.FileType)
	assert.Equal(t, "application/pdf", result.ContentType)
	assert.NotEmpty(t, result.ObjectKey)
	assert.Contains(t, result.ObjectKey, "documents/"+userID.String()+"/")
	assert.Contains(t, result.ObjectKey, ".pdf")

	mockRepo.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestUploadUseCase_Execute_InvalidFileType(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	mockStorage := new(MockStorage)
	logger := zap.NewNop()

	uc := NewUploadUseCase(mockRepo, mockStorage, 50*1024*1024, logger)

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

	// Storage and repo should not be called
	mockStorage.AssertNotCalled(t, "Upload")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestUploadUseCase_Execute_FileTooLarge(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	mockStorage := new(MockStorage)
	logger := zap.NewNop()

	uc := NewUploadUseCase(mockRepo, mockStorage, 10*1024*1024, logger) // 10MB max

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

	// Storage and repo should not be called
	mockStorage.AssertNotCalled(t, "Upload")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestUploadUseCase_Execute_DBError(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	mockStorage := new(MockStorage)
	logger := zap.NewNop()

	uc := NewUploadUseCase(mockRepo, mockStorage, 50*1024*1024, logger)

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

	// Storage upload succeeds
	mockStorage.On("Upload", ctx, mock.Anything, mock.Anything, int64(12), "application/pdf").Return(nil)

	// DB create fails
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Document")).Return(errors.New("db error"))

	// Expect rollback - storage.Delete should be called
	mockStorage.On("Delete", ctx, mock.Anything).Return(nil)

	_, err := uc.Execute(ctx, params)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save document metadata")

	mockRepo.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
	mockStorage.AssertCalled(t, "Delete", ctx, mock.Anything)
}

func TestUploadUseCase_Execute_AllFileTypes(t *testing.T) {
	testCases := []struct {
		fileType    string
		extension   string
		contentType string
	}{
		{"pdf", ".pdf", "application/pdf"},
		{"docx", ".docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
		{"pptx", ".pptx", "application/vnd.openxmlformats-officedocument.presentationml.presentation"},
		{"txt", ".txt", "text/plain"},
		{"md", ".md", "text/markdown"},
	}

	for _, tc := range testCases {
		t.Run(tc.fileType, func(t *testing.T) {
			mockRepo := new(MockDocumentRepository)
			mockStorage := new(MockStorage)
			logger := zap.NewNop()

			uc := NewUploadUseCase(mockRepo, mockStorage, 50*1024*1024, logger)

			ctx := context.Background()
			userID := uuid.New()
			fileContent := strings.NewReader("test content")

			params := UploadParams{
				UserID:      userID,
				Title:       "Test Document",
				FileName:    "test." + tc.fileType,
				FileSize:    12,
				FileType:    tc.fileType,
				FileContent: fileContent,
			}

			mockStorage.On("Upload", ctx, mock.MatchedBy(func(key string) bool {
				return strings.HasSuffix(key, tc.extension)
			}), mock.Anything, int64(12), tc.contentType).Return(nil)

			mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Document")).Return(nil)

			result, err := uc.Execute(ctx, params)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, entity.FileType(tc.fileType), result.FileType)
			assert.Equal(t, tc.contentType, result.ContentType)
			assert.Contains(t, result.ObjectKey, tc.extension)

			mockRepo.AssertExpectations(t)
			mockStorage.AssertExpectations(t)
		})
	}
}

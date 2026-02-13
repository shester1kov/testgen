package document

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestUploadUseCase_StorageUploadFailure(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	mockStorage := new(MockStorage)
	logger := zap.NewNop()

	uc := NewUploadUseCase(mockRepo, mockStorage, 10_000, logger)

	params := UploadParams{
		UserID:      uuid.New(),
		Title:       "title",
		FileName:    "file.pdf",
		FileSize:    4,
		FileType:    "pdf",
		FileContent: strings.NewReader("data"),
	}

	// Storage upload fails
	mockStorage.On("Upload", mock.Anything, mock.Anything, mock.Anything, int64(4), "application/pdf").
		Return(errors.New("storage error"))

	_, err := uc.Execute(context.Background(), params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to upload file to storage")

	// Repository should not be called if storage fails
	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestUploadUseCase_DBErrorWithRollback(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	mockStorage := new(MockStorage)
	logger := zap.NewNop()

	uc := NewUploadUseCase(mockRepo, mockStorage, 10_000, logger)

	params := UploadParams{
		UserID:      uuid.New(),
		Title:       "title",
		FileName:    "file.txt",
		FileSize:    10,
		FileType:    "txt",
		FileContent: strings.NewReader("test data"),
	}

	// Storage upload succeeds
	mockStorage.On("Upload", mock.Anything, mock.Anything, mock.Anything, int64(10), "text/plain").Return(nil)

	// DB create fails
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Document")).Return(errors.New("db error"))

	// Rollback should be called
	mockStorage.On("Delete", mock.Anything, mock.Anything).Return(nil)

	_, err := uc.Execute(context.Background(), params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save document metadata")

	// Verify rollback was attempted
	mockStorage.AssertCalled(t, "Delete", mock.Anything, mock.Anything)
}

func TestUploadUseCase_DBErrorWithRollbackFailure(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	mockStorage := new(MockStorage)
	logger := zap.NewNop()

	uc := NewUploadUseCase(mockRepo, mockStorage, 10_000, logger)

	params := UploadParams{
		UserID:      uuid.New(),
		Title:       "title",
		FileName:    "file.md",
		FileSize:    2,
		FileType:    "md",
		FileContent: strings.NewReader("ok"),
	}

	// Storage upload succeeds
	mockStorage.On("Upload", mock.Anything, mock.Anything, mock.Anything, int64(2), "text/markdown").Return(nil)

	// DB create fails
	dbErr := errors.New("database connection lost")
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Document")).Return(dbErr)

	// Rollback also fails
	mockStorage.On("Delete", mock.Anything, mock.Anything).Return(errors.New("rollback failed"))

	_, err := uc.Execute(context.Background(), params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save document metadata")

	// Even if rollback fails, the original error should be returned
	// Verify rollback was attempted
	mockStorage.AssertCalled(t, "Delete", mock.Anything, mock.Anything)
}

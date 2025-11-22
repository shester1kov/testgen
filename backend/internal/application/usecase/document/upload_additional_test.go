package document

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type failingReader struct{}

func (f failingReader) Read(p []byte) (int, error) {
	return 0, errors.New("read failed")
}

func TestUploadUseCase_MkdirAllFailure(t *testing.T) {
	mockRepo := new(MockDocumentRepository)

	tempFile, err := os.CreateTemp("", "upload-destination")
	assert.NoError(t, err)
	t.Cleanup(func() { os.Remove(tempFile.Name()) })

	uc := NewUploadUseCase(mockRepo, tempFile.Name(), 10_000)
	params := UploadParams{
		UserID:      uuid.New(),
		Title:       "title",
		FileName:    "file.pdf",
		FileSize:    4,
		FileType:    "pdf",
		FileContent: strings.NewReader("data"),
	}

	_, err = uc.Execute(context.Background(), params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create upload directory")
	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestUploadUseCase_CopyFailureCleansUp(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	tempDir := t.TempDir()
	uc := NewUploadUseCase(mockRepo, tempDir, 10_000)

	params := UploadParams{
		UserID:      uuid.New(),
		Title:       "title",
		FileName:    "file.txt",
		FileSize:    10,
		FileType:    "txt",
		FileContent: io.NopCloser(strings.NewReader("")),
	}
	// override reader with one that fails after create
	params.FileContent = failingReader{}

	_, err := uc.Execute(context.Background(), params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save file")

	entries, _ := os.ReadDir(tempDir)
	assert.Len(t, entries, 0, "temporary file should be removed on failure")
	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestUploadUseCase_PropagatesRepositoryError(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	tempDir := t.TempDir()
	uc := NewUploadUseCase(mockRepo, tempDir, 10_000)

	params := UploadParams{
		UserID:      uuid.New(),
		Title:       "title",
		FileName:    "file.md",
		FileSize:    2,
		FileType:    "md",
		FileContent: strings.NewReader("ok"),
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Document")).Return(assert.AnError)

	_, err := uc.Execute(context.Background(), params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save document metadata")

	files, _ := os.ReadDir(tempDir)
	if len(files) > 0 {
		os.Remove(filepath.Join(tempDir, files[0].Name()))
	}
}

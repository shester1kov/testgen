package document

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
)

// UploadUseCase handles document upload
type UploadUseCase struct {
	documentRepo repository.DocumentRepository
	uploadDir    string
	maxFileSize  int64
}

// NewUploadUseCase creates a new upload use case
func NewUploadUseCase(documentRepo repository.DocumentRepository, uploadDir string, maxFileSize int64) *UploadUseCase {
	return &UploadUseCase{
		documentRepo: documentRepo,
		uploadDir:    uploadDir,
		maxFileSize:  maxFileSize,
	}
}

// UploadParams contains upload parameters
type UploadParams struct {
	UserID      uuid.UUID
	Title       string
	FileName    string
	FileType    string
	FileSize    int64
	FileContent io.Reader
}

// Execute executes the upload use case
func (uc *UploadUseCase) Execute(ctx context.Context, params UploadParams) (*entity.Document, error) {
	// Validate file size
	if params.FileSize > uc.maxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", uc.maxFileSize)
	}

	// Validate file type
	validTypes := map[string]bool{"pdf": true, "docx": true, "pptx": true, "txt": true, "md": true}
	if !validTypes[params.FileType] {
		return nil, fmt.Errorf("unsupported file type: %s", params.FileType)
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uc.uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Generate unique file name
	uniqueFileName := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(params.FileName))
	filePath := filepath.Join(uc.uploadDir, uniqueFileName)

	// Save file to disk
	outFile, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, params.FileContent); err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Create document entity
	document := &entity.Document{
		ID:       uuid.New(),
		UserID:   params.UserID,
		Title:    params.Title,
		FileName: params.FileName,
		FilePath: filePath,
		FileType: entity.FileType(params.FileType),
		FileSize: params.FileSize,
		Status:   entity.StatusUploaded,
	}

	// Save document metadata to database
	if err := uc.documentRepo.Create(ctx, document); err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save document metadata: %w", err)
	}

	return document, nil
}

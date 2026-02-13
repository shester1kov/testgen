package document

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/storage"
	"go.uber.org/zap"
)

// UploadUseCase handles document upload
type UploadUseCase struct {
	documentRepo repository.DocumentRepository
	storage      storage.Storage
	maxFileSize  int64
	logger       *zap.Logger
}

// NewUploadUseCase creates a new upload use case
func NewUploadUseCase(
	documentRepo repository.DocumentRepository,
	storage storage.Storage,
	maxFileSize int64,
	logger *zap.Logger,
) *UploadUseCase {
	return &UploadUseCase{
		documentRepo: documentRepo,
		storage:      storage,
		maxFileSize:  maxFileSize,
		logger:       logger,
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
	uc.logger.Info(
		"Starting document upload",
		zap.String("user_id", params.UserID.String()),
		zap.String("filename", params.FileName),
		zap.Int64("size", params.FileSize),
	)

	// Validate file size
	if params.FileSize > uc.maxFileSize {
		uc.logger.Warn("File size exceeds limit",
			zap.Int64("size", params.FileSize),
			zap.Int64("max_size", uc.maxFileSize),
		)
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", uc.maxFileSize)
	}

	// Validate file type
	validTypes := map[string]bool{"pdf": true, "docx": true, "pptx": true, "txt": true, "md": true}
	if !validTypes[params.FileType] {
		return nil, fmt.Errorf("unsupported file type: %s", params.FileType)
	}

	// Определяем Content-Type для MinIO
	// Content-Type нужен для:
	// - Корректной отдачи файла браузеру (Content-Type header)
	// - Определения типа файла в MinIO Console
	contentType := entity.GetContentTypeByExtension(entity.FileType(params.FileType))

	// Генерируем уникальный object key для MinIO
	// Формат: "documents/{userID}/{uuid}.{ext}"
	// Пример: "documents/550e8400-e29b-41d4-a716-446655440000/abc123.pdf"
	// Преимущества:
	// - Уникальность (uuid)
	// - Изоляция по пользователям (userID в пути)
	// - Сохранение расширения для совместимости
	ext := filepath.Ext(params.FileName) // получаем расширение из оригинального файла
	objectKey := generateObjectKey(params.UserID, ext)

	// Загружаем файл в MinIO вместо локального диска
	// storage.Upload() принимает:
	// - ctx: контекст для таймаутов
	// - objectKey: уникальный ключ объекта
	// - params.FileContent: io.Reader с данными файла
	// - params.FileSize: размер для progress tracking
	// - contentType: MIME-тип для HTTP заголовков
	if err := uc.storage.Upload(ctx, objectKey, params.FileContent, params.FileSize, contentType); err != nil {
		uc.logger.Error("Failed to upload file to storage",
			zap.Error(err),
			zap.String("object_key", objectKey),
		)
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	// Create document entity
	document := &entity.Document{
		ID:          uuid.New(),
		UserID:      params.UserID,
		Title:       params.Title,
		FileName:    params.FileName,
		ObjectKey:   objectKey,
		FileType:    entity.FileType(params.FileType),
		FileSize:    params.FileSize,
		ContentType: contentType,
		Status:      entity.StatusUploaded,
	}

	// Rollback через storage.Delete() при ошибке БД
	// Если сохранение в БД провалилось, нужно удалить файл из MinIO
	// чтобы не оставлять "мусорные" файлы в хранилище
	if err := uc.documentRepo.Create(ctx, document); err != nil {
		uc.logger.Error("Failed to save document metadata to database",
			zap.Error(err),
			zap.String("object_key", objectKey),
		)

		// Пытаемся удалить файл из MinIO (rollback)
		if deleteErr := uc.storage.Delete(ctx, objectKey); deleteErr != nil {
			// Логируем ошибку удаления, но возвращаем оригинальную ошибку БД
			uc.logger.Error("Failed to rollback file from storage",
				zap.Error(deleteErr),
				zap.String("object_key", objectKey),
			)
		}

		return nil, fmt.Errorf("failed to save document metadata: %w", err)
	}

	uc.logger.Info("Document uploaded successfully",
		zap.String("document_id", document.ID.String()),
		zap.String("object_key", objectKey),
	)

	return document, nil
}

// generateObjectKey генерирует уникальный object key для MinIO
// Это бизнес-логика генерации пути, поэтому находится на уровне Use Case
// Формат: "documents/{userID}/{uuid}{ext}"
// Примеры:
//   - "documents/550e8400-e29b-41d4-a716-446655440000/abc123.pdf"
//   - "documents/550e8400-e29b-41d4-a716-446655440000/def456.docx"
func generateObjectKey(userID uuid.UUID, ext string) string {
	fileUUID := uuid.New()
	return fmt.Sprintf("documents/%s/%s%s", userID.String(), fileUUID.String(), ext)
}

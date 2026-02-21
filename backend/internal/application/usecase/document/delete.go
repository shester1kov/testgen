package document

import (
	"context"
	"fmt"
	"github.com/shester1kov/testgen-backend/internal/domain"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/storage"
	"go.uber.org/zap"
)

// DeleteUseCase handles document deletion
type DeleteUseCase struct {
	documentRepo repository.DocumentRepository
	storage      storage.Storage
	userRepo     repository.UserRepository
	logger       *zap.Logger
}

// NewDeleteUseCase creates a new delete use case
func NewDeleteUseCase(
	documentRepo repository.DocumentRepository,
	storage storage.Storage,
	userRepo repository.UserRepository,
	logger *zap.Logger,
) *DeleteUseCase {
	return &DeleteUseCase{
		documentRepo: documentRepo,
		storage:      storage,
		userRepo:     userRepo,
		logger:       logger,
	}
}

// Execute executes the delete use case
// Порядок операций:
// 1. Проверка существования документа
// 2. Проверка прав доступа (ownership)
// 3. Удаление файла из MinIO (если есть ObjectKey)
// 4. Удаление метаданных из БД (soft delete)
//
// Такой порядок гарантирует, что если удаление файла провалилось,
// метаданные остаются в БД и можно повторить попытку удаления.
func (uc *DeleteUseCase) Execute(ctx context.Context, documentID uuid.UUID, userID uuid.UUID) error {
	uc.logger.Info(
		"Starting document deletion",
		zap.String("document_id", documentID.String()),
		zap.String("user_id", userID.String()),
	)

	// Fetch document
	document, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		uc.logger.Error("Document not found",
			zap.Error(err),
			zap.String("document_id", documentID.String()),
		)
		return fmt.Errorf("%w: %v", domain.ErrDocumentNotFound, err)
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	// Verify ownership
	// Важно: проверяем, что пользователь имеет право удалить документ
	// Это защита от несанкционированного удаления чужих документов
	if document.UserID != userID && !user.IsAdmin() {
		uc.logger.Warn("Access denied: document belongs to another user",
			zap.String("document_id", documentID.String()),
			zap.String("requesting_user_id", userID.String()),
		)
		return domain.ErrAccessDenied
	}

	// Удаляем файл из MinIO ПЕРЕД удалением из БД
	// - Если удаление файла провалилось, метаданные остаются в БД можно повторить
	// - Если удаление файла успешно, но удаление из БД провалилось rollback не нужен,
	//   просто логируем и возвращаем ошибку, пользователь увидит документ и может повторить
	if document.ObjectKey != "" {
		uc.logger.Info("Deleting file from storage",
			zap.String("object_key", document.ObjectKey),
		)

		if err := uc.storage.Delete(ctx, document.ObjectKey); err != nil {
			uc.logger.Error("Failed to delete file from storage",
				zap.Error(err),
				zap.String("object_key", document.ObjectKey),
			)

			return fmt.Errorf("failed to delete file from storage: %w", err)

		}

		uc.logger.Info("File deleted from storage successfully",
			zap.String("object_key", document.ObjectKey),
		)
	} else {
		// Случай, когда ObjectKey пустой (возможно старая запись или ошибка при загрузке)
		uc.logger.Warn("Document has no object key, skipping file deletion",
			zap.String("document_id", documentID.String()),
		)
	}

	// Soft delete from database
	// repository.Delete() использует GORM soft delete (устанавливает deleted_at)
	if err := uc.documentRepo.Delete(ctx, documentID); err != nil {
		uc.logger.Error("Failed to delete document from database",
			zap.Error(err),
			zap.String("document_id", documentID.String()),
		)
		return fmt.Errorf("failed to delete document: %w", err)
	}

	uc.logger.Info("Document deleted successfully",
		zap.String("document_id", documentID.String()),
		zap.String("object_key", document.ObjectKey),
	)

	return nil
}

package document

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
)

// DeleteUseCase handles document deletion
type DeleteUseCase struct {
	documentRepo repository.DocumentRepository
}

// NewDeleteUseCase creates a new delete use case
func NewDeleteUseCase(documentRepo repository.DocumentRepository) *DeleteUseCase {
	return &DeleteUseCase{
		documentRepo: documentRepo,
	}
}

// Execute executes the delete use case
func (uc *DeleteUseCase) Execute(ctx context.Context, documentID uuid.UUID, userID uuid.UUID) error {
	// Fetch document
	document, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	// Verify ownership
	if document.UserID != userID {
		return fmt.Errorf("access denied: document belongs to another user")
	}

	// Soft delete from database
	if err := uc.documentRepo.Delete(ctx, documentID); err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	// Delete physical file (best effort, don't fail if file doesn't exist)
	if document.FilePath != "" {
		_ = os.Remove(document.FilePath)
	}

	return nil
}

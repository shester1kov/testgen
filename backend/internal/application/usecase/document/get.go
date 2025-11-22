package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
)

// GetUseCase handles retrieving a single document
type GetUseCase struct {
	documentRepo repository.DocumentRepository
}

// NewGetUseCase creates a new get use case
func NewGetUseCase(documentRepo repository.DocumentRepository) *GetUseCase {
	return &GetUseCase{
		documentRepo: documentRepo,
	}
}

// Execute executes the get use case
func (uc *GetUseCase) Execute(ctx context.Context, documentID uuid.UUID, userID uuid.UUID) (*entity.Document, error) {
	// Fetch document
	document, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	// Verify ownership
	if document.UserID != userID {
		return nil, fmt.Errorf("access denied: document belongs to another user")
	}

	return document, nil
}

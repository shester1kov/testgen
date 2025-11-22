package document

import (
	"context"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
)

// ListUseCase handles document listing
type ListUseCase struct {
	documentRepo repository.DocumentRepository
}

// NewListUseCase creates a new list use case
func NewListUseCase(documentRepo repository.DocumentRepository) *ListUseCase {
	return &ListUseCase{
		documentRepo: documentRepo,
	}
}

// ListParams contains list parameters
type ListParams struct {
	UserID   uuid.UUID
	Page     int
	PageSize int
}

// ListResult contains list results
type ListResult struct {
	Documents []*entity.Document
	Total     int64
}

// Execute executes the list use case
func (uc *ListUseCase) Execute(ctx context.Context, params ListParams) (*ListResult, error) {
	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Fetch documents
	documents, err := uc.documentRepo.FindByUserID(ctx, params.UserID, params.PageSize, offset)
	if err != nil {
		return nil, err
	}

	// Count total
	total, err := uc.documentRepo.CountByUserID(ctx, params.UserID)
	if err != nil {
		return nil, err
	}

	return &ListResult{
		Documents: documents,
		Total:     total,
	}, nil
}

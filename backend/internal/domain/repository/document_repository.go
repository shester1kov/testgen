package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

// DocumentRepository defines the interface for document data operations
type DocumentRepository interface {
	Create(ctx context.Context, document *entity.Document) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Document, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Document, error)
	FindAll(ctx context.Context, limit, offset int) ([]*entity.Document, error)
	Update(ctx context.Context, document *entity.Document) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	CountAll(ctx context.Context) (int64, error)
}

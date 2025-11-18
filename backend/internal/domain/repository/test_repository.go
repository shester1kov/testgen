package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

// TestRepository defines the interface for test data operations
type TestRepository interface {
	Create(ctx context.Context, test *entity.Test) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Test, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Test, error)
	Update(ctx context.Context, test *entity.Test) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
}

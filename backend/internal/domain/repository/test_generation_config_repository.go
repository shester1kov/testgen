package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

// TestGenerationConfigRepository defines the interface for test generation config data access.
type TestGenerationConfigRepository interface {
	// Create persists an immutable generation config record for a test.
	Create(ctx context.Context, config *entity.TestGenerationConfig) error

	// FindByTestID returns the generation config for a test, or nil if the test
	// was created manually (no config row exists).
	FindByTestID(ctx context.Context, testID uuid.UUID) (*entity.TestGenerationConfig, error)
}

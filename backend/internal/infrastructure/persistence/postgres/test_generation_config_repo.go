package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type testGenerationConfigRepository struct {
	db *gorm.DB
}

// NewTestGenerationConfigRepository creates a new test generation config repository.
func NewTestGenerationConfigRepository(db *gorm.DB) repository.TestGenerationConfigRepository {
	return &testGenerationConfigRepository{db: db}
}

// Create persists a new immutable generation config.
func (r *testGenerationConfigRepository) Create(ctx context.Context, config *entity.TestGenerationConfig) error {
	return r.db.WithContext(ctx).Create(config).Error
}

// FindByTestID returns the generation config for a test.
// Returns nil (no error) when no config row exists for the test (manually created test).
func (r *testGenerationConfigRepository) FindByTestID(ctx context.Context, testID uuid.UUID) (*entity.TestGenerationConfig, error) {
	var config entity.TestGenerationConfig
	err := r.db.WithContext(ctx).
		Preload("AcademicProfile").
		Where("test_id = ?", testID).
		First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &config, nil
}

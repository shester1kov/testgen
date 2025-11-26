package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type testRepository struct {
	db *gorm.DB
}

// NewTestRepository creates a new instance of test repository
func NewTestRepository(db *gorm.DB) repository.TestRepository {
	return &testRepository{db: db}
}

func (r *testRepository) Create(ctx context.Context, test *entity.Test) error {
	return r.db.WithContext(ctx).Create(test).Error
}

func (r *testRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Test, error) {
	var test entity.Test
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Document").
		Preload("Questions").
		Preload("Questions.Answers").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&test).Error
	if err != nil {
		return nil, err
	}
	return &test, nil
}

func (r *testRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Test, error) {
	var tests []*entity.Test
	err := r.db.WithContext(ctx).
		Preload("Document").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&tests).Error
	return tests, err
}

func (r *testRepository) Update(ctx context.Context, test *entity.Test) error {
	return r.db.WithContext(ctx).Save(test).Error
}

func (r *testRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.Test{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

func (r *testRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Test{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

func (r *testRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.Test, error) {
	var tests []*entity.Test
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Document").
		Where("deleted_at IS NULL").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&tests).Error
	return tests, err
}

func (r *testRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Test{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	return count, err
}

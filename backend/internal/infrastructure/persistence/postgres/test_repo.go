package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/cache"
	"gorm.io/gorm"
)

type testRepository struct {
	db    *gorm.DB
	cache *cache.RedisClient
}

// NewTestRepository creates a new instance of test repository
func NewTestRepository(db *gorm.DB, cache *cache.RedisClient) repository.TestRepository {
	return &testRepository{
		db:    db,
		cache: cache,
	}
}

func (r *testRepository) Create(ctx context.Context, test *entity.Test) error {
	if err := r.db.WithContext(ctx).Create(test).Error; err != nil {
		return err
	}

	// Cache the newly created test
	_ = r.cache.Set(ctx, cache.TestByIDKey(test.ID), test, cache.TestTTL)
	// Invalidate user's tests list cache
	_ = r.cache.Delete(ctx, cache.TestsByUserKey(test.UserID))

	return nil
}

func (r *testRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Test, error) {
	// Try to get from cache first
	cacheKey := cache.TestByIDKey(id)
	var test entity.Test
	err := r.cache.Get(ctx, cacheKey, &test)
	if err == nil {
		return &test, nil // Cache hit
	}

	// Cache miss or error - fetch from database
	if !errors.Is(err, cache.ErrCacheMiss) {
		// Log cache error but continue to database
	}

	err = r.db.WithContext(ctx).
		Preload("User").
		Preload("Document").
		Preload("Questions").
		Preload("Questions.Answers").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&test).Error
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests
	_ = r.cache.Set(ctx, cacheKey, &test, cache.TestTTL)

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
	if err := r.db.WithContext(ctx).Save(test).Error; err != nil {
		return err
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.TestByIDKey(test.ID))
	_ = r.cache.Delete(ctx, cache.TestsByUserKey(test.UserID))
	_ = r.cache.Delete(ctx, cache.QuestionsByTestKey(test.ID))

	return nil
}

func (r *testRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Get test for cache invalidation
	var test entity.Test
	if err := r.db.WithContext(ctx).Select("user_id").Where("id = ?", id).First(&test).Error; err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).
		Model(&entity.Test{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error; err != nil {
		return err
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.TestByIDKey(id))
	_ = r.cache.Delete(ctx, cache.TestsByUserKey(test.UserID))
	_ = r.cache.Delete(ctx, cache.QuestionsByTestKey(id))

	return nil
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

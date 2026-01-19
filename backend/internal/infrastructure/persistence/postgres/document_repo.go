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

type documentRepository struct {
	db    *gorm.DB
	cache *cache.RedisClient
}

// NewDocumentRepository creates a new instance of document repository
func NewDocumentRepository(db *gorm.DB, cache *cache.RedisClient) repository.DocumentRepository {
	return &documentRepository{
		db:    db,
		cache: cache,
	}
}

func (r *documentRepository) Create(ctx context.Context, document *entity.Document) error {
	if err := r.db.WithContext(ctx).Create(document).Error; err != nil {
		return err
	}

	// Cache the newly created document
	_ = r.cache.Set(ctx, cache.DocumentByIDKey(document.ID), document, cache.DocumentTTL)
	// Invalidate user's documents list cache
	_ = r.cache.Delete(ctx, cache.DocumentsByUserKey(document.UserID))

	return nil
}

func (r *documentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Document, error) {
	// Try to get from cache first
	cacheKey := cache.DocumentByIDKey(id)
	var document entity.Document
	err := r.cache.Get(ctx, cacheKey, &document)
	if err == nil {
		return &document, nil // Cache hit
	}

	// Cache miss or error - fetch from database
	if !errors.Is(err, cache.ErrCacheMiss) {
		// Log cache error but continue to database
	}

	err = r.db.WithContext(ctx).
		Preload("User").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&document).Error
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests
	_ = r.cache.Set(ctx, cacheKey, &document, cache.DocumentTTL)

	return &document, nil
}

func (r *documentRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Document, error) {
	var documents []*entity.Document
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&documents).Error
	return documents, err
}

func (r *documentRepository) Update(ctx context.Context, document *entity.Document) error {
	if err := r.db.WithContext(ctx).Save(document).Error; err != nil {
		return err
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.DocumentByIDKey(document.ID))
	_ = r.cache.Delete(ctx, cache.DocumentsByUserKey(document.UserID))

	return nil
}

func (r *documentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Get document for cache invalidation
	var doc entity.Document
	if err := r.db.WithContext(ctx).Select("user_id").Where("id = ?", id).First(&doc).Error; err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).
		Model(&entity.Document{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error; err != nil {
		return err
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.DocumentByIDKey(id))
	_ = r.cache.Delete(ctx, cache.DocumentsByUserKey(doc.UserID))

	return nil
}

func (r *documentRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Document{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

func (r *documentRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.Document, error) {
	var documents []*entity.Document
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("deleted_at IS NULL").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&documents).Error
	return documents, err
}

func (r *documentRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Document{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	return count, err
}

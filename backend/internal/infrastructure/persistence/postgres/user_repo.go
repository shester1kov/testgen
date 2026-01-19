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

type userRepository struct {
	db    *gorm.DB
	cache *cache.RedisClient
}

// NewUserRepository creates a new instance of user repository
func NewUserRepository(db *gorm.DB, cache *cache.RedisClient) repository.UserRepository {
	return &userRepository{
		db:    db,
		cache: cache,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}

	// Cache the newly created user
	_ = r.cache.Set(ctx, cache.UserByIDKey(user.ID), user, cache.UserTTL)
	_ = r.cache.Set(ctx, cache.UserByEmailKey(user.Email), user, cache.UserTTL)

	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	// Try to get from cache first
	cacheKey := cache.UserByIDKey(id)
	var user entity.User
	err := r.cache.Get(ctx, cacheKey, &user)
	if err == nil {
		return &user, nil // Cache hit
	}

	// Cache miss or error - fetch from database
	if !errors.Is(err, cache.ErrCacheMiss) {
		// Log cache error but continue to database
	}

	err = r.db.WithContext(ctx).Preload("Role").Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests
	_ = r.cache.Set(ctx, cacheKey, &user, cache.UserTTL)

	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	// Try to get from cache first
	cacheKey := cache.UserByEmailKey(email)
	var user entity.User
	err := r.cache.Get(ctx, cacheKey, &user)
	if err == nil {
		return &user, nil // Cache hit
	}

	// Cache miss or error - fetch from database
	if !errors.Is(err, cache.ErrCacheMiss) {
		// Log cache error but continue to database
	}

	err = r.db.WithContext(ctx).Preload("Role").Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests (both by ID and email)
	_ = r.cache.Set(ctx, cacheKey, &user, cache.UserTTL)
	_ = r.cache.Set(ctx, cache.UserByIDKey(user.ID), &user, cache.UserTTL)

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	// Get old email before update for cache invalidation
	var oldUser entity.User
	if err := r.db.WithContext(ctx).Select("email").Where("id = ?", user.ID).First(&oldUser).Error; err != nil {
		return err
	}

	// Use Model().Select() to explicitly update role_id
	// Save() doesn't work well when Role association is preloaded
	if err := r.db.WithContext(ctx).Model(user).Select("email", "password_hash", "full_name", "role_id", "updated_at").Updates(user).Error; err != nil {
		return err
	}

	// Invalidate cache (both old and new email if changed)
	_ = r.cache.Delete(ctx, cache.UserByIDKey(user.ID))
	_ = r.cache.Delete(ctx, cache.UserByEmailKey(oldUser.Email))
	if oldUser.Email != user.Email {
		_ = r.cache.Delete(ctx, cache.UserByEmailKey(user.Email))
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Get user email for cache invalidation
	var user entity.User
	if err := r.db.WithContext(ctx).Select("email").Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error; err != nil {
		return err
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.UserByIDKey(id), cache.UserByEmailKey(user.Email))

	return nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("deleted_at IS NULL").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).Error
	return users, err
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("deleted_at IS NULL").Count(&count).Error
	return count, err
}

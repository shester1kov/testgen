package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type documentRepository struct {
	db *gorm.DB
}

// NewDocumentRepository creates a new instance of document repository
func NewDocumentRepository(db *gorm.DB) repository.DocumentRepository {
	return &documentRepository{db: db}
}

func (r *documentRepository) Create(ctx context.Context, document *entity.Document) error {
	return r.db.WithContext(ctx).Create(document).Error
}

func (r *documentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Document, error) {
	var document entity.Document
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&document).Error
	if err != nil {
		return nil, err
	}
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
	return r.db.WithContext(ctx).Save(document).Error
}

func (r *documentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.Document{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

func (r *documentRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Document{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

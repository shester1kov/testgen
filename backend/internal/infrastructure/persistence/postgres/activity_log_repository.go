package postgres

import (
	"context"

	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type activityLogRepository struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) repository.ActivityLogRepository {
	return &activityLogRepository{db: db}
}

func (r *activityLogRepository) Create(ctx context.Context, log *entity.ActivityLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

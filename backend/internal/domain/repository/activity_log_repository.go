package repository

import (
	"context"

	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

type ActivityLogRepository interface {
	Create(ctx context.Context, log *entity.ActivityLog) error
}

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

// RoleRepository defines the interface for role data access
type RoleRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)
	FindByName(ctx context.Context, name entity.RoleName) (*entity.Role, error)
	List(ctx context.Context) ([]*entity.Role, error)
}

package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	var role entity.Role
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindByName(ctx context.Context, name entity.RoleName) (*entity.Role, error) {
	var role entity.Role
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) List(ctx context.Context) ([]*entity.Role, error) {
	var roles []*entity.Role
	if err := r.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

package repository

import (
	"context"

	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

// AcademicProfileRepository defines the interface for academic profile data access.
type AcademicProfileRepository interface {
	// FindByCode returns a profile with its formulation examples and question examples
	// pre-loaded, looked up by the machine-readable code (e.g. "technological").
	FindByCode(ctx context.Context, code string) (*entity.AcademicProfile, error)

	// List returns all profiles ordered by sort_order, without child relations.
	List(ctx context.Context) ([]*entity.AcademicProfile, error)
}

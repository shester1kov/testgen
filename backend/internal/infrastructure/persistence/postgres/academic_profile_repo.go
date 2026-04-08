package postgres

import (
	"context"

	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type academicProfileRepository struct {
	db *gorm.DB
}

// NewAcademicProfileRepository creates a new academic profile repository.
func NewAcademicProfileRepository(db *gorm.DB) repository.AcademicProfileRepository {
	return &academicProfileRepository{db: db}
}

// FindByCode returns the profile with its formulation examples and question examples.
func (r *academicProfileRepository) FindByCode(ctx context.Context, code string) (*entity.AcademicProfile, error) {
	var profile entity.AcademicProfile
	err := r.db.WithContext(ctx).
		Preload("Formulations", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_num ASC")
		}).
		Preload("QuestionExamples", func(db *gorm.DB) *gorm.DB {
			return db.Order("question_type ASC, order_num ASC")
		}).
		Where("code = ?", code).
		First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// List returns all profiles ordered by sort_order without child relations.
func (r *academicProfileRepository) List(ctx context.Context) ([]*entity.AcademicProfile, error) {
	var profiles []*entity.AcademicProfile
	if err := r.db.WithContext(ctx).Order("sort_order ASC").Find(&profiles).Error; err != nil {
		return nil, err
	}
	return profiles, nil
}

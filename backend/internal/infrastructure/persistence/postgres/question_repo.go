package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type questionRepository struct {
	db *gorm.DB
}

// NewQuestionRepository creates a new question repository
func NewQuestionRepository(db *gorm.DB) repository.QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) Create(ctx context.Context, question *entity.Question) error {
	return r.db.WithContext(ctx).Create(question).Error
}

func (r *questionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Question, error) {
	var question entity.Question
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) FindByTestID(ctx context.Context, testID uuid.UUID) ([]*entity.Question, error) {
	var questions []*entity.Question
	err := r.db.WithContext(ctx).
		Where("test_id = ?", testID).
		Order("order_num ASC").
		Find(&questions).Error
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *questionRepository) Update(ctx context.Context, question *entity.Question) error {
	return r.db.WithContext(ctx).Save(question).Error
}

func (r *questionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Question{}, "id = ?", id).Error
}

func (r *questionRepository) CountByTestID(ctx context.Context, testID uuid.UUID) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Question{}).
		Where("test_id = ?", testID).
		Count(&count).Error
	return int(count), err
}

func (r *questionRepository) ReorderQuestions(ctx context.Context, testID uuid.UUID, questionIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, qid := range questionIDs {
			if err := tx.Model(&entity.Question{}).
				Where("id = ? AND test_id = ?", qid, testID).
				Update("order_num", i+1).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

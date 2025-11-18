package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type answerRepository struct {
	db *gorm.DB
}

// NewAnswerRepository creates a new answer repository
func NewAnswerRepository(db *gorm.DB) repository.AnswerRepository {
	return &answerRepository{db: db}
}

func (r *answerRepository) Create(ctx context.Context, answer *entity.Answer) error {
	return r.db.WithContext(ctx).Create(answer).Error
}

func (r *answerRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Answer, error) {
	var answer entity.Answer
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&answer).Error
	if err != nil {
		return nil, err
	}
	return &answer, nil
}

func (r *answerRepository) FindByQuestionID(ctx context.Context, questionID uuid.UUID) ([]*entity.Answer, error) {
	var answers []*entity.Answer
	err := r.db.WithContext(ctx).
		Where("question_id = ?", questionID).
		Order("order_num ASC").
		Find(&answers).Error
	if err != nil {
		return nil, err
	}
	return answers, nil
}

func (r *answerRepository) Update(ctx context.Context, answer *entity.Answer) error {
	return r.db.WithContext(ctx).Save(answer).Error
}

func (r *answerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Answer{}, "id = ?", id).Error
}

func (r *answerRepository) DeleteByQuestionID(ctx context.Context, questionID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Answer{}, "question_id = ?", questionID).Error
}

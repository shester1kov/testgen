package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

// QuestionRepository defines the interface for question data operations
type QuestionRepository interface {
	// Create creates a new question
	Create(ctx context.Context, question *entity.Question) error

	// FindByID retrieves a question by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Question, error)

	// FindByTestID retrieves all questions for a specific test
	FindByTestID(ctx context.Context, testID uuid.UUID) ([]*entity.Question, error)

	// Update updates an existing question
	Update(ctx context.Context, question *entity.Question) error

	// Delete deletes a question by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// CountByTestID counts questions for a specific test
	CountByTestID(ctx context.Context, testID uuid.UUID) (int, error)

	// ReorderQuestions updates the order of questions
	ReorderQuestions(ctx context.Context, testID uuid.UUID, questionIDs []uuid.UUID) error
}

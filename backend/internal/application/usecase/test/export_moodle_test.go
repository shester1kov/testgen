package test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/moodle"
	"github.com/stretchr/testify/require"
)

type mockTestRepository struct {
	repository.TestRepository
	findByID func(ctx context.Context, id uuid.UUID) (*entity.Test, error)
}

func (m *mockTestRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Test, error) {
	return m.findByID(ctx, id)
}

type mockQuestionRepository struct {
	repository.QuestionRepository
	findByTestID func(ctx context.Context, testID uuid.UUID) ([]*entity.Question, error)
}

func (m *mockQuestionRepository) FindByTestID(ctx context.Context, testID uuid.UUID) ([]*entity.Question, error) {
	return m.findByTestID(ctx, testID)
}

type mockAnswerRepository struct {
	repository.AnswerRepository
	findByQuestionID func(ctx context.Context, questionID uuid.UUID) ([]*entity.Answer, error)
}

func (m *mockAnswerRepository) FindByQuestionID(ctx context.Context, questionID uuid.UUID) ([]*entity.Answer, error) {
	return m.findByQuestionID(ctx, questionID)
}

func TestExportMoodleUseCase(t *testing.T) {
	userID := uuid.New()
	testID := uuid.New()
	questionID := uuid.New()

	sampleTest := &entity.Test{ID: testID, UserID: userID}
	questions := []*entity.Question{{ID: questionID, QuestionText: "Q1", QuestionType: entity.QuestionTypeSingleChoice, Points: 1, OrderNum: 1}}
	answers := []*entity.Answer{{ID: uuid.New(), QuestionID: questionID, AnswerText: "A1", IsCorrect: true, OrderNum: 1}}

	testCases := []struct {
		name        string
		setup       func() *ExportMoodleUseCase
		expectedErr string
	}{
		{
			name: "exports xml successfully",
			setup: func() *ExportMoodleUseCase {
				return &ExportMoodleUseCase{
					testRepo: &mockTestRepository{findByID: func(ctx context.Context, id uuid.UUID) (*entity.Test, error) {
						return sampleTest, nil
					}},
					questionRepo: &mockQuestionRepository{findByTestID: func(ctx context.Context, testID uuid.UUID) ([]*entity.Question, error) {
						return questions, nil
					}},
					answerRepo: &mockAnswerRepository{findByQuestionID: func(ctx context.Context, questionID uuid.UUID) ([]*entity.Answer, error) {
						return answers, nil
					}},
					xmlExporter: moodle.NewMoodleXMLExporter(),
				}
			},
		},
		{
			name: "returns error when test not found",
			setup: func() *ExportMoodleUseCase {
				return &ExportMoodleUseCase{
					testRepo: &mockTestRepository{findByID: func(ctx context.Context, id uuid.UUID) (*entity.Test, error) {
						return nil, errors.New("missing")
					}},
					questionRepo: &mockQuestionRepository{findByTestID: func(context.Context, uuid.UUID) ([]*entity.Question, error) { return nil, nil }},
					answerRepo:   &mockAnswerRepository{findByQuestionID: func(context.Context, uuid.UUID) ([]*entity.Answer, error) { return nil, nil }},
					xmlExporter:  moodle.NewMoodleXMLExporter(),
				}
			},
			expectedErr: "test not found",
		},
		{
			name: "prevents access for other user",
			setup: func() *ExportMoodleUseCase {
				wrongOwner := *sampleTest
				wrongOwner.UserID = uuid.New()
				return &ExportMoodleUseCase{
					testRepo: &mockTestRepository{findByID: func(context.Context, uuid.UUID) (*entity.Test, error) {
						return &wrongOwner, nil
					}},
					questionRepo: &mockQuestionRepository{findByTestID: func(context.Context, uuid.UUID) ([]*entity.Question, error) { return questions, nil }},
					answerRepo:   &mockAnswerRepository{findByQuestionID: func(context.Context, uuid.UUID) ([]*entity.Answer, error) { return answers, nil }},
					xmlExporter:  moodle.NewMoodleXMLExporter(),
				}
			},
			expectedErr: "unauthorized",
		},
		{
			name: "errors when questions lookup fails",
			setup: func() *ExportMoodleUseCase {
				return &ExportMoodleUseCase{
					testRepo: &mockTestRepository{findByID: func(context.Context, uuid.UUID) (*entity.Test, error) { return sampleTest, nil }},
					questionRepo: &mockQuestionRepository{findByTestID: func(context.Context, uuid.UUID) ([]*entity.Question, error) {
						return nil, errors.New("boom")
					}},
					answerRepo:  &mockAnswerRepository{findByQuestionID: func(context.Context, uuid.UUID) ([]*entity.Answer, error) { return nil, nil }},
					xmlExporter: moodle.NewMoodleXMLExporter(),
				}
			},
			expectedErr: "failed to retrieve questions",
		},
		{
			name: "errors when no questions present",
			setup: func() *ExportMoodleUseCase {
				return &ExportMoodleUseCase{
					testRepo: &mockTestRepository{findByID: func(context.Context, uuid.UUID) (*entity.Test, error) { return sampleTest, nil }},
					questionRepo: &mockQuestionRepository{findByTestID: func(context.Context, uuid.UUID) ([]*entity.Question, error) {
						return []*entity.Question{}, nil
					}},
					answerRepo:  &mockAnswerRepository{findByQuestionID: func(context.Context, uuid.UUID) ([]*entity.Answer, error) { return nil, nil }},
					xmlExporter: moodle.NewMoodleXMLExporter(),
				}
			},
			expectedErr: "test has no questions",
		},
		{
			name: "errors when answer retrieval fails",
			setup: func() *ExportMoodleUseCase {
				return &ExportMoodleUseCase{
					testRepo:     &mockTestRepository{findByID: func(context.Context, uuid.UUID) (*entity.Test, error) { return sampleTest, nil }},
					questionRepo: &mockQuestionRepository{findByTestID: func(context.Context, uuid.UUID) ([]*entity.Question, error) { return questions, nil }},
					answerRepo: &mockAnswerRepository{findByQuestionID: func(context.Context, uuid.UUID) ([]*entity.Answer, error) {
						return nil, errors.New("answers error")
					}},
					xmlExporter: moodle.NewMoodleXMLExporter(),
				}
			},
			expectedErr: "failed to retrieve answers",
		},
		{
			name: "errors when exporter fails",
			setup: func() *ExportMoodleUseCase {
				badQuestion := []*entity.Question{{ID: questionID, QuestionText: "Q1", QuestionType: entity.QuestionType(""), Points: 1}}
				return &ExportMoodleUseCase{
					testRepo:     &mockTestRepository{findByID: func(context.Context, uuid.UUID) (*entity.Test, error) { return sampleTest, nil }},
					questionRepo: &mockQuestionRepository{findByTestID: func(context.Context, uuid.UUID) ([]*entity.Question, error) { return badQuestion, nil }},
					answerRepo:   &mockAnswerRepository{findByQuestionID: func(context.Context, uuid.UUID) ([]*entity.Answer, error) { return answers, nil }},
					xmlExporter:  moodle.NewMoodleXMLExporter(),
				}
			},
			expectedErr: "failed to convert question",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			uc := tt.setup()
			content, err := uc.Execute(context.Background(), testID, userID)

			if tt.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			require.NoError(t, err)
			require.Contains(t, content, "<quiz>")
		})
	}
}

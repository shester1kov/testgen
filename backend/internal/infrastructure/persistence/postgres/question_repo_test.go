package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupQuestionTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{SkipDefaultTransaction: true})
	require.NoError(t, err)

	err = db.Exec(`
                CREATE TABLE tests (
                        id TEXT PRIMARY KEY,
                        user_id TEXT,
                        document_id TEXT,
                        title TEXT,
                        status TEXT,
                        deleted_at DATETIME,
                        created_at DATETIME,
                        updated_at DATETIME
                );
        `).Error
	require.NoError(t, err)

	err = db.Exec(`
                CREATE TABLE questions (
                        id TEXT PRIMARY KEY,
                        test_id TEXT NOT NULL,
                        question_text TEXT NOT NULL,
                        question_type TEXT,
                        difficulty TEXT,
                        points REAL,
                        order_num INTEGER,
                        created_at DATETIME,
                        updated_at DATETIME
                );
        `).Error
	require.NoError(t, err)

	return db
}

func seedQuestions(t *testing.T, db *gorm.DB, testID uuid.UUID) []*entity.Question {
	baseTime := time.Time{}
	questions := []*entity.Question{
		{
			ID:           uuid.New(),
			TestID:       testID,
			QuestionText: "Q1",
			QuestionType: entity.QuestionTypeSingleChoice,
			Difficulty:   entity.DifficultyEasy,
			Points:       1,
			OrderNum:     2,
			CreatedAt:    baseTime,
			UpdatedAt:    baseTime,
		},
		{
			ID:           uuid.New(),
			TestID:       testID,
			QuestionText: "Q2",
			QuestionType: entity.QuestionTypeMultipleChoice,
			Difficulty:   entity.DifficultyMedium,
			Points:       2,
			OrderNum:     1,
			CreatedAt:    baseTime,
			UpdatedAt:    baseTime,
		},
	}

	for _, q := range questions {
		require.NoError(t, db.Create(q).Error)
	}
	return questions
}

func TestQuestionRepository_CreateAndFind(t *testing.T) {
	db := setupQuestionTestDB(t)
	repo := NewQuestionRepository(db)
	testID := uuid.New()

	q := &entity.Question{
		ID:           uuid.New(),
		TestID:       testID,
		QuestionText: "What is Go?",
		QuestionType: entity.QuestionTypeShortAnswer,
		Difficulty:   entity.DifficultyHard,
		Points:       3,
		OrderNum:     1,
	}

	err := repo.Create(context.Background(), q)
	require.NoError(t, err)

	fetched, err := repo.FindByID(context.Background(), q.ID)
	assert.NoError(t, err)
	assert.Equal(t, q.QuestionText, fetched.QuestionText)

	missingID := uuid.New()
	missing, err := repo.FindByID(context.Background(), missingID)
	assert.Nil(t, missing)
	assert.Error(t, err)
}

func TestQuestionRepository_FindAndCountByTest(t *testing.T) {
	db := setupQuestionTestDB(t)
	repo := NewQuestionRepository(db)
	testID := uuid.New()
	seedQuestions(t, db, testID)

	questions, err := repo.FindByTestID(context.Background(), testID)
	require.NoError(t, err)
	require.Len(t, questions, 2)
	assert.Equal(t, "Q2", questions[0].QuestionText)
	assert.Equal(t, "Q1", questions[1].QuestionText)

	count, err := repo.CountByTestID(context.Background(), testID)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	// drop table to force error branch
	require.NoError(t, db.Exec("DROP TABLE questions;").Error)
	_, err = repo.CountByTestID(context.Background(), testID)
	assert.Error(t, err)
}

func TestQuestionRepository_ReorderQuestions(t *testing.T) {
	db := setupQuestionTestDB(t)
	repo := NewQuestionRepository(db)
	testID := uuid.New()
	questions := seedQuestions(t, db, testID)

	newOrder := []uuid.UUID{questions[0].ID, questions[1].ID}
	err := repo.ReorderQuestions(context.Background(), testID, newOrder)
	assert.NoError(t, err)

	var reordered []*entity.Question
	require.NoError(t, db.Order("order_num asc").Find(&reordered).Error)
	assert.Equal(t, 1, reordered[0].OrderNum)
	assert.Equal(t, questions[0].ID, reordered[0].ID)

	// cause error by dropping table inside transaction
	require.NoError(t, db.Exec("DROP TABLE questions;").Error)
	err = repo.ReorderQuestions(context.Background(), testID, newOrder)
	assert.Error(t, err)
}

func TestQuestionRepository_Update(t *testing.T) {
	db := setupQuestionTestDB(t)
	repo := NewQuestionRepository(db)
	testID := uuid.New()

	q := &entity.Question{
		ID:           uuid.New(),
		TestID:       testID,
		QuestionText: "Original Question",
		QuestionType: entity.QuestionTypeSingleChoice,
		Difficulty:   entity.DifficultyEasy,
		Points:       1,
		OrderNum:     1,
	}

	err := repo.Create(context.Background(), q)
	require.NoError(t, err)

	// Update the question
	q.QuestionText = "Updated Question"
	q.Difficulty = entity.DifficultyHard
	q.Points = 5

	err = repo.Update(context.Background(), q)
	require.NoError(t, err)

	// Verify update
	fetched, err := repo.FindByID(context.Background(), q.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Question", fetched.QuestionText)
	assert.Equal(t, entity.DifficultyHard, fetched.Difficulty)
	assert.Equal(t, float64(5), fetched.Points)
}

func TestQuestionRepository_Delete(t *testing.T) {
	db := setupQuestionTestDB(t)
	repo := NewQuestionRepository(db)
	testID := uuid.New()

	q := &entity.Question{
		ID:           uuid.New(),
		TestID:       testID,
		QuestionText: "Question to delete",
		QuestionType: entity.QuestionTypeTrueFalse,
		Difficulty:   entity.DifficultyMedium,
		Points:       2,
		OrderNum:     1,
	}

	err := repo.Create(context.Background(), q)
	require.NoError(t, err)

	// Delete the question
	err = repo.Delete(context.Background(), q.ID)
	require.NoError(t, err)

	// Verify deletion
	fetched, err := repo.FindByID(context.Background(), q.ID)
	assert.Nil(t, fetched)
	assert.Error(t, err)

	// Count should be 0
	count, err := repo.CountByTestID(context.Background(), testID)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

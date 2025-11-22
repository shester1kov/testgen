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

func setupAnswerTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{SkipDefaultTransaction: true})
	require.NoError(t, err)

	err = db.Exec(`
                CREATE TABLE answers (
                        id TEXT PRIMARY KEY,
                        question_id TEXT NOT NULL,
                        answer_text TEXT NOT NULL,
                        is_correct BOOLEAN,
                        order_num INTEGER,
                        created_at DATETIME
                );
        `).Error
	require.NoError(t, err)

	return db
}

func TestAnswerRepository_CreateFindAndDelete(t *testing.T) {
	db := setupAnswerTestDB(t)
	repo := NewAnswerRepository(db)
	questionID := uuid.New()

	answer := &entity.Answer{
		ID:         uuid.New(),
		QuestionID: questionID,
		AnswerText: "42",
		IsCorrect:  true,
		OrderNum:   1,
		CreatedAt:  time.Time{},
	}

	err := repo.Create(context.Background(), answer)
	require.NoError(t, err)

	fetched, err := repo.FindByID(context.Background(), answer.ID)
	assert.NoError(t, err)
	assert.Equal(t, "42", fetched.AnswerText)

	list, err := repo.FindByQuestionID(context.Background(), questionID)
	assert.NoError(t, err)
	require.Len(t, list, 1)
	assert.True(t, list[0].IsCorrect)

	err = repo.Delete(context.Background(), answer.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(context.Background(), answer.ID)
	assert.Error(t, err)
}

func TestAnswerRepository_UpdateAndDeleteByQuestion(t *testing.T) {
	db := setupAnswerTestDB(t)
	repo := NewAnswerRepository(db)
	questionID := uuid.New()

	answer := &entity.Answer{
		ID:         uuid.New(),
		QuestionID: questionID,
		AnswerText: "Old",
		IsCorrect:  false,
		OrderNum:   1,
		CreatedAt:  time.Time{},
	}
	require.NoError(t, db.Create(answer).Error)

	answer.MarkCorrect()
	answer.AnswerText = "New"
	err := repo.Update(context.Background(), answer)
	assert.NoError(t, err)

	fetched, err := repo.FindByID(context.Background(), answer.ID)
	assert.NoError(t, err)
	assert.Equal(t, "New", fetched.AnswerText)
	assert.True(t, fetched.IsCorrect)

	err = repo.DeleteByQuestionID(context.Background(), questionID)
	assert.NoError(t, err)

	// force error path by dropping table
	require.NoError(t, db.Exec("DROP TABLE answers;").Error)
	err = repo.DeleteByQuestionID(context.Background(), questionID)
	assert.Error(t, err)
}

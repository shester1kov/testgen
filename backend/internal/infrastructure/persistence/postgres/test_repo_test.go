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

func setupTestRepoDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{SkipDefaultTransaction: true})
	require.NoError(t, err)

	err = db.Exec(`
                CREATE TABLE users (
                        id TEXT PRIMARY KEY,
                        email TEXT,
                        password_hash TEXT,
                        full_name TEXT,
                        role_id TEXT,
                        created_at DATETIME,
                        updated_at DATETIME
                );
        `).Error
	require.NoError(t, err)

	err = db.Exec(`
                CREATE TABLE tests (
                        id TEXT PRIMARY KEY,
                        user_id TEXT NOT NULL,
                        document_id TEXT,
                        title TEXT,
                        description TEXT,
                        total_questions INTEGER,
                        status TEXT,
                        moodle_synced BOOLEAN,
                        moodle_test_id TEXT,
                        created_at DATETIME,
                        updated_at DATETIME,
                        deleted_at DATETIME
                );
        `).Error
	require.NoError(t, err)

	err = db.Exec(`
                CREATE TABLE documents (
                        id TEXT PRIMARY KEY,
                        user_id TEXT,
                        title TEXT
                );
        `).Error
	require.NoError(t, err)

	err = db.Exec(`
                CREATE TABLE questions (
                        id TEXT PRIMARY KEY,
                        test_id TEXT,
                        question_text TEXT,
                        order_num INTEGER
                );
        `).Error
	require.NoError(t, err)

	err = db.Exec(`
                CREATE TABLE answers (
                        id TEXT PRIMARY KEY,
                        question_id TEXT,
                        answer_text TEXT,
                        is_correct BOOLEAN,
                        order_num INTEGER
                );
        `).Error
	require.NoError(t, err)

	return db
}

func seedTest(t *testing.T, db *gorm.DB, includeRelations bool) *entity.Test {
	testID := uuid.New()
	docID := uuid.New()
	now := time.Time{}
	userID := uuid.New()

	require.NoError(t, db.Exec("INSERT INTO users (id, email, password_hash, full_name, role_id) VALUES (?, ?, 'hash', 'User', ?)", userID.String(), "tester@example.com", uuid.New().String()).Error)

	if includeRelations {
		require.NoError(t, db.Exec("INSERT INTO documents (id, user_id, title) VALUES (?, ?, ?)", docID.String(), userID.String(), "Doc").Error)
		require.NoError(t, db.Exec("INSERT INTO questions (id, test_id, question_text, order_num) VALUES (?, ?, ?, 1)", uuid.New().String(), testID.String(), "Q1").Error)
		require.NoError(t, db.Exec("INSERT INTO answers (id, question_id, answer_text, is_correct, order_num) VALUES (?, ?, ?, 1, 1)", uuid.New().String(), uuid.New().String(), "A1").Error)
	}

	test := &entity.Test{
		ID:             testID,
		UserID:         userID,
		DocumentID:     &docID,
		Title:          "Generated Test",
		Description:    "desc",
		TotalQuestions: 1,
		Status:         entity.TestStatusDraft,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	require.NoError(t, db.Create(test).Error)
	return test
}

func TestTestRepository_CreateAndFetch(t *testing.T) {
	db := setupTestRepoDB(t)
	repo := NewTestRepository(db)

	test := seedTest(t, db, true)

	fetched, err := repo.FindByID(context.Background(), test.ID)
	assert.NoError(t, err)
	assert.Equal(t, test.Title, fetched.Title)
	assert.NotNil(t, fetched.Document)

	list, err := repo.FindByUserID(context.Background(), test.UserID, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	test.Title = "Updated"
	assert.NoError(t, repo.Update(context.Background(), test))

	err = repo.Delete(context.Background(), test.ID)
	assert.NoError(t, err)

	count, err := repo.CountByUserID(context.Background(), test.UserID)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, count)
}

func TestTestRepository_ErrorBranches(t *testing.T) {
	db := setupTestRepoDB(t)
	repo := NewTestRepository(db)
	test := seedTest(t, db, false)

	require.NoError(t, db.Exec("DROP TABLE tests;").Error)

	_, err := repo.FindByID(context.Background(), test.ID)
	assert.Error(t, err)

	_, err = repo.FindByUserID(context.Background(), test.UserID, 5, 0)
	assert.Error(t, err)
}

func TestTestRepository_Create(t *testing.T) {
	db := setupTestRepoDB(t)
	repo := NewTestRepository(db)

	userID := uuid.New()
	require.NoError(t, db.Exec("INSERT INTO users (id, email, password_hash, full_name, role_id) VALUES (?, ?, 'hash', 'User', ?)", userID.String(), "test@example.com", uuid.New().String()).Error)

	newTest := &entity.Test{
		ID:             uuid.New(),
		UserID:         userID,
		Title:          "New Test",
		Description:    "Test description",
		TotalQuestions: 5,
		Status:         entity.TestStatusDraft,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := repo.Create(context.Background(), newTest)
	require.NoError(t, err)

	// Verify it was created
	fetched, err := repo.FindByID(context.Background(), newTest.ID)
	require.NoError(t, err)
	assert.Equal(t, newTest.Title, fetched.Title)
	assert.Equal(t, newTest.Description, fetched.Description)
	assert.Equal(t, newTest.TotalQuestions, fetched.TotalQuestions)
}

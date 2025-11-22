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

func setupDocumentTestDB(t *testing.T) *gorm.DB {
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
                CREATE TABLE documents (
                        id TEXT PRIMARY KEY,
                        user_id TEXT NOT NULL,
                        title TEXT NOT NULL,
                        file_name TEXT,
                        file_path TEXT,
                        file_type TEXT,
                        file_size INTEGER,
                        parsed_text TEXT,
                        status TEXT,
                        error_msg TEXT,
                        created_at DATETIME,
                        updated_at DATETIME,
                        deleted_at DATETIME
                );
        `).Error
	require.NoError(t, err)

	return db
}

func createDocument(t *testing.T, db *gorm.DB, userID uuid.UUID) *entity.Document {
	require.NoError(t, db.Exec(
		"INSERT INTO users (id, email, password_hash, full_name, role_id) VALUES (?, ?, 'hash', 'User', ?)",
		userID.String(),
		"user@example.com",
		uuid.New().String(),
	).Error)

	doc := &entity.Document{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     "Sample",
		FileName:  "file.txt",
		FilePath:  "/tmp/file.txt",
		FileType:  entity.FileTypeTXT,
		FileSize:  10,
		Status:    entity.StatusUploaded,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	require.NoError(t, db.Create(doc).Error)
	return doc
}

func TestDocumentRepository_CRUD(t *testing.T) {
	db := setupDocumentTestDB(t)
	repo := NewDocumentRepository(db)
	userID := uuid.New()

	require.NoError(t, db.Exec(
		"INSERT INTO users (id, email, password_hash, full_name, role_id) VALUES (?, ?, 'hash', 'User', ?)",
		userID.String(),
		"creator@example.com",
		uuid.New().String(),
	).Error)

	doc := &entity.Document{
		ID:       uuid.New(),
		UserID:   userID,
		Title:    "My Doc",
		FileName: "doc.txt",
		FilePath: "/tmp/doc.txt",
		FileType: entity.FileTypeTXT,
		FileSize: 20,
	}

	require.NoError(t, repo.Create(context.Background(), doc))

	fetched, err := repo.FindByID(context.Background(), doc.ID)
	assert.NoError(t, err)
	assert.Equal(t, doc.Title, fetched.Title)

	docs, err := repo.FindByUserID(context.Background(), userID, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, docs, 1)

	doc.Title = "Updated"
	assert.NoError(t, repo.Update(context.Background(), doc))

	err = repo.Delete(context.Background(), doc.ID)
	assert.NoError(t, err)

	count, err := repo.CountByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, count)

	_, err = repo.FindByID(context.Background(), uuid.New())
	assert.Error(t, err)
}

func TestDocumentRepository_CountAndFindErrors(t *testing.T) {
	db := setupDocumentTestDB(t)
	repo := NewDocumentRepository(db)
	userID := uuid.New()
	createDocument(t, db, userID)

	count, err := repo.CountByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, count)

	require.NoError(t, db.Exec("DROP TABLE documents;").Error)
	_, err = repo.FindByUserID(context.Background(), userID, 5, 0)
	assert.Error(t, err)
}

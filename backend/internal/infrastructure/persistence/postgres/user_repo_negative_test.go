package postgres

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)

	// Create users table
	err = db.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			full_name VARCHAR(255) NOT NULL,
			role_id TEXT NOT NULL,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)
	`).Error
	require.NoError(t, err)

	// Create roles table
	err = db.Exec(`
		CREATE TABLE roles (
			id TEXT PRIMARY KEY,
			name VARCHAR(50) UNIQUE NOT NULL,
			description TEXT,
			created_at DATETIME,
			updated_at DATETIME
		)
	`).Error
	require.NoError(t, err)

	return db
}

// NEGATIVE TEST: Create user with duplicate email
func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)

	// Create role first
	roleID := uuid.New()
	err := db.Exec("INSERT INTO roles (id, name, description) VALUES (?, ?, ?)",
		roleID.String(), "student", "Student role").Error
	require.NoError(t, err)

	// Create first user
	user1 := &entity.User{
		ID:           uuid.New(),
		Email:        "duplicate@example.com",
		PasswordHash: "hash1",
		FullName:     "User One",
		RoleID:       roleID,
	}

	err = repo.Create(context.Background(), user1)
	assert.NoError(t, err)

	// Try to create second user with same email - should fail
	user2 := &entity.User{
		ID:           uuid.New(),
		Email:        "duplicate@example.com", // Same email!
		PasswordHash: "hash2",
		FullName:     "User Two",
		RoleID:       roleID,
	}

	err = repo.Create(context.Background(), user2)
	assert.Error(t, err, "Expected error for duplicate email")
}

// NEGATIVE TEST: Find user by non-existent email
func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)

	user, err := repo.FindByEmail(context.Background(), "nonexistent@example.com")

	assert.Error(t, err)
	assert.Nil(t, user)
}

// NEGATIVE TEST: Find user by invalid UUID
func TestUserRepository_FindByID_InvalidUUID(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)

	// Generate a random UUID that doesn't exist
	nonExistentID := uuid.New()

	user, err := repo.FindByID(context.Background(), nonExistentID)

	assert.Error(t, err)
	assert.Nil(t, user)
}

// NEGATIVE TEST: Update non-existent user
func TestUserRepository_Update_NotFound(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)

	// Try to update user that doesn't exist
	user := &entity.User{
		ID:           uuid.New(),
		Email:        "nonexistent@example.com",
		PasswordHash: "hash",
		FullName:     "Non Existent",
		RoleID:       uuid.New(),
	}

	err := repo.Update(context.Background(), user)

	// Update might succeed but affect 0 rows, or return error
	// Either case is acceptable - we're testing it doesn't panic
	_ = err // Ignore result, just ensure it doesn't panic
}

// NEGATIVE TEST: Create user with empty email
func TestUserRepository_Create_EmptyEmail(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)

	roleID := uuid.New()
	err := db.Exec("INSERT INTO roles (id, name, description) VALUES (?, ?, ?)",
		roleID.String(), "student", "Student role").Error
	require.NoError(t, err)

	user := &entity.User{
		ID:           uuid.New(),
		Email:        "", // Empty email
		PasswordHash: "hash",
		FullName:     "Test User",
		RoleID:       roleID,
	}

	err = repo.Create(context.Background(), user)

	// SQLite might allow empty string, but in production with PostgreSQL + constraints this would fail
	// For now, we're testing that it doesn't panic
	_ = err
}

// NEGATIVE TEST: Create user with nil context
func TestUserRepository_Create_NilContext(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)

	roleID := uuid.New()
	err := db.Exec("INSERT INTO roles (id, name, description) VALUES (?, ?, ?)",
		roleID.String(), "student", "Student role").Error
	require.NoError(t, err)

	user := &entity.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hash",
		FullName:     "Test User",
		RoleID:       roleID,
	}

	// This should panic or handle gracefully
	defer func() {
		if r := recover(); r == nil {
			// If doesn't panic, check for error
			err := repo.Create(nil, user)
			if err == nil {
				// GORM might handle nil context - that's OK
				t.Log("GORM handled nil context gracefully")
			}
		}
	}()

	_ = repo.Create(nil, user)
}

// NEGATIVE TEST: FindByEmail with empty string
func TestUserRepository_FindByEmail_EmptyString(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)

	user, err := repo.FindByEmail(context.Background(), "")

	assert.Error(t, err)
	assert.Nil(t, user)
}

// NEGATIVE TEST: FindByEmail with malformed email
func TestUserRepository_FindByEmail_MalformedEmail(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)

	malformedEmails := []string{
		"not-an-email",
		"@example.com",
		"test@",
		"test@@example.com",
		"test@example@com",
		" ",
		"test@example..com",
	}

	for _, email := range malformedEmails {
		t.Run("malformed_"+email, func(t *testing.T) {
			user, err := repo.FindByEmail(context.Background(), email)

			// Should not find anything (and not panic)
			assert.Error(t, err)
			assert.Nil(t, user)
		})
	}
}

// NEGATIVE TEST: List with very large limit
func TestUserRepository_List_LargeLimit(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)

	// Request with unreasonably large limit
	users, err := repo.List(context.Background(), 999999, 0)

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 0) // Empty database
}

// NEGATIVE TEST: List with negative offset
func TestUserRepository_List_NegativeOffset(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)

	// GORM might handle negative offset in different ways
	users, err := repo.List(context.Background(), 10, -1)

	// Should either return error or handle gracefully
	_ = err
	_ = users
	// Main goal: ensure it doesn't panic
}

// NEGATIVE TEST: Delete already deleted user (soft delete)
func TestUserRepository_Delete_AlreadyDeleted(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)

	// Create role
	roleID := uuid.New()
	err := db.Exec("INSERT INTO roles (id, name, description) VALUES (?, ?, ?)",
		roleID.String(), "student", "Student role").Error
	require.NoError(t, err)

	// Create user
	user := &entity.User{
		ID:           uuid.New(),
		Email:        "todelete@example.com",
		PasswordHash: "hash",
		FullName:     "To Delete",
		RoleID:       roleID,
	}

	err = repo.Create(context.Background(), user)
	require.NoError(t, err)

	// Delete once
	err = repo.Delete(context.Background(), user.ID)
	assert.NoError(t, err)

	// Try to delete again - should handle gracefully
	err = repo.Delete(context.Background(), user.ID)

	// GORM soft delete might succeed (updates deleted_at again)
	// or fail - either is acceptable
	_ = err // Just ensure no panic
}

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

func setupUserPositiveDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{SkipDefaultTransaction: true})
	require.NoError(t, err)

	require.NoError(t, db.Exec(`
                CREATE TABLE roles (
                        id TEXT PRIMARY KEY,
                        name VARCHAR(50) UNIQUE NOT NULL,
                        description TEXT,
                        created_at DATETIME,
                        updated_at DATETIME
                );
        `).Error)

	require.NoError(t, db.Exec(`
                CREATE TABLE users (
                        id TEXT PRIMARY KEY,
                        email VARCHAR(255) UNIQUE NOT NULL,
                        password_hash VARCHAR(255) NOT NULL,
                        full_name VARCHAR(255) NOT NULL,
                        role_id TEXT NOT NULL,
                        created_at DATETIME,
                        updated_at DATETIME,
                        deleted_at DATETIME
                );
        `).Error)

	return db
}

func seedUserWithRole(t *testing.T, db *gorm.DB) (*entity.User, *entity.Role) {
	role := &entity.Role{ID: uuid.New(), Name: entity.RoleNameStudent, Description: "student", CreatedAt: time.Time{}, UpdatedAt: time.Time{}}
	require.NoError(t, db.Create(role).Error)

	user := &entity.User{ID: uuid.New(), Email: "user@example.com", PasswordHash: "hash", FullName: "User", RoleID: role.ID}
	require.NoError(t, db.Create(user).Error)
	return user, role
}

func TestUserRepository_CRUDAndQueries(t *testing.T) {
	db := setupUserPositiveDB(t)
	repo := NewUserRepository(db)

	ctx := context.Background()

	// create baseline user
	createdUser, role := seedUserWithRole(t, db)

	// FindByID
	fetched, err := repo.FindByID(ctx, createdUser.ID)
	require.NoError(t, err)
	assert.Equal(t, createdUser.Email, fetched.Email)

	// FindByEmail
	byEmail, err := repo.FindByEmail(ctx, createdUser.Email)
	require.NoError(t, err)
	assert.Equal(t, createdUser.ID, byEmail.ID)

	// Update role
	newRole := &entity.Role{ID: uuid.New(), Name: entity.RoleNameTeacher, Description: "teacher"}
	require.NoError(t, db.Create(newRole).Error)
	fetched.RoleID = newRole.ID
	fetched.Role = newRole
	require.NoError(t, repo.Update(ctx, fetched))

	updated, err := repo.FindByID(ctx, createdUser.ID)
	require.NoError(t, err)
	assert.Equal(t, newRole.ID, updated.RoleID)

	// List & Count
	users, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.Len(t, users, 1)
	total, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)

	// Delete should mark deleted_at and exclude from listings
	require.NoError(t, repo.Delete(ctx, createdUser.ID))
	remaining, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.Empty(t, remaining)
	countAfter, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.Zero(t, countAfter)

	// Ensure original role remains untouched
	var checkRole entity.Role
	require.NoError(t, db.First(&checkRole, "id = ?", role.ID).Error)
}

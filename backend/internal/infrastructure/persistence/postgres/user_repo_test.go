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

// TestUserRepository_UpdateRoleWithPreloadedAssociation tests that role_id is properly updated
// even when the Role association is preloaded (regression test for role update bug)
func TestUserRepository_UpdateRoleWithPreloadedAssociation(t *testing.T) {
	db := setupUserPositiveDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create student role
	studentRole := &entity.Role{
		ID:          uuid.New(),
		Name:        entity.RoleNameStudent,
		Description: "Student role",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	require.NoError(t, db.Create(studentRole).Error)

	// Create teacher role
	teacherRole := &entity.Role{
		ID:          uuid.New(),
		Name:        entity.RoleNameTeacher,
		Description: "Teacher role",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	require.NoError(t, db.Create(teacherRole).Error)

	// Create user with student role
	user := &entity.User{
		ID:           uuid.New(),
		Email:        "student@test.com",
		PasswordHash: "hash",
		FullName:     "Test Student",
		RoleID:       studentRole.ID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	require.NoError(t, db.Create(user).Error)

	// Fetch user with preloaded Role association (simulates real handler behavior)
	fetchedUser, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, studentRole.ID, fetchedUser.RoleID)
	assert.NotNil(t, fetchedUser.Role)
	assert.Equal(t, entity.RoleNameStudent, fetchedUser.Role.Name)

	// Update user role to teacher (with preloaded association present)
	fetchedUser.RoleID = teacherRole.ID
	fetchedUser.UpdatedAt = time.Now()
	err = repo.Update(ctx, fetchedUser)
	require.NoError(t, err)

	// Verify role was actually updated in database
	verifyUser, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, teacherRole.ID, verifyUser.RoleID, "role_id should be updated to teacher role")
	assert.NotNil(t, verifyUser.Role)
	assert.Equal(t, entity.RoleNameTeacher, verifyUser.Role.Name, "preloaded Role should reflect teacher role")

	// Verify in raw database
	var dbUser entity.User
	require.NoError(t, db.First(&dbUser, "id = ?", user.ID).Error)
	assert.Equal(t, teacherRole.ID, dbUser.RoleID, "role_id should be persisted in database")
}

// TestUserRepository_UpdateMultipleFields tests that Update properly updates all specified fields
func TestUserRepository_UpdateMultipleFields(t *testing.T) {
	db := setupUserPositiveDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create role
	role := &entity.Role{
		ID:          uuid.New(),
		Name:        entity.RoleNameStudent,
		Description: "Student",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	require.NoError(t, db.Create(role).Error)

	// Create user
	user := &entity.User{
		ID:           uuid.New(),
		Email:        "old@test.com",
		PasswordHash: "oldhash",
		FullName:     "Old Name",
		RoleID:       role.ID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	require.NoError(t, db.Create(user).Error)

	// Fetch and update multiple fields
	fetchedUser, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)

	fetchedUser.Email = "new@test.com"
	fetchedUser.FullName = "New Name"
	fetchedUser.PasswordHash = "newhash"
	fetchedUser.UpdatedAt = time.Now()

	err = repo.Update(ctx, fetchedUser)
	require.NoError(t, err)

	// Verify all fields updated
	updated, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "new@test.com", updated.Email)
	assert.Equal(t, "New Name", updated.FullName)
	assert.Equal(t, "newhash", updated.PasswordHash)
	assert.True(t, updated.UpdatedAt.After(user.UpdatedAt))
}

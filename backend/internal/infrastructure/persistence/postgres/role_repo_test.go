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

func setupRoleTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		// Disable default transaction for better performance in tests
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)

	// Manually create table without PostgreSQL-specific defaults
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

func seedRoles(t *testing.T, db *gorm.DB) {
	now := time.Time{}
	roles := []*entity.Role{
		{
			ID:          uuid.New(),
			Name:        entity.RoleNameAdmin,
			Description: "Administrator with full access",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          uuid.New(),
			Name:        entity.RoleNameTeacher,
			Description: "Teacher who creates tests",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          uuid.New(),
			Name:        entity.RoleNameStudent,
			Description: "Student (default role)",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	for _, role := range roles {
		err := db.Create(role).Error
		require.NoError(t, err)
	}
}

func TestRoleRepository_FindByID(t *testing.T) {
	db := setupRoleTestDB(t)
	seedRoles(t, db)
	repo := NewRoleRepository(db)

	tests := []struct {
		name    string
		setup   func() uuid.UUID
		wantErr bool
	}{
		{
			name: "should find existing role by ID",
			setup: func() uuid.UUID {
				var role entity.Role
				err := db.Where("name = ?", entity.RoleNameAdmin).First(&role).Error
				require.NoError(t, err)
				return role.ID
			},
			wantErr: false,
		},
		{
			name: "should return error for non-existent role ID",
			setup: func() uuid.UUID {
				return uuid.New()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.setup()
			role, err := repo.FindByID(context.Background(), id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, role)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, role)
				assert.Equal(t, id, role.ID)
			}
		})
	}
}

func TestRoleRepository_FindByName(t *testing.T) {
	db := setupRoleTestDB(t)
	seedRoles(t, db)
	repo := NewRoleRepository(db)

	tests := []struct {
		name     string
		roleName entity.RoleName
		wantErr  bool
	}{
		{
			name:     "should find admin role by name",
			roleName: entity.RoleNameAdmin,
			wantErr:  false,
		},
		{
			name:     "should find teacher role by name",
			roleName: entity.RoleNameTeacher,
			wantErr:  false,
		},
		{
			name:     "should find student role by name",
			roleName: entity.RoleNameStudent,
			wantErr:  false,
		},
		{
			name:     "should return error for non-existent role name",
			roleName: "invalid_role",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role, err := repo.FindByName(context.Background(), tt.roleName)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, role)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, role)
				assert.Equal(t, tt.roleName, role.Name)
			}
		})
	}
}

func TestRoleRepository_List(t *testing.T) {
	db := setupRoleTestDB(t)
	seedRoles(t, db)
	repo := NewRoleRepository(db)

	roles, err := repo.List(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, roles)
	assert.Len(t, roles, 3)

	// Verify all three roles are present
	roleNames := make(map[entity.RoleName]bool)
	for _, role := range roles {
		roleNames[role.Name] = true
	}

	assert.True(t, roleNames[entity.RoleNameAdmin])
	assert.True(t, roleNames[entity.RoleNameTeacher])
	assert.True(t, roleNames[entity.RoleNameStudent])
}

func TestRoleRepository_List_EmptyDatabase(t *testing.T) {
	db := setupRoleTestDB(t)
	repo := NewRoleRepository(db)

	roles, err := repo.List(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, roles)
	assert.Len(t, roles, 0)
}

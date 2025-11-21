package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUser_SetPassword(t *testing.T) {
	user := &User{}
	password := "testpassword123"

	err := user.SetPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.PasswordHash)
	assert.NotEqual(t, password, user.PasswordHash)

	// Verify password can be checked
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	assert.NoError(t, err)
}

func TestUser_CheckPassword(t *testing.T) {
	user := &User{}
	password := "testpassword123"

	err := user.SetPassword(password)
	assert.NoError(t, err)

	// Correct password
	assert.True(t, user.CheckPassword(password))

	// Wrong password
	assert.False(t, user.CheckPassword("wrongpassword"))
}

func TestUser_IsTeacher(t *testing.T) {
	teacherRole := &Role{
		ID:   uuid.New(),
		Name: RoleNameTeacher,
	}
	user := &User{
		RoleID: teacherRole.ID,
		Role:   teacherRole,
	}
	assert.True(t, user.IsTeacher())
	assert.False(t, user.IsStudent())
	assert.False(t, user.IsAdmin())
}

func TestUser_IsStudent(t *testing.T) {
	studentRole := &Role{
		ID:   uuid.New(),
		Name: RoleNameStudent,
	}
	user := &User{
		RoleID: studentRole.ID,
		Role:   studentRole,
	}
	assert.False(t, user.IsTeacher())
	assert.True(t, user.IsStudent())
	assert.False(t, user.IsAdmin())
}

func TestUser_IsAdmin(t *testing.T) {
	adminRole := &Role{
		ID:   uuid.New(),
		Name: RoleNameAdmin,
	}
	user := &User{
		RoleID: adminRole.ID,
		Role:   adminRole,
	}
	assert.False(t, user.IsTeacher())
	assert.False(t, user.IsStudent())
	assert.True(t, user.IsAdmin())
}

func TestUser_TableName(t *testing.T) {
	user := User{}
	assert.Equal(t, "users", user.TableName())
}

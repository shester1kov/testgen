package entity

import (
	"testing"

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
	user := &User{Role: RoleTeacher}
	assert.True(t, user.IsTeacher())
	assert.False(t, user.IsStudent())
	assert.False(t, user.IsAdmin())
}

func TestUser_IsStudent(t *testing.T) {
	user := &User{Role: RoleStudent}
	assert.False(t, user.IsTeacher())
	assert.True(t, user.IsStudent())
	assert.False(t, user.IsAdmin())
}

func TestUser_IsAdmin(t *testing.T) {
	user := &User{Role: RoleAdmin}
	assert.False(t, user.IsTeacher())
	assert.False(t, user.IsStudent())
	assert.True(t, user.IsAdmin())
}

func TestUser_TableName(t *testing.T) {
	user := User{}
	assert.Equal(t, "users", user.TableName())
}

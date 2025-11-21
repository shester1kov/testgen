package entity

import (
	"time"

	"github.com/google/uuid"
)

// RoleName represents the name of a role
type RoleName string

const (
	RoleNameAdmin   RoleName = "admin"
	RoleNameTeacher RoleName = "teacher"
	RoleNameStudent RoleName = "student"
)

// Role represents a user role in the system
type Role struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        RoleName   `json:"name" gorm:"type:varchar(50);uniqueIndex;not null"`
	Description string     `json:"description" gorm:"type:text"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (Role) TableName() string {
	return "roles"
}

// IsAdmin checks if the role is admin
func (r *Role) IsAdmin() bool {
	return r.Name == RoleNameAdmin
}

// IsTeacher checks if the role is teacher
func (r *Role) IsTeacher() bool {
	return r.Name == RoleNameTeacher
}

// IsStudent checks if the role is student
func (r *Role) IsStudent() bool {
	return r.Name == RoleNameStudent
}

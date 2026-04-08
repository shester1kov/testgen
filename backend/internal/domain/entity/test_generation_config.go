package entity

import (
	"time"

	"github.com/google/uuid"
)

// TestGenerationConfig is an immutable record of the LLM parameters used
// when a test was AI-generated. Absent for manually created tests.
type TestGenerationConfig struct {
	TestID             uuid.UUID `gorm:"type:uuid;primary_key"`
	AcademicProfileID  uuid.UUID `gorm:"type:uuid;not null"`
	Difficulty         string    `gorm:"type:varchar(50);not null"`
	NumQuestions       int       `gorm:"not null"`
	LLMProvider        string    `gorm:"type:varchar(50);not null"`
	Language           string    `gorm:"type:varchar(10);not null;default:'ru'"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`

	AcademicProfile AcademicProfile `gorm:"foreignKey:AcademicProfileID"`
}

// TableName specifies the table name for GORM.
func (TestGenerationConfig) TableName() string {
	return "test_generation_configs"
}

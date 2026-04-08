package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ExampleAnswer is a single answer option inside a question example stored in the DB.
type ExampleAnswer struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

// ExampleAnswers is a slice of ExampleAnswer that serialises to/from PostgreSQL JSONB.
type ExampleAnswers []ExampleAnswer

// Value implements driver.Valuer for GORM/database/sql.
func (a ExampleAnswers) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan implements sql.Scanner for GORM/database/sql.
func (a *ExampleAnswers) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("ExampleAnswers.Scan: expected []byte")
	}
	return json.Unmarshal(b, a)
}

// FormulationExample is a ПЛОХО/ХОРОШО phrasing pair for a given academic profile.
type FormulationExample struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	AcademicProfileID  uuid.UUID `gorm:"type:uuid;not null"`
	BadExample         string    `gorm:"type:text;not null"`
	GoodExample        string    `gorm:"type:text;not null"`
	OrderNum           int       `gorm:"not null;default:0"`
}

// TableName specifies the table name for GORM.
func (FormulationExample) TableName() string {
	return "academic_profile_formulation_examples"
}

// QuestionExample is a domain-appropriate example question stored per profile and type.
// The difficulty field is NOT stored here; it is injected at prompt-build time.
type QuestionExample struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	AcademicProfileID  uuid.UUID      `gorm:"type:uuid;not null"`
	QuestionType       string         `gorm:"type:varchar(50);not null"`
	QuestionText       string         `gorm:"type:text;not null"`
	Answers            ExampleAnswers `gorm:"type:jsonb;not null"`
	Explanation        string         `gorm:"type:text;not null"`
	OrderNum           int            `gorm:"not null;default:0"`
}

// TableName specifies the table name for GORM.
func (QuestionExample) TableName() string {
	return "academic_profile_question_examples"
}

// AcademicProfile is the DB entity for an educational discipline profile.
type AcademicProfile struct {
	ID                 uuid.UUID            `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Code               string               `gorm:"type:varchar(50);uniqueIndex;not null"`
	DisplayName        string               `gorm:"type:varchar(100);not null"`
	Description        string               `gorm:"type:text"`
	ProfileInstruction string               `gorm:"type:text"`
	Temperature        float64              `gorm:"type:numeric(3,2);not null;default:0.50"`
	SortOrder          int                  `gorm:"not null;default:0"`
	CreatedAt          time.Time            `gorm:"autoCreateTime"`
	Formulations       []FormulationExample `gorm:"foreignKey:AcademicProfileID"`
	QuestionExamples   []QuestionExample    `gorm:"foreignKey:AcademicProfileID"`
}

// TableName specifies the table name for GORM.
func (AcademicProfile) TableName() string {
	return "academic_profiles"
}

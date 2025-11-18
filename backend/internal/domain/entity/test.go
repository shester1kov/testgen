package entity

import (
	"time"

	"github.com/google/uuid"
)

type TestStatus string

const (
	TestStatusDraft     TestStatus = "draft"
	TestStatusPublished TestStatus = "published"
	TestStatusArchived  TestStatus = "archived"
)

type Test struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID         uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	DocumentID     *uuid.UUID `json:"document_id,omitempty" gorm:"type:uuid;index"`
	Title          string     `json:"title" gorm:"type:varchar(500);not null"`
	Description    string     `json:"description,omitempty" gorm:"type:text"`
	TotalQuestions int        `json:"total_questions" gorm:"default:0"`
	Status         TestStatus `json:"status" gorm:"type:varchar(50);default:'draft';index"`
	MoodleSynced   bool       `json:"moodle_synced" gorm:"default:false"`
	MoodleTestID   string     `json:"moodle_test_id,omitempty" gorm:"type:varchar(255)"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relations
	User      User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Document  *Document  `json:"document,omitempty" gorm:"foreignKey:DocumentID"`
	Questions []Question `json:"questions,omitempty" gorm:"foreignKey:TestID"`
}

// TableName specifies the table name for GORM
func (Test) TableName() string {
	return "tests"
}

// IsPublished checks if test is published
func (t *Test) IsPublished() bool {
	return t.Status == TestStatusPublished
}

// Publish marks test as published
func (t *Test) Publish() {
	t.Status = TestStatusPublished
}

// Archive marks test as archived
func (t *Test) Archive() {
	t.Status = TestStatusArchived
}

// UpdateQuestionsCount updates total questions count
func (t *Test) UpdateQuestionsCount(count int) {
	t.TotalQuestions = count
}

// MarkMoodleSynced marks test as synced with Moodle
func (t *Test) MarkMoodleSynced(moodleID string) {
	t.MoodleSynced = true
	t.MoodleTestID = moodleID
}

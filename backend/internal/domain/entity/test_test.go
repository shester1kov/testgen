package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTest_TableName(t *testing.T) {
	test := Test{}
	assert.Equal(t, "tests", test.TableName())
}

func TestTest_IsPublished(t *testing.T) {
	test := &Test{Status: TestStatusPublished}
	assert.True(t, test.IsPublished())

	test.Status = TestStatusDraft
	assert.False(t, test.IsPublished())
}

func TestTest_Publish(t *testing.T) {
	test := &Test{Status: TestStatusDraft}
	test.Publish()
	assert.Equal(t, TestStatusPublished, test.Status)
}

func TestTest_Archive(t *testing.T) {
	test := &Test{Status: TestStatusPublished}
	test.Archive()
	assert.Equal(t, TestStatusArchived, test.Status)
}

func TestTest_UpdateQuestionsCount(t *testing.T) {
	test := &Test{TotalQuestions: 0}
	test.UpdateQuestionsCount(10)
	assert.Equal(t, 10, test.TotalQuestions)
}

func TestTest_MarkMoodleSynced(t *testing.T) {
	test := &Test{MoodleSynced: false}
	moodleTestID := "moodle-123"
	test.MarkMoodleSynced(moodleTestID)

	assert.True(t, test.MoodleSynced)
	assert.Equal(t, moodleTestID, test.MoodleTestID)
}

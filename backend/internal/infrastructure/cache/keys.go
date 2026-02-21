package cache

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Cache key prefixes and TTL constants
const (
	// Key prefixes
	UserByIDPrefix        = "user:id:"
	UserByEmailPrefix     = "user:email:"
	DocumentByIDPrefix    = "document:id:"
	DocumentsByUserPrefix = "documents:user:"
	TestByIDPrefix        = "test:id:"
	TestsByUserPrefix     = "tests:user:"
	QuestionsByTestPrefix = "questions:test:"

	// TTL values
	UserTTL      = 30 * time.Minute
	DocumentTTL  = 1 * time.Hour
	TestTTL      = 30 * time.Minute
	QuestionsTTL = 1 * time.Hour
)

// Key builders for type safety
func UserByIDKey(userID uuid.UUID) string {
	return fmt.Sprintf("%s%s", UserByIDPrefix, userID.String())
}

func UserByEmailKey(email string) string {
	return fmt.Sprintf("%s%s", UserByEmailPrefix, email)
}

func DocumentByIDKey(docID uuid.UUID) string {
	return fmt.Sprintf("%s%s", DocumentByIDPrefix, docID.String())
}

func DocumentsByUserKey(userID uuid.UUID) string {
	return fmt.Sprintf("%s%s", DocumentsByUserPrefix, userID.String())
}

func TestByIDKey(testID uuid.UUID) string {
	return fmt.Sprintf("%s%s", TestByIDPrefix, testID.String())
}

func TestsByUserKey(userID uuid.UUID) string {
	return fmt.Sprintf("%s%s", TestsByUserPrefix, userID.String())
}

func QuestionsByTestKey(testID uuid.UUID) string {
	return fmt.Sprintf("%s%s", QuestionsByTestPrefix, testID.String())
}

// Pattern builders for batch operations
func UserPattern() string {
	return "user:*"
}

func DocumentPattern(userID uuid.UUID) string {
	return fmt.Sprintf("document*%s*", userID.String())
}

func TestPattern(userID uuid.UUID) string {
	return fmt.Sprintf("test*%s*", userID.String())
}

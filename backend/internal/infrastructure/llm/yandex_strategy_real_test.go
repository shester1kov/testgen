package llm

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestYandexGPTStrategy_RealAPI tests with real YandexGPT API
// This test is skipped unless YANDEX_GPT_API_KEY and YANDEX_GPT_FOLDER_ID are set
// To run: YANDEX_GPT_API_KEY=your-key YANDEX_GPT_FOLDER_ID=your-folder go test -v -run TestYandexGPTStrategy_RealAPI
func TestYandexGPTStrategy_RealAPI(t *testing.T) {
	apiKey := os.Getenv("YANDEX_GPT_API_KEY")
	folderID := os.Getenv("YANDEX_GPT_FOLDER_ID")
	model := os.Getenv("YANDEX_GPT_MODEL")

	if apiKey == "" || folderID == "" {
		t.Skip("Skipping real API test: YANDEX_GPT_API_KEY or YANDEX_GPT_FOLDER_ID not set")
	}

	if model == "" {
		model = "yandexgpt-lite"
	}

	t.Logf("Testing with real YandexGPT API")
	t.Logf("Folder ID: %s", folderID)
	t.Logf("Model: %s", model)

	strategy := NewYandexGPTStrategy(apiKey, folderID, model)

	t.Run("generates questions from simple text", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		params := GenerationParams{
			Text: `Go — это статически типизированный компилируемый язык программирования,
разработанный внутри компании Google. Синтаксис Go похож на C, но с дополнительными
возможностями, такими как сборка мусора, типобезопасность, некоторые возможности
динамической типизации, дополнительные встроенные типы.`,
			NumQuestions:  3,
			QuestionTypes: []QuestionType{SingleChoice},
			Difficulty:    "easy",
			Language:      "ru",
		}

		t.Logf("Sending request to YandexGPT...")
		questions, err := strategy.GenerateQuestions(ctx, params)

		require.NoError(t, err, "Should successfully generate questions")
		require.NotEmpty(t, questions, "Should return at least one question")
		require.LessOrEqual(t, len(questions), 3, "Should not return more than requested questions")

		// Validate each question
		for i, q := range questions {
			t.Logf("Question %d:", i+1)
			t.Logf("  Text: %s", q.QuestionText)
			t.Logf("  Type: %s", q.QuestionType)
			t.Logf("  Difficulty: %s", q.Difficulty)
			t.Logf("  Answers: %d", len(q.Answers))

			require.NotEmpty(t, q.QuestionText, "Question text should not be empty")
			require.NotEmpty(t, q.Answers, "Should have at least one answer")

			// Check that there's at least one correct answer
			hasCorrect := false
			for _, a := range q.Answers {
				if a.IsCorrect {
					hasCorrect = true
					break
				}
			}
			require.True(t, hasCorrect, "Question should have at least one correct answer")
		}
	})

	t.Run("handles different difficulty levels", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		params := GenerationParams{
			Text:          "Python - это высокоуровневый язык программирования.",
			NumQuestions:  1,
			QuestionTypes: []QuestionType{SingleChoice},
			Difficulty:    "hard",
			Language:      "ru",
		}

		questions, err := strategy.GenerateQuestions(ctx, params)

		require.NoError(t, err)
		require.NotEmpty(t, questions)
		t.Logf("Generated hard question: %s", questions[0].QuestionText)
	})

	t.Run("handles English text", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		params := GenerationParams{
			Text:          "JavaScript is a high-level, interpreted programming language that conforms to the ECMAScript specification.",
			NumQuestions:  2,
			QuestionTypes: []QuestionType{SingleChoice},
			Difficulty:    "medium",
			Language:      "en",
		}

		questions, err := strategy.GenerateQuestions(ctx, params)

		require.NoError(t, err)
		require.NotEmpty(t, questions)
		t.Logf("Generated %d questions from English text", len(questions))
	})
}

// TestYandexGPTStrategy_RealAPI_ErrorCases tests error handling with real API
func TestYandexGPTStrategy_RealAPI_ErrorCases(t *testing.T) {
	apiKey := os.Getenv("YANDEX_GPT_API_KEY")
	folderID := os.Getenv("YANDEX_GPT_FOLDER_ID")

	if apiKey == "" || folderID == "" {
		t.Skip("Skipping real API error test: YANDEX_GPT_API_KEY or YANDEX_GPT_FOLDER_ID not set")
	}

	t.Run("invalid API key returns error", func(t *testing.T) {
		strategy := NewYandexGPTStrategy("invalid-key", folderID, "yandexgpt-lite")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		params := GenerationParams{
			Text:         "Test text",
			NumQuestions: 1,
		}

		_, err := strategy.GenerateQuestions(ctx, params)

		require.Error(t, err, "Should return error for invalid API key")
		t.Logf("Expected error received: %v", err)
	})

	t.Run("invalid folder ID returns error", func(t *testing.T) {
		strategy := NewYandexGPTStrategy(apiKey, "invalid-folder", "yandexgpt-lite")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		params := GenerationParams{
			Text:         "Test text",
			NumQuestions: 1,
		}

		_, err := strategy.GenerateQuestions(ctx, params)

		require.Error(t, err, "Should return error for invalid folder ID")
		t.Logf("Expected error received: %v", err)
	})
}

// TestYandexGPTStrategy_RealAPI_Performance tests response time
func TestYandexGPTStrategy_RealAPI_Performance(t *testing.T) {
	apiKey := os.Getenv("YANDEX_GPT_API_KEY")
	folderID := os.Getenv("YANDEX_GPT_FOLDER_ID")
	model := os.Getenv("YANDEX_GPT_MODEL")

	if apiKey == "" || folderID == "" {
		t.Skip("Skipping real API performance test: environment variables not set")
	}

	if model == "" {
		model = "yandexgpt-lite"
	}

	strategy := NewYandexGPTStrategy(apiKey, folderID, model)

	t.Run("measures response time for small request", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		params := GenerationParams{
			Text:         "Короткий текст для быстрого теста.",
			NumQuestions: 1,
			Difficulty:   "easy",
		}

		start := time.Now()
		questions, err := strategy.GenerateQuestions(ctx, params)
		duration := time.Since(start)

		require.NoError(t, err)
		require.NotEmpty(t, questions)

		t.Logf("Response time: %v", duration)
		t.Logf("Questions generated: %d", len(questions))

		// YandexGPT lite should respond within reasonable time
		if duration > 30*time.Second {
			t.Logf("WARNING: Response time is high (%v). Consider checking network or API load.", duration)
		}
	})
}

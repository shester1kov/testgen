package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestYandexGPTStrategy_BuildPrompt(t *testing.T) {
	strategy := NewYandexGPTStrategy("test-key", "test-folder", "yandexgpt-lite")

	t.Run("builds prompt with all parameters", func(t *testing.T) {
		params := GenerationParams{
			Text:          "Test document text",
			NumQuestions:  5,
			QuestionTypes: []QuestionType{SingleChoice, MultipleChoice},
			Difficulty:    "hard",
			Language:      "ru",
		}

		prompt := strategy.buildPrompt(params)

		require.Contains(t, prompt, "5 тестовых вопросов")
		require.Contains(t, prompt, "Test document text")
		require.Contains(t, prompt, "single_choice, multiple_choice")
		require.Contains(t, prompt, "hard")
		require.Contains(t, prompt, "ru")
	})

	t.Run("uses default values when not provided", func(t *testing.T) {
		params := GenerationParams{
			Text:         "Test text",
			NumQuestions: 3,
		}

		prompt := strategy.buildPrompt(params)

		require.Contains(t, prompt, "single_choice")
		require.Contains(t, prompt, "medium")
		require.Contains(t, prompt, "ru")
	})

	t.Run("handles multiple question types", func(t *testing.T) {
		params := GenerationParams{
			Text:          "Test text",
			NumQuestions:  2,
			QuestionTypes: []QuestionType{TrueFalse, ShortAnswer},
		}

		prompt := strategy.buildPrompt(params)

		require.Contains(t, prompt, "true_false, short_answer")
	})
}

func TestYandexGPTStrategy_ParseQuestions(t *testing.T) {
	strategy := NewYandexGPTStrategy("test-key", "test-folder", "yandexgpt-lite")

	t.Run("parses valid JSON response", func(t *testing.T) {
		jsonResponse := `{
			"questions": [
				{
					"question": "What is Go?",
					"type": "single_choice",
					"difficulty": "easy",
					"answers": [
						{"text": "A programming language", "is_correct": true},
						{"text": "A database", "is_correct": false},
						{"text": "A framework", "is_correct": false},
						{"text": "An IDE", "is_correct": false}
					],
					"explanation": "Go is a programming language developed by Google"
				},
				{
					"question": "Is Go statically typed?",
					"type": "true_false",
					"difficulty": "medium",
					"answers": [
						{"text": "True", "is_correct": true},
						{"text": "False", "is_correct": false}
					]
				}
			]
		}`

		params := GenerationParams{NumQuestions: 2}
		questions, err := strategy.parseQuestions(jsonResponse, params)

		require.NoError(t, err)
		require.Len(t, questions, 2)

		// Check first question
		require.Equal(t, "What is Go?", questions[0].QuestionText)
		require.Equal(t, SingleChoice, questions[0].QuestionType)
		require.Equal(t, "easy", questions[0].Difficulty)
		require.Len(t, questions[0].Answers, 4)
		require.True(t, questions[0].Answers[0].IsCorrect)
		require.Equal(t, "Go is a programming language developed by Google", questions[0].Explanation)

		// Check second question
		require.Equal(t, "Is Go statically typed?", questions[1].QuestionText)
		require.Equal(t, TrueFalse, questions[1].QuestionType)
		require.Len(t, questions[1].Answers, 2)
	})

	t.Run("parses JSON wrapped in markdown code blocks", func(t *testing.T) {
		jsonResponse := "```json\n" + `{
			"questions": [
				{
					"question": "Test question",
					"type": "single_choice",
					"difficulty": "easy",
					"answers": [
						{"text": "Answer 1", "is_correct": true}
					]
				}
			]
		}` + "\n```"

		params := GenerationParams{NumQuestions: 1}
		questions, err := strategy.parseQuestions(jsonResponse, params)

		require.NoError(t, err)
		require.Len(t, questions, 1)
		require.Equal(t, "Test question", questions[0].QuestionText)
	})

	t.Run("parses JSON with just backticks", func(t *testing.T) {
		jsonResponse := "```\n" + `{
			"questions": [
				{
					"question": "Another test",
					"type": "multiple_choice",
					"difficulty": "hard",
					"answers": [
						{"text": "Option A", "is_correct": true},
						{"text": "Option B", "is_correct": true},
						{"text": "Option C", "is_correct": false}
					]
				}
			]
		}` + "\n```"

		params := GenerationParams{NumQuestions: 1}
		questions, err := strategy.parseQuestions(jsonResponse, params)

		require.NoError(t, err)
		require.Len(t, questions, 1)
		require.Equal(t, MultipleChoice, questions[0].QuestionType)
	})

	t.Run("returns error on invalid JSON", func(t *testing.T) {
		invalidJSON := "This is not JSON"

		params := GenerationParams{NumQuestions: 1}
		_, err := strategy.parseQuestions(invalidJSON, params)

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to parse JSON")
	})

	t.Run("returns error when no questions in response", func(t *testing.T) {
		emptyResponse := `{"questions": []}`

		params := GenerationParams{NumQuestions: 1}
		_, err := strategy.parseQuestions(emptyResponse, params)

		require.Error(t, err)
		require.Contains(t, err.Error(), "no questions generated")
	})

	t.Run("handles questions without explanation", func(t *testing.T) {
		jsonResponse := `{
			"questions": [
				{
					"question": "Simple question",
					"type": "single_choice",
					"difficulty": "easy",
					"answers": [
						{"text": "Yes", "is_correct": true},
						{"text": "No", "is_correct": false}
					]
				}
			]
		}`

		params := GenerationParams{NumQuestions: 1}
		questions, err := strategy.parseQuestions(jsonResponse, params)

		require.NoError(t, err)
		require.Len(t, questions, 1)
		require.Empty(t, questions[0].Explanation)
	})
}

func TestNewYandexGPTStrategy_DefaultModel(t *testing.T) {
	t.Run("uses provided model", func(t *testing.T) {
		strategy := NewYandexGPTStrategy("key", "folder", "yandexgpt")
		require.Equal(t, "yandexgpt", strategy.model)
	})

	t.Run("uses default model when empty", func(t *testing.T) {
		strategy := NewYandexGPTStrategy("key", "folder", "")
		require.Equal(t, "yandexgpt-lite", strategy.model)
	})

	t.Run("sets all fields correctly", func(t *testing.T) {
		strategy := NewYandexGPTStrategy("my-key", "my-folder", "my-model")
		require.Equal(t, "my-key", strategy.apiKey)
		require.Equal(t, "my-folder", strategy.folderID)
		require.Equal(t, "my-model", strategy.model)
		require.NotNil(t, strategy.client)
	})
}

func TestYandexGPTStrategy_GenerateQuestions_Integration(t *testing.T) {
	t.Run("successfully generates questions with valid response", func(t *testing.T) {
		// Create mock server that returns valid YandexGPT response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request headers
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
			require.Equal(t, "Api-Key test-api-key", r.Header.Get("Authorization"))

			// Return valid YandexGPT response
			response := YandexGPTResponse{
				Result: YandexResult{
					Alternatives: []YandexAlternative{
						{
							Message: YandexMessage{
								Role: "assistant",
								Text: `{
									"questions": [
										{
											"question": "Что такое Go?",
											"type": "single_choice",
											"difficulty": "easy",
											"answers": [
												{"text": "Язык программирования", "is_correct": true},
												{"text": "База данных", "is_correct": false},
												{"text": "Фреймворк", "is_correct": false},
												{"text": "IDE", "is_correct": false}
											],
											"explanation": "Go - это язык программирования"
										}
									]
								}`,
							},
							Status: "ALTERNATIVE_STATUS_FINAL",
						},
					},
					Usage: YandexUsage{
						InputTextTokens:  "100",
						CompletionTokens: "200",
						TotalTokens:      "300",
					},
					ModelVersion: "yandexgpt-lite/latest",
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		// Create strategy and override baseURL to use mock server
		strategy := NewYandexGPTStrategy("test-api-key", "test-folder", "yandexgpt-lite")
		strategy.baseURL = server.URL
		strategy.client = &http.Client{
			Timeout: 5 * time.Second,
		}

		// Test generation
		ctx := context.Background()
		params := GenerationParams{
			Text:          "Тестовый текст о языке программирования Go",
			NumQuestions:  1,
			QuestionTypes: []QuestionType{SingleChoice},
			Difficulty:    "easy",
			Language:      "ru",
		}

		questions, err := strategy.GenerateQuestions(ctx, params)

		require.NoError(t, err)
		require.Len(t, questions, 1)
		require.Equal(t, "Что такое Go?", questions[0].QuestionText)
		require.Equal(t, SingleChoice, questions[0].QuestionType)
		require.Equal(t, "easy", questions[0].Difficulty)
		require.Len(t, questions[0].Answers, 4)
		require.True(t, questions[0].Answers[0].IsCorrect)
	})

	t.Run("handles YandexGPT API error", func(t *testing.T) {
		// Create mock server that returns error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": {"code": "UNAUTHENTICATED", "message": "Invalid API key"}}`))
		}))
		defer server.Close()

		strategy := NewYandexGPTStrategy("invalid-key", "test-folder", "yandexgpt-lite")
		strategy.baseURL = server.URL
		strategy.client = &http.Client{
			Timeout: 5 * time.Second,
		}

		ctx := context.Background()
		params := GenerationParams{
			Text:         "Test text",
			NumQuestions: 1,
		}

		_, err := strategy.GenerateQuestions(ctx, params)

		require.Error(t, err)
		require.Contains(t, err.Error(), "yandexgpt API error")
		require.Contains(t, err.Error(), "401")
	})

	t.Run("handles malformed JSON in generated text", func(t *testing.T) {
		// Create mock server that returns valid YandexGPT response but with invalid question JSON
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := YandexGPTResponse{
				Result: YandexResult{
					Alternatives: []YandexAlternative{
						{
							Message: YandexMessage{
								Role: "assistant",
								Text: "Это не валидный JSON, а просто текст",
							},
							Status: "ALTERNATIVE_STATUS_FINAL",
						},
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		strategy := NewYandexGPTStrategy("test-key", "test-folder", "yandexgpt-lite")
		strategy.baseURL = server.URL
		strategy.client = &http.Client{
			Timeout: 5 * time.Second,
		}

		ctx := context.Background()
		params := GenerationParams{
			Text:         "Test text",
			NumQuestions: 1,
		}

		_, err := strategy.GenerateQuestions(ctx, params)

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to parse generated questions")
	})

	t.Run("handles empty alternatives in response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := YandexGPTResponse{
				Result: YandexResult{
					Alternatives: []YandexAlternative{}, // Empty alternatives
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		strategy := NewYandexGPTStrategy("test-key", "test-folder", "yandexgpt-lite")
		strategy.baseURL = server.URL
		strategy.client = &http.Client{
			Timeout: 5 * time.Second,
		}

		ctx := context.Background()
		params := GenerationParams{
			Text:         "Test text",
			NumQuestions: 1,
		}

		_, err := strategy.GenerateQuestions(ctx, params)

		require.Error(t, err)
		require.Contains(t, err.Error(), "no alternatives in response")
	})

	t.Run("handles context timeout", func(t *testing.T) {
		// Create mock server that delays response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second) // Delay longer than context timeout
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		strategy := NewYandexGPTStrategy("test-key", "test-folder", "yandexgpt-lite")
		strategy.baseURL = server.URL
		strategy.client = &http.Client{
			Timeout: 5 * time.Second,
		}

		// Create context with short timeout
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		params := GenerationParams{
			Text:         "Test text",
			NumQuestions: 1,
		}

		_, err := strategy.GenerateQuestions(ctx, params)

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to send request")
	})

	t.Run("validates required API key", func(t *testing.T) {
		strategy := NewYandexGPTStrategy("", "test-folder", "yandexgpt-lite")

		ctx := context.Background()
		params := GenerationParams{
			Text:         "Test text",
			NumQuestions: 1,
		}

		_, err := strategy.GenerateQuestions(ctx, params)

		require.Error(t, err)
		require.Contains(t, err.Error(), "API key not configured")
	})

	t.Run("validates required folder ID", func(t *testing.T) {
		strategy := NewYandexGPTStrategy("test-key", "", "yandexgpt-lite")

		ctx := context.Background()
		params := GenerationParams{
			Text:         "Test text",
			NumQuestions: 1,
		}

		_, err := strategy.GenerateQuestions(ctx, params)

		require.Error(t, err)
		require.Contains(t, err.Error(), "folder ID not configured")
	})
}

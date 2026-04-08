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

func TestBuildSystemPrompt_Difficulty(t *testing.T) {
	t.Run("easy: all single_choice, no multiple_choice mention", func(t *testing.T) {
		params := GenerationParams{
			NumQuestions:        5,
			Difficulty:          "easy",
			MultipleChoiceCount: 0,
		}

		prompt := BuildSystemPrompt(params)

		require.Contains(t, prompt, "5 тестовых вопросов типа single_choice")
		require.Contains(t, prompt, "easy")
		require.NotContains(t, prompt, "multiple_choice")
	})

	t.Run("medium: split between single_choice and multiple_choice", func(t *testing.T) {
		params := GenerationParams{
			NumQuestions:        10,
			Difficulty:          "medium",
			MultipleChoiceCount: 3, // ceil(10 * 0.3)
		}

		prompt := BuildSystemPrompt(params)

		require.Contains(t, prompt, "10 тестовых вопросов")
		require.Contains(t, prompt, "7 вопросов типа single_choice")
		require.Contains(t, prompt, "3 вопросов типа multiple_choice")
		require.Contains(t, prompt, "medium")
	})

	t.Run("hard: majority multiple_choice", func(t *testing.T) {
		params := GenerationParams{
			NumQuestions:        10,
			Difficulty:          "hard",
			MultipleChoiceCount: 6, // ceil(10 * 0.6)
		}

		prompt := BuildSystemPrompt(params)

		require.Contains(t, prompt, "10 тестовых вопросов")
		require.Contains(t, prompt, "4 вопросов типа single_choice")
		require.Contains(t, prompt, "6 вопросов типа multiple_choice")
		require.Contains(t, prompt, "hard")
	})

	t.Run("defaults difficulty to medium when empty", func(t *testing.T) {
		params := GenerationParams{
			NumQuestions:        3,
			MultipleChoiceCount: 0,
		}

		prompt := BuildSystemPrompt(params)

		require.Contains(t, prompt, "medium")
	})

	t.Run("includes self-contained question rules", func(t *testing.T) {
		params := GenerationParams{
			NumQuestions:        3,
			Difficulty:          "easy",
			MultipleChoiceCount: 0,
		}

		prompt := BuildSystemPrompt(params)

		require.Contains(t, prompt, "САМОДОСТАТОЧНЫМ")
		require.Contains(t, prompt, "В примере выше")
	})

	t.Run("includes multiple_choice JSON example when count > 0 and ProfileData is set", func(t *testing.T) {
		params := GenerationParams{
			NumQuestions:        5,
			Difficulty:          "medium",
			MultipleChoiceCount: 2,
			ProfileData: &ProfileData{
				QuestionExamples: []QuestionExample{
					{QuestionType: SingleChoice, QuestionText: "Вопрос single", Answers: []ExampleAnswer{{Text: "А", IsCorrect: true}}, Explanation: "Пояснение"},
					{QuestionType: MultipleChoice, QuestionText: "Вопрос multiple", Answers: []ExampleAnswer{{Text: "А", IsCorrect: true}, {Text: "Б", IsCorrect: false}}, Explanation: "Пояснение"},
				},
			},
		}

		prompt := BuildSystemPrompt(params)

		require.Contains(t, prompt, `"type": "multiple_choice"`)
		require.Contains(t, prompt, `"type": "single_choice"`)
	})

	t.Run("excludes multiple_choice JSON example when count is 0 and ProfileData is set", func(t *testing.T) {
		params := GenerationParams{
			NumQuestions:        5,
			Difficulty:          "easy",
			MultipleChoiceCount: 0,
			ProfileData: &ProfileData{
				QuestionExamples: []QuestionExample{
					{QuestionType: SingleChoice, QuestionText: "Вопрос single", Answers: []ExampleAnswer{{Text: "А", IsCorrect: true}}, Explanation: "Пояснение"},
					{QuestionType: MultipleChoice, QuestionText: "Вопрос multiple", Answers: []ExampleAnswer{{Text: "А", IsCorrect: true}}, Explanation: "Пояснение"},
				},
			},
		}

		prompt := BuildSystemPrompt(params)

		require.NotContains(t, prompt, `"type": "multiple_choice"`)
		require.Contains(t, prompt, `"type": "single_choice"`)
	})

	t.Run("no ProfileData produces empty examples section", func(t *testing.T) {
		params := GenerationParams{
			NumQuestions:        5,
			Difficulty:          "easy",
			MultipleChoiceCount: 2,
		}

		prompt := BuildSystemPrompt(params)

		require.NotContains(t, prompt, `"type": "single_choice"`)
		require.NotContains(t, prompt, `"type": "multiple_choice"`)
	})
}

func TestBuildSystemPrompt_AcademicProfiles(t *testing.T) {
	// Helper: minimal ProfileData with the instruction text and one example per type.
	makeProfile := func(code AcademicProfile, instruction string, singleText, multipleText string) *ProfileData {
		examples := []QuestionExample{
			{
				QuestionType: SingleChoice,
				QuestionText: singleText,
				Answers:      []ExampleAnswer{{Text: "А", IsCorrect: true}, {Text: "Б", IsCorrect: false}},
				Explanation:  "Пояснение",
			},
		}
		if multipleText != "" {
			examples = append(examples, QuestionExample{
				QuestionType: MultipleChoice,
				QuestionText: multipleText,
				Answers:      []ExampleAnswer{{Text: "А", IsCorrect: true}, {Text: "Б", IsCorrect: false}},
				Explanation:  "Пояснение",
			})
		}
		return &ProfileData{
			Code:        code,
			Instruction: instruction,
			Formulations: []FormulationExample{
				{Bad: "ПЛОХО для " + string(code), Good: "ХОРОШО для " + string(code)},
			},
			QuestionExamples: examples,
		}
	}

	t.Run("no ProfileData produces no profile block", func(t *testing.T) {
		params := GenerationParams{
			NumQuestions:        5,
			Difficulty:          "medium",
			MultipleChoiceCount: 0,
		}

		prompt := BuildSystemPrompt(params)

		require.NotContains(t, prompt, "ПРОФИЛЬ ДИСЦИПЛИНЫ")
	})

	t.Run("ProfileData with empty instruction produces no profile block", func(t *testing.T) {
		params := GenerationParams{
			NumQuestions:        5,
			Difficulty:          "medium",
			MultipleChoiceCount: 0,
			ProfileData:         &ProfileData{Code: ProfileUniversal},
		}

		prompt := BuildSystemPrompt(params)

		require.NotContains(t, prompt, "ПРОФИЛЬ ДИСЦИПЛИНЫ")
	})

	t.Run("technological ProfileData includes technical domain instruction", func(t *testing.T) {
		instruction := "ПРОФИЛЬ ДИСЦИПЛИНЫ: технологический (программирование, инженерия, математика, информатика).\nАкцентируй вопросы на точных определениях, алгоритмах, формулах."
		params := GenerationParams{
			NumQuestions:        5,
			Difficulty:          "medium",
			MultipleChoiceCount: 0,
			ProfileData:         makeProfile(ProfileTechnological, instruction, "Вопрос про алгоритм", ""),
		}

		prompt := BuildSystemPrompt(params)

		require.Contains(t, prompt, "ПРОФИЛЬ ДИСЦИПЛИНЫ")
		require.Contains(t, prompt, "технологический")
		require.Contains(t, prompt, "алгоритм")
	})

	t.Run("natural_science ProfileData includes science domain instruction", func(t *testing.T) {
		instruction := "ПРОФИЛЬ ДИСЦИПЛИНЫ: естественно-научный (физика, химия, биология).\nАкцентируй вопросы на законах природы."
		params := GenerationParams{
			NumQuestions:        5,
			Difficulty:          "easy",
			MultipleChoiceCount: 0,
			ProfileData:         makeProfile(ProfileNaturalScience, instruction, "Закон Бойля–Мариотта", ""),
		}

		prompt := BuildSystemPrompt(params)

		require.Contains(t, prompt, "ПРОФИЛЬ ДИСЦИПЛИНЫ")
		require.Contains(t, prompt, "естественно-научный")
		require.Contains(t, prompt, "законах природы")
	})

	t.Run("humanities ProfileData includes humanities domain instruction", func(t *testing.T) {
		instruction := "ПРОФИЛЬ ДИСЦИПЛИНЫ: гуманитарный (история, литература, философия).\nАкцентируй вопросы на авторстве, периодизации."
		params := GenerationParams{
			NumQuestions:        5,
			Difficulty:          "hard",
			MultipleChoiceCount: 3,
			ProfileData:         makeProfile(ProfileHumanities, instruction, "Достоевский", "Эпоха Просвещения"),
		}

		prompt := BuildSystemPrompt(params)

		require.Contains(t, prompt, "ПРОФИЛЬ ДИСЦИПЛИНЫ")
		require.Contains(t, prompt, "гуманитарный")
		require.Contains(t, prompt, "периодизации")
	})

	t.Run("social_economic ProfileData includes social-economic domain instruction", func(t *testing.T) {
		instruction := "ПРОФИЛЬ ДИСЦИПЛИНЫ: социально-экономический.\nАкцентируй вопросы на нормативной базы."
		params := GenerationParams{
			NumQuestions:        5,
			Difficulty:          "medium",
			MultipleChoiceCount: 0,
			ProfileData:         makeProfile(ProfileSocialEconomic, instruction, "Вопрос про ВВП", ""),
		}

		prompt := BuildSystemPrompt(params)

		require.Contains(t, prompt, "ПРОФИЛЬ ДИСЦИПЛИНЫ")
		require.Contains(t, prompt, "социально-экономический")
		require.Contains(t, prompt, "нормативной базы")
	})

	t.Run("creative ProfileData includes creative domain instruction", func(t *testing.T) {
		instruction := "ПРОФИЛЬ ДИСЦИПЛИНЫ: творческий (изобразительное искусство, дизайн).\nАкцентируй вопросы на художественных стилях."
		params := GenerationParams{
			NumQuestions:        5,
			Difficulty:          "easy",
			MultipleChoiceCount: 0,
			ProfileData:         makeProfile(ProfileCreative, instruction, "Импрессионизм", ""),
		}

		prompt := BuildSystemPrompt(params)

		require.Contains(t, prompt, "ПРОФИЛЬ ДИСЦИПЛИНЫ")
		require.Contains(t, prompt, "творческий")
		require.Contains(t, prompt, "художественных стилях")
	})

	t.Run("ProfileData question example text appears in prompt", func(t *testing.T) {
		params := GenerationParams{
			NumQuestions:        5,
			Difficulty:          "medium",
			MultipleChoiceCount: 0,
			ProfileData:         makeProfile(ProfileNaturalScience, "Блок профиля", "Закон Бойля–Мариотта", ""),
		}

		prompt := BuildSystemPrompt(params)

		require.Contains(t, prompt, "Закон Бойля–Мариотта")
	})

	t.Run("ProfileData with multiple_choice example appears when count > 0", func(t *testing.T) {
		params := GenerationParams{
			NumQuestions:        5,
			Difficulty:          "medium",
			MultipleChoiceCount: 2,
			ProfileData:         makeProfile(ProfileHumanities, "Блок профиля", "Достоевский", "Черты эпохи Просвещения"),
		}

		prompt := BuildSystemPrompt(params)

		require.Contains(t, prompt, "Достоевский")
		require.Contains(t, prompt, "Черты эпохи Просвещения")
	})

	t.Run("difficulty ratio counts are preserved regardless of ProfileData", func(t *testing.T) {
		params := GenerationParams{
			NumQuestions:        10,
			Difficulty:          "hard",
			MultipleChoiceCount: 6,
			ProfileData:         makeProfile(ProfileSocialEconomic, "ПРОФИЛЬ ДИСЦИПЛИНЫ: социально-экономический.", "Вопрос", "МногоВыборВопрос"),
		}

		prompt := BuildSystemPrompt(params)

		require.Contains(t, prompt, "4 вопросов типа single_choice")
		require.Contains(t, prompt, "6 вопросов типа multiple_choice")
		require.Contains(t, prompt, "социально-экономический")
	})

	t.Run("ProfileData formulation examples appear in prompt", func(t *testing.T) {
		params := GenerationParams{
			NumQuestions:        5,
			Difficulty:          "medium",
			MultipleChoiceCount: 0,
			ProfileData: &ProfileData{
				Code:        ProfileTechnological,
				Instruction: "Блок профиля",
				Formulations: []FormulationExample{
					{Bad: "Плохой вопрос про полиморфизм", Good: "Хороший вопрос про полиморфизм в ООП"},
				},
			},
		}

		prompt := BuildSystemPrompt(params)

		require.Contains(t, prompt, "ПРИМЕРЫ ПРАВИЛЬНОЙ И НЕПРАВИЛЬНОЙ ФОРМУЛИРОВКИ")
		require.Contains(t, prompt, "Плохой вопрос про полиморфизм")
		require.Contains(t, prompt, "Хороший вопрос про полиморфизм в ООП")
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

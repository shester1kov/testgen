package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// YandexGPTStrategy implements LLM strategy for YandexGPT API
type YandexGPTStrategy struct {
	apiKey   string
	folderID string
	model    string
	baseURL  string // Base URL for API (for testing)
	client   *http.Client
}

// YandexGPTRequest represents the request structure for YandexGPT API
type YandexGPTRequest struct {
	ModelURI          string                   `json:"modelUri"`
	CompletionOptions YandexCompletionOptions  `json:"completionOptions"`
	Messages          []YandexMessage          `json:"messages"`
}

// YandexCompletionOptions represents completion options
type YandexCompletionOptions struct {
	Stream      bool    `json:"stream"`
	Temperature float64 `json:"temperature"`
	MaxTokens   string  `json:"maxTokens"` // YandexGPT API requires string, not int
}

// YandexMessage represents a message in the conversation
type YandexMessage struct {
	Role string `json:"role"` // system, user, assistant
	Text string `json:"text"`
}

// YandexGPTResponse represents the response from YandexGPT API
type YandexGPTResponse struct {
	Result YandexResult `json:"result"`
}

// YandexResult contains the generated result
type YandexResult struct {
	Alternatives []YandexAlternative `json:"alternatives"`
	Usage        YandexUsage         `json:"usage"`
	ModelVersion string              `json:"modelVersion"`
}

// YandexAlternative represents one generated alternative
type YandexAlternative struct {
	Message YandexMessage `json:"message"`
	Status  string        `json:"status"` // ALTERNATIVE_STATUS_FINAL, etc.
}

// YandexUsage tracks token usage
// Note: YandexGPT API returns tokens as strings, not integers
type YandexUsage struct {
	InputTextTokens  string `json:"inputTextTokens"`
	CompletionTokens string `json:"completionTokens"`
	TotalTokens      string `json:"totalTokens"`
}

// QuestionResponse represents the structured JSON response from LLM
type QuestionResponse struct {
	Questions []struct {
		Question   string `json:"question"`
		Type       string `json:"type"`
		Difficulty string `json:"difficulty"`
		Answers    []struct {
			Text      string `json:"text"`
			IsCorrect bool   `json:"is_correct"`
		} `json:"answers"`
		Explanation string `json:"explanation,omitempty"`
	} `json:"questions"`
}

// NewYandexGPTStrategy creates a new YandexGPT strategy
func NewYandexGPTStrategy(apiKey, folderID, model string) *YandexGPTStrategy {
	if model == "" {
		model = "yandexgpt-lite" // Default to lite model (cheaper, faster)
	}

	return &YandexGPTStrategy{
		apiKey:   apiKey,
		folderID: folderID,
		model:    model,
		baseURL:  "https://llm.api.cloud.yandex.net/foundationModels/v1/completion",
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// GenerateQuestions generates questions using YandexGPT API
func (s *YandexGPTStrategy) GenerateQuestions(ctx context.Context, params GenerationParams) ([]GeneratedQuestion, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("yandexgpt API key not configured")
	}
	if s.folderID == "" {
		return nil, fmt.Errorf("yandexgpt folder ID not configured")
	}

	// Build the prompt
	prompt := s.buildPrompt(params)

	// Prepare the request
	reqBody := YandexGPTRequest{
		ModelURI: fmt.Sprintf("gpt://%s/%s", s.folderID, s.model),
		CompletionOptions: YandexCompletionOptions{
			Stream:      false,
			Temperature: 0.6, // Match Python example
			MaxTokens:   "2000",
		},
		Messages: []YandexMessage{
			{
				Role: "system",
				Text: "Ты - профессиональный создатель тестовых вопросов для образовательных целей. Генерируй качественные вопросы на русском языке в формате JSON.",
			},
			{
				Role: "user",
				Text: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST",
		s.baseURL,
		bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Api-Key %s", s.apiKey))

	// Log request for debugging

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()


	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}


	// Log first 500 chars of response for debugging
	bodyPreview := string(body)
	if len(bodyPreview) > 500 {
		bodyPreview = bodyPreview[:500] + "..."
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yandexgpt API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse YandexGPT response
	var yandexResp YandexGPTResponse
	if err := json.Unmarshal(body, &yandexResp); err != nil {
		return nil, fmt.Errorf("failed to parse yandex response: %w", err)
	}


	// Extract the generated text
	if len(yandexResp.Result.Alternatives) == 0 {
		return nil, fmt.Errorf("no alternatives in response")
	}

	generatedText := yandexResp.Result.Alternatives[0].Message.Text

	// Log first 300 chars of generated text
	textPreview := generatedText
	if len(textPreview) > 300 {
		textPreview = textPreview[:300] + "..."
	}

	// Parse the JSON from the generated text
	questions, err := s.parseQuestions(generatedText, params)
	if err != nil {
		return nil, fmt.Errorf("failed to parse generated questions: %w", err)
	}

	return questions, nil
}

// buildPrompt creates a prompt for question generation
func (s *YandexGPTStrategy) buildPrompt(params GenerationParams) string {
	var questionTypesStr string
	if len(params.QuestionTypes) > 0 {
		types := make([]string, len(params.QuestionTypes))
		for i, qt := range params.QuestionTypes {
			types[i] = string(qt)
		}
		questionTypesStr = strings.Join(types, ", ")
	} else {
		questionTypesStr = string(SingleChoice)
	}

	difficulty := params.Difficulty
	if difficulty == "" {
		difficulty = "medium"
	}

	language := params.Language
	if language == "" {
		language = "ru"
	}

	prompt := fmt.Sprintf(`На основе следующего текста создай %d тестовых вопросов.

ТЕКСТ:
%s

ТРЕБОВАНИЯ:
- Типы вопросов: %s
- Сложность: %s
- Язык: %s
- Для каждого вопроса типа single_choice создай 4 варианта ответа (1 правильный, 3 неправильных)
- Для каждого вопроса типа multiple_choice создай 5-6 вариантов (2-3 правильных, 2-3 неправильных)
- Для true_false создай только 2 варианта: "Верно" и "Неверно"

ВАЖНО - ПРАВИЛА ФОРМУЛИРОВКИ ВОПРОСОВ:
1. Каждый вопрос должен быть САМОДОСТАТОЧНЫМ и понятным без ссылок на текст
2. НЕ используй фразы типа "В примере выше", "Как показано в коде", "Согласно тексту лекции"
3. Если в тексте есть конкретный пример кода или ситуации - включи его ПОЛНОСТЬЮ в текст вопроса
4. Вопрос должен содержать всю необходимую информацию для ответа
5. Формулируй вопросы в общем виде, проверяя понимание концепций, а не запоминание примеров

ПРИМЕРЫ:
❌ ПЛОХО: "В приведённом выше примере наследования, какой метод будет вызван?"
✅ ХОРОШО: "В следующем коде:\nclass Parent { void foo() {...} }\nclass Child extends Parent { void foo() {...} }\nChild obj = new Child();\nКакой метод будет вызван при obj.foo()?"

❌ ПЛОХО: "Согласно лекции, что такое полиморфизм?"
✅ ХОРОШО: "Что такое полиморфизм в объектно-ориентированном программировании?"

ФОРМАТ ОТВЕТА (строго JSON):
{
  "questions": [
    {
      "question": "Текст вопроса",
      "type": "single_choice",
      "difficulty": "%s",
      "answers": [
        {"text": "Вариант ответа 1", "is_correct": true},
        {"text": "Вариант ответа 2", "is_correct": false},
        {"text": "Вариант ответа 3", "is_correct": false},
        {"text": "Вариант ответа 4", "is_correct": false}
      ],
      "explanation": "Краткое объяснение правильного ответа"
    }
  ]
}

Верни ТОЛЬКО валидный JSON без дополнительного текста.`,
		params.NumQuestions,
		params.Text,
		questionTypesStr,
		difficulty,
		language,
		difficulty,
	)

	return prompt
}

// parseQuestions parses the generated JSON into GeneratedQuestion structs
func (s *YandexGPTStrategy) parseQuestions(text string, params GenerationParams) ([]GeneratedQuestion, error) {

	// Sometimes LLM wraps JSON in markdown code blocks
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)


	var qResponse QuestionResponse
	if err := json.Unmarshal([]byte(text), &qResponse); err != nil {
		// Log first 500 chars of problematic text
		preview := text
		if len(preview) > 500 {
			preview = preview[:500] + "..."
		}
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}


	if len(qResponse.Questions) == 0 {
		return nil, fmt.Errorf("no questions generated")
	}

	result := make([]GeneratedQuestion, 0, len(qResponse.Questions))
	for _, q := range qResponse.Questions {
		answers := make([]GeneratedAnswer, len(q.Answers))
		for i, a := range q.Answers {
			answers[i] = GeneratedAnswer{
				Text:      a.Text,
				IsCorrect: a.IsCorrect,
			}
		}

		result = append(result, GeneratedQuestion{
			QuestionText: q.Question,
			QuestionType: QuestionType(q.Type),
			Difficulty:   q.Difficulty,
			Answers:      answers,
			Explanation:  q.Explanation,
		})
	}

	return result, nil
}

// GetProviderName returns the provider name
func (s *YandexGPTStrategy) GetProviderName() string {
	return "yandexgpt"
}

package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	log.Printf("[YandexGPT] Step 1: Preparing to send request")
	log.Printf("[YandexGPT] Model URI: gpt://%s/%s", s.folderID, s.model)
	log.Printf("[YandexGPT] Text length: %d chars", len(params.Text))
	log.Printf("[YandexGPT] Questions: %d", params.NumQuestions)
	log.Printf("[YandexGPT] Difficulty: %s", params.Difficulty)
	log.Printf("[YandexGPT] Prompt length: %d chars", len(prompt))

	// Send request
	log.Printf("[YandexGPT] Step 2: Calling YandexGPT API...")
	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("[YandexGPT] ERROR in Step 2: HTTP request failed: %v", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("[YandexGPT] Step 3: Response received successfully")
	log.Printf("[YandexGPT] Status code: %d", resp.StatusCode)
	log.Printf("[YandexGPT] Content-Type: %s", resp.Header.Get("Content-Type"))

	// Read response body
	log.Printf("[YandexGPT] Step 4: Reading response body...")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[YandexGPT] ERROR in Step 4: Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("[YandexGPT] Step 5: Response body read successfully")
	log.Printf("[YandexGPT] Body length: %d bytes", len(body))

	// Log first 500 chars of response for debugging
	bodyPreview := string(body)
	if len(bodyPreview) > 500 {
		bodyPreview = bodyPreview[:500] + "..."
	}
	log.Printf("[YandexGPT] Body preview: %s", bodyPreview)

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("[YandexGPT] ERROR: Non-OK status code")
		log.Printf("[YandexGPT] Status: %d", resp.StatusCode)
		log.Printf("[YandexGPT] Full response: %s", string(body))
		return nil, fmt.Errorf("yandexgpt API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse YandexGPT response
	log.Printf("[YandexGPT] Step 6: Parsing YandexGPT JSON response...")
	var yandexResp YandexGPTResponse
	if err := json.Unmarshal(body, &yandexResp); err != nil {
		log.Printf("[YandexGPT] ERROR in Step 6: Failed to unmarshal YandexGPT response: %v", err)
		log.Printf("[YandexGPT] Response body: %s", string(body))
		return nil, fmt.Errorf("failed to parse yandex response: %w", err)
	}

	log.Printf("[YandexGPT] Step 7: YandexGPT response parsed successfully")
	log.Printf("[YandexGPT] Number of alternatives: %d", len(yandexResp.Result.Alternatives))

	// Extract the generated text
	if len(yandexResp.Result.Alternatives) == 0 {
		log.Printf("[YandexGPT] ERROR: No alternatives in response")
		return nil, fmt.Errorf("no alternatives in response")
	}

	generatedText := yandexResp.Result.Alternatives[0].Message.Text
	log.Printf("[YandexGPT] Step 8: Extracted generated text")
	log.Printf("[YandexGPT] Generated text length: %d chars", len(generatedText))

	// Log first 300 chars of generated text
	textPreview := generatedText
	if len(textPreview) > 300 {
		textPreview = textPreview[:300] + "..."
	}
	log.Printf("[YandexGPT] Generated text preview: %s", textPreview)

	// Parse the JSON from the generated text
	log.Printf("[YandexGPT] Step 9: Parsing questions from generated text...")
	questions, err := s.parseQuestions(generatedText, params)
	if err != nil {
		log.Printf("[YandexGPT] ERROR in Step 9: Failed to parse questions: %v", err)
		return nil, fmt.Errorf("failed to parse generated questions: %w", err)
	}

	log.Printf("[YandexGPT] Step 10: SUCCESS - Generated %d questions", len(questions))
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
	log.Printf("[YandexGPT] parseQuestions: Starting to parse generated text")
	log.Printf("[YandexGPT] parseQuestions: Original text length: %d", len(text))

	// Sometimes LLM wraps JSON in markdown code blocks
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	log.Printf("[YandexGPT] parseQuestions: After cleanup, text length: %d", len(text))

	var qResponse QuestionResponse
	if err := json.Unmarshal([]byte(text), &qResponse); err != nil {
		log.Printf("[YandexGPT] parseQuestions: ERROR - Failed to unmarshal JSON: %v", err)
		// Log first 500 chars of problematic text
		preview := text
		if len(preview) > 500 {
			preview = preview[:500] + "..."
		}
		log.Printf("[YandexGPT] parseQuestions: Problematic text: %s", preview)
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	log.Printf("[YandexGPT] parseQuestions: Successfully parsed JSON")
	log.Printf("[YandexGPT] parseQuestions: Number of questions in response: %d", len(qResponse.Questions))

	if len(qResponse.Questions) == 0 {
		log.Printf("[YandexGPT] parseQuestions: ERROR - No questions in parsed response")
		return nil, fmt.Errorf("no questions generated")
	}

	result := make([]GeneratedQuestion, 0, len(qResponse.Questions))
	for i, q := range qResponse.Questions {
		log.Printf("[YandexGPT] parseQuestions: Processing question %d: type=%s, answers=%d",
			i+1, q.Type, len(q.Answers))

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

	log.Printf("[YandexGPT] parseQuestions: Successfully created %d GeneratedQuestion structs", len(result))
	return result, nil
}

// GetProviderName returns the provider name
func (s *YandexGPTStrategy) GetProviderName() string {
	return "yandexgpt"
}

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

	"go.uber.org/zap"
)

// YandexGPTStrategy implements LLM strategy for YandexGPT API
type YandexGPTStrategy struct {
	apiKey   string
	folderID string
	model    string
	baseURL  string // Base URL for API (for testing)
	client   *http.Client
	logger   *zap.Logger
}

// YandexGPTRequest represents the request structure for YandexGPT API
type YandexGPTRequest struct {
	ModelURI          string                  `json:"modelUri"`
	CompletionOptions YandexCompletionOptions `json:"completionOptions"`
	Messages          []YandexMessage         `json:"messages"`
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

// NewYandexGPTStrategy creates a new YandexGPT strategy.
// The logger parameter is optional; if omitted, a no-op logger is used.
func NewYandexGPTStrategy(apiKey, folderID, model string, logger ...*zap.Logger) *YandexGPTStrategy {
	if model == "" {
		model = "yandexgpt-lite" // Default to lite model (cheaper, faster)
	}

	var l *zap.Logger
	if len(logger) > 0 && logger[0] != nil {
		l = logger[0]
	} else {
		l = zap.NewNop()
	}

	return &YandexGPTStrategy{
		apiKey:   apiKey,
		folderID: folderID,
		model:    model,
		baseURL:  "https://llm.api.cloud.yandex.net/foundationModels/v1/completion",
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		logger: l,
	}
}

// GenerateQuestions generates questions using YandexGPT API
func (s *YandexGPTStrategy) GenerateQuestions(ctx context.Context, params GenerationParams) ([]GeneratedQuestion, error) {
	if s.apiKey == "" {
		s.logger.Error("yandexgpt API key not configured")
		return nil, fmt.Errorf("yandexgpt API key not configured")
	}
	if s.folderID == "" {
		s.logger.Error("yandexgpt folder ID not configured")
		return nil, fmt.Errorf("yandexgpt folder ID not configured")
	}

	temperature := params.Temperature
	if temperature == 0 {
		temperature = 0.5 // safe default when caller does not set profile-based temperature
	}

	// Prepare the request
	reqBody := YandexGPTRequest{
		ModelURI: fmt.Sprintf("gpt://%s/%s", s.folderID, s.model),
		CompletionOptions: YandexCompletionOptions{
			Stream:      false,
			Temperature: temperature,
			MaxTokens:   "8000",
		},
		Messages: []YandexMessage{
			{
				Role: "system",
				Text: BuildSystemPrompt(params),
			},
			{
				Role: "user",
				Text: params.Text,
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

	s.logger.Debug("sending request to yandexgpt API",
		zap.String("model", s.model),
		zap.Int("num_questions", params.NumQuestions),
		zap.String("difficulty", params.Difficulty),
		zap.String("profile", string(params.Profile)),
	)

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

	// Check status code
	if resp.StatusCode != http.StatusOK {
		s.logger.Error("yandexgpt API returned non-200 status",
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(body)),
		)
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

	s.logger.Debug("received response from yandexgpt",
		zap.Int("text_length", len(generatedText)),
	)

	// Parse the JSON from the generated text
	questions, err := s.parseQuestions(generatedText, params)
	if err != nil {
		return nil, fmt.Errorf("failed to parse generated questions: %w", err)
	}

	return questions, nil
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
		preview := text
		if len(preview) > 500 {
			preview = preview[:500] + "..."
		}
		s.logger.Error("failed to parse JSON from yandexgpt response",
			zap.Error(err),
			zap.String("response_preview", preview),
		)
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

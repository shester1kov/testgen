package moodle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client represents a Moodle web service client
type Client struct {
	baseURL     string
	token       string
	importToken string
	httpClient  *http.Client
}

// NewClient creates a new Moodle client
func NewClient(baseURL, token, importToken string) *Client {
	return &Client{
		baseURL:     baseURL,
		token:       token,
		importToken: importToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// UploadQuizRequest represents quiz upload request
type UploadQuizRequest struct {
	CourseID   string `json:"courseid"`
	CourseName string `json:"-"` // Not used in new implementation
	QuizName   string `json:"quizname"`
	XMLContent string `json:"xmlcontent"`
}

// UploadQuizResponse represents quiz upload response
type UploadQuizResponse struct {
	QuizID           string `json:"quiz_id"` // Empty for question-only import
	CourseID         string `json:"course_id"`
	CategoryID       string `json:"category_id"`
	CategoryName     string `json:"category_name"`
	QuestionsImported int   `json:"questions_imported"`
	QuestionBankURL  string `json:"question_bank_url"`
	Note             string `json:"note"`
	Success          bool   `json:"success"`
	Message          string `json:"message"`
}

// UploadQuiz uploads a quiz to Moodle using custom import endpoint
func (c *Client) UploadQuiz(ctx context.Context, req UploadQuizRequest) (*UploadQuizResponse, error) {
	// Prepare JSON request body
	requestBody := map[string]string{
		"token":      c.importToken,
		"courseid":   req.CourseID,
		"quizname":   req.QuizName,
		"xmlcontent": req.XMLContent,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Build request URL to custom import endpoint
	requestURL := fmt.Sprintf("%s/local/testgen_import.php", c.baseURL)

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moodle API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var result UploadQuizResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("moodle import failed: %s", result.Message)
	}

	return &result, nil
}

// GetCourses retrieves list of courses from Moodle
func (c *Client) GetCourses(ctx context.Context) ([]Course, error) {
	params := url.Values{}
	params.Set("wstoken", c.token)
	params.Set("wsfunction", "core_course_get_courses")
	params.Set("moodlewsrestformat", "json")

	requestURL := fmt.Sprintf("%s/webservice/rest/server.php?%s", c.baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moodle API returned status %d", resp.StatusCode)
	}

	var courses []Course
	if err := json.NewDecoder(resp.Body).Decode(&courses); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return courses, nil
}

// Course represents a Moodle course
type Course struct {
	ID        int    `json:"id"`
	ShortName string `json:"shortname"`
	FullName  string `json:"fullname"`
	Category  int    `json:"categoryid"`
}

// ValidateConnection checks if the Moodle connection is valid
func (c *Client) ValidateConnection(ctx context.Context) error {
	params := url.Values{}
	params.Set("wstoken", c.token)
	params.Set("wsfunction", "core_webservice_get_site_info")
	params.Set("moodlewsrestformat", "json")

	requestURL := fmt.Sprintf("%s/webservice/rest/server.php?%s", c.baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("moodle API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var siteInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&siteInfo); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for error in response
	if errMsg, ok := siteInfo["error"].(string); ok {
		return fmt.Errorf("moodle error: %s", errMsg)
	}

	return nil
}

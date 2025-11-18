package moodle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

// Client represents a Moodle web service client
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewClient creates a new Moodle client
func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// UploadQuizRequest represents quiz upload request
type UploadQuizRequest struct {
	CourseName string
	QuizName   string
	XMLContent string
}

// UploadQuizResponse represents quiz upload response
type UploadQuizResponse struct {
	QuizID   string `json:"quiz_id"`
	CourseID string `json:"course_id"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
}

// UploadQuiz uploads a quiz to Moodle
// TODO: Implement actual Moodle web service integration
// This is a placeholder implementation
func (c *Client) UploadQuiz(ctx context.Context, req UploadQuizRequest) (*UploadQuizResponse, error) {
	// Prepare multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add token
	if err := writer.WriteField("wstoken", c.token); err != nil {
		return nil, fmt.Errorf("failed to write token field: %w", err)
	}

	// Add web service function
	if err := writer.WriteField("wsfunction", "core_course_import_course"); err != nil {
		return nil, fmt.Errorf("failed to write wsfunction field: %w", err)
	}

	// Add format
	if err := writer.WriteField("moodlewsrestformat", "json"); err != nil {
		return nil, fmt.Errorf("failed to write format field: %w", err)
	}

	// Add quiz XML content as file
	part, err := writer.CreateFormFile("file", "quiz.xml")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.WriteString(part, req.XMLContent); err != nil {
		return nil, fmt.Errorf("failed to write XML content: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Build request URL
	requestURL := fmt.Sprintf("%s/webservice/rest/server.php", c.baseURL)

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("moodle API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var result UploadQuizResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
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

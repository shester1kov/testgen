package dto

// CreateTestRequest represents test creation request
type CreateTestRequest struct {
	Title       string  `json:"title" validate:"required,min=3"`
	Description string  `json:"description"`
	DocumentID  *string `json:"document_id" validate:"omitempty,uuid"`
}

// GenerateTestRequest represents test generation request
type GenerateTestRequest struct {
	DocumentID    string   `json:"document_id" validate:"required,uuid"`
	Title         string   `json:"title" validate:"required,min=3"`
	NumQuestions  int      `json:"num_questions" validate:"required,min=1,max=50"`
	QuestionTypes []string `json:"question_types"`
	Difficulty    string   `json:"difficulty" validate:"required,oneof=easy medium hard"`
	LLMProvider   string   `json:"llm_provider" validate:"omitempty,oneof=perplexity openai yandexgpt"`
}

// TestResponse represents test response
type TestResponse struct {
	ID             string          `json:"id"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	TotalQuestions int             `json:"total_questions"`
	Status         string          `json:"status"`
	MoodleSynced   bool            `json:"moodle_synced"`
	CreatedAt      string          `json:"created_at"`
	Questions      []QuestionDTO   `json:"questions,omitempty"`
}

// QuestionDTO represents question data
type QuestionDTO struct {
	ID           string      `json:"id"`
	QuestionText string      `json:"question_text"`
	QuestionType string      `json:"question_type"`
	Difficulty   string      `json:"difficulty"`
	Points       float64     `json:"points"`
	OrderNum     int         `json:"order_num"`
	Answers      []AnswerDTO `json:"answers"`
}

// AnswerDTO represents answer data
type AnswerDTO struct {
	ID         string `json:"id"`
	AnswerText string `json:"answer_text"`
	IsCorrect  bool   `json:"is_correct"`
	OrderNum   int    `json:"order_num"`
}

// TestListResponse represents list of tests
type TestListResponse struct {
	Tests    []TestResponse `json:"tests"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// GenerateQuestionsResponse represents the response after generating questions
type GenerateQuestionsResponse struct {
	Message   string        `json:"message"`
	Count     int           `json:"count"`
	Questions []QuestionDTO `json:"questions"`
}

// SyncMoodleRequest represents Moodle sync request
type SyncMoodleRequest struct {
	CourseName string `json:"course_name" validate:"required"`
}

// MoodleSyncResponse represents the response after syncing with Moodle
type MoodleSyncResponse struct {
	Message  string `json:"message"`
	MoodleID string `json:"moodle_id"`
	CourseID string `json:"course_id"`
}

// MoodleCoursesResponse represents the response for getting Moodle courses
type MoodleCoursesResponse struct {
	Courses []MoodleCourse `json:"courses"`
}

// MoodleCourse represents a Moodle course
type MoodleCourse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ShortName string `json:"short_name,omitempty"`
}

// MoodleConnectionResponse represents the response for Moodle connection validation
type MoodleConnectionResponse struct {
	Connected bool   `json:"connected"`
	Message   string `json:"message,omitempty"`
	Error     string `json:"error,omitempty"`
}

package dto

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Code    string `json:"code" example:"VALIDATION_ERROR"`
	Message string `json:"message" example:"Invalid input data"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Operation completed successfully"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// Error codes constants
const (
	// Authentication errors
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeTokenExpired       = "TOKEN_EXPIRED"
	ErrCodeInvalidToken       = "INVALID_TOKEN"

	// Authorization errors
	ErrCodeForbidden         = "FORBIDDEN"
	ErrCodeInsufficientPerms = "INSUFFICIENT_PERMISSIONS"

	// Validation errors
	ErrCodeValidationError = "VALIDATION_ERROR"
	ErrCodeInvalidInput    = "INVALID_INPUT"
	ErrCodeInvalidUUID     = "INVALID_UUID"
	ErrCodeInvalidRole     = "INVALID_ROLE"

	// Resource errors
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeAlreadyExists = "ALREADY_EXISTS"
	ErrCodeConflict      = "CONFLICT"

	// Database errors
	ErrCodeDatabaseError = "DATABASE_ERROR"
	ErrCodeInternalError = "INTERNAL_ERROR"

	// User errors
	ErrCodeUserNotFound     = "USER_NOT_FOUND"
	ErrCodeUserExists       = "USER_ALREADY_EXISTS"
	ErrCodeInvalidPassword  = "INVALID_PASSWORD"
	ErrCodeRoleNotFound     = "ROLE_NOT_FOUND"
	ErrCodeCannotDeleteRole = "CANNOT_DELETE_ROLE"

	// Document errors
	ErrCodeDocumentNotFound  = "DOCUMENT_NOT_FOUND"
	ErrCodeInvalidFileType   = "INVALID_FILE_TYPE"
	ErrCodeFileTooLarge      = "FILE_TOO_LARGE"
	ErrCodeUploadFailed      = "UPLOAD_FAILED"
	ErrCodeParsingFailed     = "PARSING_FAILED"
	ErrCodeDocumentExists    = "DOCUMENT_ALREADY_EXISTS"
	ErrCodeDocumentInUse     = "DOCUMENT_IN_USE"
	ErrCodeInvalidDocumentID = "INVALID_DOCUMENT_ID"
	ErrCodeDocumentNotParsed = "DOCUMENT_NOT_PARSED"

	// Test errors
	ErrCodeTestNotFound        = "TEST_NOT_FOUND"
	ErrCodeInvalidTestID       = "INVALID_TEST_ID"
	ErrCodeGenerationFailed    = "GENERATION_FAILED"
	ErrCodeExportFailed        = "EXPORT_FAILED"
	ErrCodeInvalidProvider     = "INVALID_PROVIDER"
	ErrCodeTestHasNoQuestions  = "TEST_HAS_NO_QUESTIONS"
	ErrCodeMoodleSyncFailed    = "MOODLE_SYNC_FAILED"
	ErrCodeMoodleUploadFailed  = "MOODLE_UPLOAD_FAILED"
	ErrCodeMoodleNotConnected  = "MOODLE_NOT_CONNECTED"

	// Pagination errors
	ErrCodeInvalidLimit  = "INVALID_LIMIT"
	ErrCodeInvalidOffset = "INVALID_OFFSET"
)

// NewErrorResponse creates a new error response
func NewErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string) SuccessResponse {
	return SuccessResponse{
		Success: true,
		Message: message,
	}
}

// NewMessageResponse creates a new message response
func NewMessageResponse(message string) MessageResponse {
	return MessageResponse{
		Message: message,
	}
}

// DashboardStatsResponse represents dashboard statistics
type DashboardStatsResponse struct {
	DocumentsCount int64 `json:"documents_count" example:"15"`
	TestsCount     int64 `json:"tests_count" example:"8"`
	QuestionsCount int64 `json:"questions_count" example:"120"`
}

package domain

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrRoleNotFound         = errors.New("role not found")
	ErrInvalidRoleName      = errors.New("invalid role name")
	ErrDocumentNotFound     = errors.New("document not found")
	ErrTestNotFound         = errors.New("test not found")
	ErrQuestionNotFound     = errors.New("question not found")
	ErrAccessDenied         = errors.New("access denied")
	ErrInvalidCredentials   = errors.New("invalid email or password")
	ErrAlreadyExists        = errors.New("already exists")
	ErrFileSizeExceeds      = errors.New("file size exceeds limit")
	ErrUnsupportedFileType  = errors.New("unsupported file type")
	ErrNotParsedYet         = errors.New("not parsed yet")
	ErrAlreadyParsed        = errors.New("already parsed")
	ErrInvalidLLMProvider   = errors.New("invalid LLM provider")
	ErrFailedToGenerate     = errors.New("failed to generate")
	ErrNoQuestions          = errors.New("no questions")
	ErrUnsupportedExport    = errors.New("unsupported export format")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrTitleTooShort        = errors.New("title must be at least")
	ErrQuestionTextTooShort = errors.New("question text must be at least")
)

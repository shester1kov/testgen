package errors

import "fmt"

// AppError represents an application error
type AppError struct {
	Code    int
	Message string
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common error constructors
func BadRequest(message string) *AppError {
	return &AppError{
		Code:    400,
		Message: message,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		Code:    401,
		Message: message,
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		Code:    403,
		Message: message,
	}
}

func NotFound(message string) *AppError {
	return &AppError{
		Code:    404,
		Message: message,
	}
}

func Conflict(message string) *AppError {
	return &AppError{
		Code:    409,
		Message: message,
	}
}

func InternalServerError(message string, err error) *AppError {
	return &AppError{
		Code:    500,
		Message: message,
		Err:     err,
	}
}

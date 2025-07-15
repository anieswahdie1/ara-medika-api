package errors

import (
	"fmt"
	"net/http"
)

type APIError struct {
	Status    int    `json:"status"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
	Details   any    `json:"details,omitempty"`
}

// Implement error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("status=%d, code=%s, message=%s", e.Status, e.ErrorCode, e.Message)
}

// Predefined error constructors
func NewBadRequestError(code, message string, details any) *APIError {
	return &APIError{
		Status:    http.StatusBadRequest,
		ErrorCode: code,
		Message:   message,
		Details:   details,
	}
}

func NewUnauthorizedError(code, message string) *APIError {
	return &APIError{
		Status:    http.StatusUnauthorized,
		ErrorCode: code,
		Message:   message,
	}
}

func NewNotFoundError(code, message string) *APIError {
	return &APIError{
		Status:    http.StatusNotFound,
		ErrorCode: code,
		Message:   message,
	}
}

func NewInternalServerError(code, message string) *APIError {
	return &APIError{
		Status:    http.StatusInternalServerError,
		ErrorCode: code,
		Message:   message,
	}
}

func NewValidationError(errors map[string]string) *APIError {
	return NewBadRequestError(
		"validation_failed",
		"Validation failed",
		errors,
	)
}

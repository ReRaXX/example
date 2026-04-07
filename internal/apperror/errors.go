package apperror

import (
	"errors"
	"net/http"
)

// AppError represents a custom application error with HTTP status code
type AppError struct {
	Message    string
	StatusCode int
	Err        error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError
func NewAppError(message string, statusCode int, err error) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// Predefined error types
var (
	// ValidationError represents a validation error (400)
	ValidationError = func(message string, err error) *AppError {
		return NewAppError(message, http.StatusBadRequest, err)
	}

	// UnauthorizedError represents an unauthorized error (401)
	UnauthorizedError = func(message string, err error) *AppError {
		return NewAppError(message, http.StatusUnauthorized, err)
	}

	// ForbiddenError represents a forbidden error (403)
	ForbiddenError = func(message string, err error) *AppError {
		return NewAppError(message, http.StatusForbidden, err)
	}

	// NotFoundError represents a not found error (404)
	NotFoundError = func(message string, err error) *AppError {
		return NewAppError(message, http.StatusNotFound, err)
	}

	// ConflictError represents a conflict error (409)
	ConflictError = func(message string, err error) *AppError {
		return NewAppError(message, http.StatusConflict, err)
	}

	// InternalServerError represents an internal server error (500)
	InternalServerError = func(message string, err error) *AppError {
		return NewAppError(message, http.StatusInternalServerError, err)
	}
)

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetStatusCode returns the HTTP status code from an error
func GetStatusCode(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.StatusCode
	}
	return http.StatusInternalServerError
}

// pkg/errors/errors.go
package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Standard error types
var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrNotFound         = errors.New("resource not found")
	ErrInternalServer   = errors.New("internal server error")
	ErrDuplicateEntry   = errors.New("duplicate entry")
	ErrValidationFailed = errors.New("validation failed")
	ErrDatabase         = errors.New("database error")
	ErrExternalService  = errors.New("external service error")
)

// AppError extends the standard error with additional context
type AppError struct {
	Err     error  // Original error
	Message string // User-friendly message
	Code    string // Error code
	Status  int    // HTTP status code
}

// Error returns the error message
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Is reports whether target is in the error chain
func (e *AppError) Is(target error) bool {
	if e.Err == nil {
		return false
	}
	return errors.Is(e.Err, target)
}

// Helper functions to create specific error types
func NewBadRequestError(message string, err error) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    "BAD_REQUEST",
		Status:  http.StatusBadRequest,
	}
}

func NewUnauthorizedError(message string, err error) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    "UNAUTHORIZED",
		Status:  http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string, err error) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    "FORBIDDEN",
		Status:  http.StatusForbidden,
	}
}

func NewNotFoundError(message string, err error) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    "NOT_FOUND",
		Status:  http.StatusNotFound,
	}
}

func NewConflictError(message string, err error) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    "CONFLICT",
		Status:  http.StatusConflict,
	}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    "INTERNAL_SERVER_ERROR",
		Status:  http.StatusInternalServerError,
	}
}

func NewValidationError(message string, err error) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    "VALIDATION_ERROR",
		Status:  http.StatusBadRequest,
	}
}

func NewDatabaseError(message string, err error) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    "DATABASE_ERROR",
		Status:  http.StatusInternalServerError,
	}
}

func NewExternalServiceError(message string, err error) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    "EXTERNAL_SERVICE_ERROR",
		Status:  http.StatusInternalServerError,
	}
}

// WrapError adds context to an existing error
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	
	// If it's already an AppError, just update the message
	var appErr *AppError
	if errors.As(err, &appErr) {
		if message != "" {
			appErr.Message = message + ": " + appErr.Message
		}
		return appErr
	}
	
	// Otherwise wrap it as an internal error
	return &AppError{
		Err:     err,
		Message: message,
		Code:    "INTERNAL_SERVER_ERROR",
		Status:  http.StatusInternalServerError,
	}
}

// PostgreSQL error codes
const (
	PgErrUniqueViolation      = "23505"
	PgErrForeignKeyViolation  = "23503"
	PgErrCheckViolation       = "23514"
)

// IsPgError checks if an error is a PostgreSQL error with a specific code
func IsPgError(err error, code string) bool {
	pgErr, ok := err.(interface {
		Code() string
	})
	
	if ok && pgErr.Code() == code {
		return true
	}
	return false
}
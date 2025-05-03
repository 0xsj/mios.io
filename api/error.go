package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Existing error types
var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrNotFound         = errors.New("resource not found")
	ErrInternalServer   = errors.New("internal server error")
	ErrDuplicateEntry   = errors.New("duplicate entry")
	ErrValidationFailed = errors.New("validation failed")
	ErrDatabase         = errors.New("database error")
	ErrInternal         = errors.New("internal server error")
	ErrExternalService  = errors.New("external service error")
)

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Status  int    `json:"-"`                 // HTTP status code, not shown in response
	Code    string `json:"code"`              // Application-specific error code
	Message string `json:"message"`           // User-friendly error message
	Details any    `json:"details,omitempty"` // Optional details about the error
}

// Common error responses
var (
	ErrBadRequestResponse = ErrorResponse{
		Status:  http.StatusBadRequest,
		Code:    "BAD_REQUEST",
		Message: "The request was invalid",
	}

	ErrUnauthorizedResponse = ErrorResponse{
		Status:  http.StatusUnauthorized,
		Code:    "UNAUTHORIZED",
		Message: "Authentication is required",
	}

	ErrNotFoundResponse = ErrorResponse{
		Status:  http.StatusNotFound,
		Code:    "NOT_FOUND",
		Message: "The requested resource was not found",
	}

	ErrServerResponse = ErrorResponse{
		Status:  http.StatusInternalServerError,
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "An unexpected error occurred",
	}
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

// Helper functions to create specific error types
func NewBadRequestError(message string, err error) error {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    "BAD_REQUEST",
		Status:  http.StatusBadRequest,
	}
}

func NewUnauthorizedError(message string, err error) error {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    "UNAUTHORIZED",
		Status:  http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string, err error) error {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    "FORBIDDEN",
		Status:  http.StatusForbidden,
	}
}

func NewNotFoundError(message string, err error) error {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    "NOT_FOUND",
		Status:  http.StatusNotFound,
	}
}

func NewConflictError(message string, err error) error {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    "CONFLICT",
		Status:  http.StatusConflict,
	}
}

func NewInternalError(message string, err error) error {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    "INTERNAL_SERVER_ERROR",
		Status:  http.StatusInternalServerError,
	}
}

func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		if message != "" {
			appErr.Message = message + ": " + appErr.Message
		}
		return appErr
	}

	// Otherwise wrap it as an internal error
	return NewInternalError(message, err)
}

func RespondWithError(c *gin.Context, err ErrorResponse, details ...any) {
	if len(details) > 0 {
		err.Details = details[0]
	}
	c.JSON(err.Status, err)
}

// Enhanced error handler that works with both existing errors and new AppError type
func HandleError(c *gin.Context, err error) {
	// First check if it's our new AppError type
	var appErr *AppError
	if errors.As(err, &appErr) {
		// Use the embedded details
		c.JSON(appErr.Status, ErrorResponse{
			Status:  appErr.Status,
			Code:    appErr.Code,
			Message: appErr.Message,
		})
		return
	}

	// Fall back to the existing error handling logic
	switch {
	case errors.Is(err, ErrInvalidInput), errors.Is(err, ErrValidationFailed):
		RespondWithError(c, ErrBadRequestResponse, err.Error())
	case errors.Is(err, ErrUnauthorized):
		RespondWithError(c, ErrUnauthorizedResponse)
	case errors.Is(err, ErrForbidden):
		RespondWithError(c, ErrorResponse{
			Status:  http.StatusForbidden,
			Code:    "FORBIDDEN",
			Message: "You don't have permission to access this resource",
		})
	case errors.Is(err, ErrNotFound):
		RespondWithError(c, ErrNotFoundResponse)
	case errors.Is(err, ErrDuplicateEntry):
		RespondWithError(c, ErrorResponse{
			Status:  http.StatusConflict,
			Code:    "CONFLICT",
			Message: "The resource already exists",
		})
	default:
		RespondWithError(c, ErrServerResponse)
	}
}

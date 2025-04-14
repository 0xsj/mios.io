package api

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrNotFound         = errors.New("resource not found")
	ErrInternalServer   = errors.New("internal server error")
	ErrDuplicateEntry   = errors.New("duplicate entry")
	ErrValidationFailed = errors.New("validation failed")
)

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Status  int    `json:"-"`          // HTTP status code, not shown in response
	Code    string `json:"code"`       // Application-specific error code
	Message string `json:"message"`    // User-friendly error message
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


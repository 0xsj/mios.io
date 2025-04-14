package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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

func RespondWithError(c *gin.Context, err ErrorResponse, details ...any) {
	if len(details) > 0 {
		err.Details = details[0]
	}
	c.JSON(err.Status, err)
}

func HandleError(c *gin.Context, err error) {
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
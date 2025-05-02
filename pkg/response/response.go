// pkg/response/response.go
package response

import (
	stderrors "errors" // Standard Go errors package with alias
	"net/http"

	"github.com/0xsj/gin-sqlc/pkg/errors" // Your custom errors package
	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Status  int    `json:"-"`                 // HTTP status code, not shown in response
	Code    string `json:"code"`              // Application-specific error code
	Message string `json:"message"`           // User-friendly error message
	Details any    `json:"details,omitempty"` // Optional details about the error
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	PerPage      int `json:"per_page"`
	TotalRecords int `json:"total_records"`
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

	ErrForbiddenResponse = ErrorResponse{
		Status:  http.StatusForbidden,
		Code:    "FORBIDDEN",
		Message: "You don't have permission to access this resource",
	}

	ErrNotFoundResponse = ErrorResponse{
		Status:  http.StatusNotFound,
		Code:    "NOT_FOUND",
		Message: "The requested resource was not found",
	}

	ErrConflictResponse = ErrorResponse{
		Status:  http.StatusConflict,
		Code:    "CONFLICT",
		Message: "The resource already exists",
	}

	ErrInternalServerResponse = ErrorResponse{
		Status:  http.StatusInternalServerError,
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "An unexpected error occurred",
	}

	ErrServiceUnavailableResponse = ErrorResponse{
		Status:  http.StatusServiceUnavailable,
		Code:    "SERVICE_UNAVAILABLE",
		Message: "The service is currently unavailable",
	}
)

// Success sends a successful response
func Success(c *gin.Context, data any, message string, statusCode ...int) {
	resp := Response{
		Success: true,
		Data:    data,
	}

	if message != "" {
		resp.Message = message
	}

	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	c.JSON(code, resp)
}

// WithPagination sends a response with pagination metadata
func WithPagination(c *gin.Context, data any, meta PaginationMeta) {
	resp := Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
	c.JSON(http.StatusOK, resp)
}

// NoContent sends a 204 No Content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error sends an error response
func Error(c *gin.Context, err ErrorResponse, details ...any) {
	if len(details) > 0 {
		err.Details = details[0]
	}
	c.JSON(err.Status, err)
}

// HandleError maps application errors to HTTP responses
func HandleError(c *gin.Context, err error) {
	// First check if it's our new AppError type
	var appErr *errors.AppError
	if stderrors.As(err, &appErr) {
		// Use the embedded details
		c.JSON(appErr.Status, ErrorResponse{
			Status:  appErr.Status,
			Code:    appErr.Code,
			Message: appErr.Message,
		})
		return
	}

	// Fall back to standard error types
	switch {
	case stderrors.Is(err, errors.ErrInvalidInput), stderrors.Is(err, errors.ErrValidationFailed):
		Error(c, ErrBadRequestResponse, err.Error())
	case stderrors.Is(err, errors.ErrUnauthorized):
		Error(c, ErrUnauthorizedResponse)
	case stderrors.Is(err, errors.ErrForbidden):
		Error(c, ErrForbiddenResponse)
	case stderrors.Is(err, errors.ErrNotFound):
		Error(c, ErrNotFoundResponse)
	case stderrors.Is(err, errors.ErrDuplicateEntry):
		Error(c, ErrConflictResponse)
	case stderrors.Is(err, errors.ErrDatabase) || stderrors.Is(err, errors.ErrExternalService):
		Error(c, ErrInternalServerResponse)
	default:
		Error(c, ErrInternalServerResponse)
	}
}
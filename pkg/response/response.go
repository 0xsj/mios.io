package response

import (
	stderrors "errors"
	"net/http"

	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/pkg/errors"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type PaginationMeta struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	PerPage      int `json:"per_page"`
	TotalRecords int `json:"total_records"`
}

var (
	ErrBadRequestResponse = ErrorResponse{
		Code:    "BAD_REQUEST",
		Message: "The request was invalid",
	}

	ErrUnauthorizedResponse = ErrorResponse{
		Code:    "UNAUTHORIZED",
		Message: "Authentication is required",
	}

	ErrForbiddenResponse = ErrorResponse{
		Code:    "FORBIDDEN",
		Message: "You don't have permission to access this resource",
	}

	ErrNotFoundResponse = ErrorResponse{
		Code:    "NOT_FOUND",
		Message: "The requested resource was not found",
	}

	ErrConflictResponse = ErrorResponse{
		Code:    "CONFLICT",
		Message: "The resource already exists",
	}

	ErrInternalServerResponse = ErrorResponse{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "An unexpected error occurred",
	}

	ErrServiceUnavailableResponse = ErrorResponse{
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

func WithPagination(c *gin.Context, data any, meta PaginationMeta) {
	resp := Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
	c.JSON(http.StatusOK, resp)
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Error(c *gin.Context, err ErrorResponse, details ...any) {
	if len(details) > 0 {
		err.Details = details[0]
	}

	statusCode := http.StatusInternalServerError

	switch err.Code {
	case "BAD_REQUEST", "VALIDATION_ERROR":
		statusCode = http.StatusBadRequest
	case "UNAUTHORIZED":
		statusCode = http.StatusUnauthorized
	case "FORBIDDEN":
		statusCode = http.StatusForbidden
	case "NOT_FOUND":
		statusCode = http.StatusNotFound
	case "CONFLICT":
		statusCode = http.StatusConflict
	case "SERVICE_UNAVAILABLE":
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, err)
}

func HandleError(c *gin.Context, err error, logger log.Logger) {
	var appErr *errors.AppError
	if stderrors.As(err, &appErr) {
		appErr.Log(logger)

		c.JSON(appErr.Status, ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
		})
		return
	}

	switch {
	case stderrors.Is(err, errors.ErrInvalidInput), stderrors.Is(err, errors.ErrValidationFailed):
		logger.Warn("Bad request error:", err)
		Error(c, ErrBadRequestResponse, err.Error())
	case stderrors.Is(err, errors.ErrUnauthorized):
		logger.Warn("Unauthorized error:", err)
		Error(c, ErrUnauthorizedResponse)
	case stderrors.Is(err, errors.ErrForbidden):
		logger.Warn("Forbidden error:", err)
		Error(c, ErrForbiddenResponse)
	case stderrors.Is(err, errors.ErrNotFound):
		logger.Info("Not found error:", err)
		Error(c, ErrNotFoundResponse)
	case stderrors.Is(err, errors.ErrDuplicateEntry):
		logger.Warn("Conflict error:", err)
		Error(c, ErrConflictResponse)
	case stderrors.Is(err, errors.ErrDatabase) || stderrors.Is(err, errors.ErrExternalService):
		logger.Error("Database/external service error:", err)
		Error(c, ErrInternalServerResponse)
	default:
		logger.Error("Unhandled error:", err)
		Error(c, ErrInternalServerResponse)
	}
}

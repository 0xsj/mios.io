// pkg/errors/errors.go
package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/0xsj/gin-sqlc/log"
)

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

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

type AppError struct {
	Err      error
	Message  string
	Code     string
	Status   int
	LogLevel LogLevel
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Is(target error) bool {
	if e.Err == nil {
		return false
	}
	return errors.Is(e.Err, target)
}

func (e *AppError) Log(logger log.Logger) {
	errMsg := fmt.Sprintf("Error: %s (Code: %s, Status: %d)",
		e.Message, e.Code, e.Status)

	if e.Err != nil {
		errMsg = fmt.Sprintf("%s, Cause: %v", errMsg, e.Err)
	}

	switch e.LogLevel {
	case LogLevelDebug:
		logger.Debug(errMsg)
	case LogLevelInfo:
		logger.Info(errMsg)
	case LogLevelWarn:
		logger.Warn(errMsg)
	case LogLevelError:
		logger.Error(errMsg)
	case LogLevelFatal:
		logger.Fatal(errMsg)
	default:
		logger.Error(errMsg)
	}
}

func NewBadRequestError(message string, err error) *AppError {
	return &AppError{
		Err:      err,
		Message:  message,
		Code:     "BAD_REQUEST",
		Status:   http.StatusBadRequest,
		LogLevel: LogLevelWarn,
	}
}

func NewUnauthorizedError(message string, err error) *AppError {
	return &AppError{
		Err:      err,
		Message:  message,
		Code:     "UNAUTHORIZED",
		Status:   http.StatusUnauthorized,
		LogLevel: LogLevelWarn,
	}
}

func NewForbiddenError(message string, err error) *AppError {
	return &AppError{
		Err:      err,
		Message:  message,
		Code:     "FORBIDDEN",
		Status:   http.StatusForbidden,
		LogLevel: LogLevelWarn,
	}
}

func NewNotFoundError(message string, err error) *AppError {
	return &AppError{
		Err:      err,
		Message:  message,
		Code:     "NOT_FOUND",
		Status:   http.StatusNotFound,
		LogLevel: LogLevelInfo,
	}
}

func NewConflictError(message string, err error) *AppError {
	return &AppError{
		Err:      err,
		Message:  message,
		Code:     "CONFLICT",
		Status:   http.StatusConflict,
		LogLevel: LogLevelWarn,
	}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Err:      err,
		Message:  message,
		Code:     "INTERNAL_SERVER_ERROR",
		Status:   http.StatusInternalServerError,
		LogLevel: LogLevelError,
	}
}

func NewValidationError(message string, err error) *AppError {
	return &AppError{
		Err:      err,
		Message:  message,
		Code:     "VALIDATION_ERROR",
		Status:   http.StatusBadRequest,
		LogLevel: LogLevelInfo,
	}
}

func NewDatabaseError(message string, err error) *AppError {
	return &AppError{
		Err:      err,
		Message:  message,
		Code:     "DATABASE_ERROR",
		Status:   http.StatusInternalServerError,
		LogLevel: LogLevelError,
	}
}

func NewExternalServiceError(message string, err error) *AppError {
	return &AppError{
		Err:      err,
		Message:  message,
		Code:     "EXTERNAL_SERVICE_ERROR",
		Status:   http.StatusInternalServerError,
		LogLevel: LogLevelError,
	}
}

func Wrap(err error, message string) error {
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

	return &AppError{
		Err:      err,
		Message:  message,
		Code:     "INTERNAL_SERVER_ERROR",
		Status:   http.StatusInternalServerError,
		LogLevel: LogLevelError,
	}
}

func WrapWith(err error, message string, errType *AppError) error {
	if err == nil {
		return nil
	}

	return &AppError{
		Err:      err,
		Message:  message,
		Code:     errType.Code,
		Status:   errType.Status,
		LogLevel: errType.LogLevel,
	}
}

const (
	PgErrUniqueViolation     = "23505"
	PgErrForeignKeyViolation = "23503"
	PgErrCheckViolation      = "23514"
)

func IsPgError(err error, code string) bool {
	pgErr, ok := err.(interface {
		Code() string
	})

	if ok && pgErr.Code() == code {
		return true
	}
	return false
}

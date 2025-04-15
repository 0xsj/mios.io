package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRespondWithError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		errorResponse  ErrorResponse
		details        any
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "basic error",
			errorResponse: ErrorResponse{
				Status:  http.StatusBadRequest,
				Code:    "BAD_REQUEST",
				Message: "Invalid input",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":"BAD_REQUEST","message":"Invalid input"}`,
		},
		{
			name: "error with details",
			errorResponse: ErrorResponse{
				Status:  http.StatusNotFound,
				Code:    "NOT_FOUND",
				Message: "Resource not found",
			},
			details:        "User with ID 123 not found",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"code":"NOT_FOUND","message":"Resource not found","details":"User with ID 123 not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.details != nil {
				RespondWithError(c, tt.errorResponse, tt.details)
			} else {
				RespondWithError(c, tt.errorResponse)
			}

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestHandleError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "invalid input error",
			err:            ErrInvalidInput,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "BAD_REQUEST",
		},
		{
			name:           "unauthorized error",
			err:            ErrUnauthorized,
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "UNAUTHORIZED",
		},
		{
			name:           "forbidden error",
			err:            ErrForbidden,
			expectedStatus: http.StatusForbidden,
			expectedCode:   "FORBIDDEN",
		},
		{
			name:           "not found error",
			err:            ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedCode:   "NOT_FOUND",
		},
		{
			name:           "duplicate entry error",
			err:            ErrDuplicateEntry,
			expectedStatus: http.StatusConflict,
			expectedCode:   "CONFLICT",
		},
		{
			name:           "unknown error",
			err:            errors.New("some unknown error"),
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_SERVER_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			HandleError(c, tt.err)

			assert.Equal(t, tt.expectedStatus, w.Code)
			
			assert.Contains(t, w.Body.String(), tt.expectedCode)
		})
	}
}
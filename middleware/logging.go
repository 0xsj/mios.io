// middleware/logging.go
package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/0xsj/gin-sqlc/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
	statusCode int
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggingMiddleware(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate a request ID
		requestID := uuid.New().String()[:8]

		// Add request ID to context
		c.Set("request_id", requestID)

		// Create request-specific logger
		reqLogger := logger.With("request_id", requestID)

		// Log the start of the request
		reqLogger.Infof("Request started: %s %s", c.Request.Method, c.Request.URL.Path)

		// Time the request
		start := time.Now()

		// Store the original body if needed for debugging
		var requestBody []byte
		if c.Request.Body != nil && shouldLogRequestBody(c.Request.URL.Path) {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create a response writer that captures the response
		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w

		// Process the request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Determine log level based on status
		status := c.Writer.Status()

		// Get a human-readable status description
		statusText := http.StatusText(status)

		// Create a structured log entry with duration
		reqLoggerWithDuration := reqLogger.With("duration_ms", duration.Milliseconds())

		// Log with a more informative format
		if status >= 500 {
			reqLoggerWithDuration.Errorf("Request completed: %s %s, status: %d %s",
				c.Request.Method,
				c.Request.URL.Path,
				status,
				statusText)
		} else if status >= 400 {
			reqLoggerWithDuration.Warnf("Request completed: %s %s, status: %d %s",
				c.Request.Method,
				c.Request.URL.Path,
				status,
				statusText)
		} else {
			reqLoggerWithDuration.Infof("Request completed: %s %s, status: %d %s",
				c.Request.Method,
				c.Request.URL.Path,
				status,
				statusText)
		}

		// Only log error responses at debug level to avoid cluttering logs
		if status >= 400 && len(w.body.Bytes()) > 0 && len(w.body.Bytes()) < 1000 {
			reqLogger.Debugf("Response body: %s", w.body.String())
		}
	}
}

// Only log bodies for certain paths to avoid logging sensitive data
func shouldLogRequestBody(path string) bool {
	// Don't log auth endpoints to avoid capturing passwords
	if strings.Contains(path, "/auth/") {
		return false
	}
	return true
}

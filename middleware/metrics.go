// middleware/metrics.go
package middleware

import (
	"strconv"
	"time"

	"github.com/0xsj/gin-sqlc/pkg/metrics"
	"github.com/gin-gonic/gin"
)

// MetricsMiddleware creates a middleware that collects HTTP metrics
func MetricsMiddleware(m *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Increment active connections
		m.HTTPActiveConnections.Inc()
		defer m.HTTPActiveConnections.Dec()

		// Get request size
		requestSize := int64(0)
		if c.Request.ContentLength > 0 {
			requestSize = c.Request.ContentLength
		}

		// Create a response writer that captures the response size
		writer := &metricsResponseWriter{
			ResponseWriter: c.Writer,
			statusCode:     200,
			bytesWritten:   0,
		}
		c.Writer = writer

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get normalized endpoint for metrics
		endpoint := metrics.NormalizeEndpoint(c.FullPath())
		if endpoint == "" {
			endpoint = "unknown"
		}

		// Record metrics
		m.RecordHTTPRequest(
			c.Request.Method,
			endpoint,
			writer.statusCode,
			duration,
			requestSize,
			int64(writer.bytesWritten),
		)
	}
}

// metricsResponseWriter wraps gin.ResponseWriter to capture response size and status
type metricsResponseWriter struct {
	gin.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (w *metricsResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *metricsResponseWriter) Write(data []byte) (int, error) {
	size, err := w.ResponseWriter.Write(data)
	w.bytesWritten += size
	return size, err
}

// RecoveryWithMetrics creates a recovery middleware that also records panic metrics
func RecoveryWithMetrics(m *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				m.PanicTotal.Inc()
				m.RecordError("panic", "middleware", "critical")
				
				// Re-panic to let the original recovery middleware handle it
				panic(r)
			}
		}()
		c.Next()
	}
}

// RateLimitMetricsWrapper wraps rate limiting middleware to add metrics
func RateLimitMetricsWrapper(m *metrics.Metrics, next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		endpoint := metrics.NormalizeEndpoint(c.FullPath())
		
		// Call the rate limit middleware
		next(c)
		
		// Check if request was blocked by rate limiting
		if c.IsAborted() && c.Writer.Status() == 429 {
			m.RecordRateLimit(endpoint, false, 0, "unknown")
		} else {
			// Extract remaining count from headers if available
			remaining := 0
			if remainingHeader := c.GetHeader("X-RateLimit-Remaining"); remainingHeader != "" {
				if r, err := strconv.Atoi(remainingHeader); err == nil {
					remaining = r
				}
			}
			m.RecordRateLimit(endpoint, true, remaining, "unknown")
		}
	}
}
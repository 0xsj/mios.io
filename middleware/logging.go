// middleware/logging.go
package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/0xsj/gin-sqlc/log"
	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggingMiddleware(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}
		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w
		
		c.Next()
		
		duration := time.Since(start)
		status := c.Writer.Status()
		
		if status >= 500 {
			logger.Errorf("Request: %s %s, Status: %d, Duration: %s, Error: %s",
				c.Request.Method, c.Request.URL.Path, status, duration, w.body.String())
		} else if status >= 400 {
			logger.Warnf("Request: %s %s, Status: %d, Duration: %s",
				c.Request.Method, c.Request.URL.Path, status, duration)
		} else {
			logger.Infof("Request: %s %s, Status: %d, Duration: %s",
				c.Request.Method, c.Request.URL.Path, status, duration)
		}
	}
}
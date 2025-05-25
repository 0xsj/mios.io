package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimitConfig struct {
	RequestsPerMinute int           // Number of requests allowed per minute
	BurstSize         int           // Number of requests allowed in burst
	KeyGenerator      KeyGenerator  // Function to generate rate limit key
	SkipSuccessful    bool          // Skip counting successful requests (2xx responses)
	WindowSize        time.Duration // Time window for rate limiting
}

type KeyGenerator func (c *gin.Context) string
// middleware/rate_limit.go
package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/pkg/redis"
	"github.com/0xsj/gin-sqlc/pkg/response"
	"github.com/gin-gonic/gin"
)

type RateLimitConfig struct {
	RequestsPerMinute int           // Number of requests allowed per minute
	BurstSize         int           // Number of requests allowed in burst
	KeyGenerator      KeyGenerator  // Function to generate rate limit key
	SkipSuccessful    bool          // Skip counting successful requests (2xx responses)
	WindowSize        time.Duration // Time window for rate limiting
}

type KeyGenerator func(c *gin.Context) string

// Default key generators
func IPBasedKeyGenerator(c *gin.Context) string {
	return fmt.Sprintf("rate_limit:ip:%s", c.ClientIP())
}

func UserBasedKeyGenerator(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists {
		return IPBasedKeyGenerator(c) // Fallback to IP
	}
	return fmt.Sprintf("rate_limit:user:%v", userID)
}

func EndpointBasedKeyGenerator(c *gin.Context) string {
	return fmt.Sprintf("rate_limit:endpoint:%s:%s:%s", c.Request.Method, c.FullPath(), c.ClientIP())
}

// Pre-configured rate limit configs
func DefaultRateLimit() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerMinute: 60,
		BurstSize:         10,
		KeyGenerator:      IPBasedKeyGenerator,
		SkipSuccessful:    false,
		WindowSize:        time.Minute,
	}
}

func StrictRateLimit() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerMinute: 30,
		BurstSize:         5,
		KeyGenerator:      IPBasedKeyGenerator,
		SkipSuccessful:    false,
		WindowSize:        time.Minute,
	}
}

func AuthenticatedUserRateLimit() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerMinute: 120,
		BurstSize:         20,
		KeyGenerator:      UserBasedKeyGenerator,
		SkipSuccessful:    true, // Don't count successful requests for authenticated users
		WindowSize:        time.Minute,
	}
}

func ExpensiveOperationRateLimit() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerMinute: 10,
		BurstSize:         2,
		KeyGenerator:      UserBasedKeyGenerator,
		SkipSuccessful:    false,
		WindowSize:        time.Minute,
	}
}

type RateLimiter struct {
	redisClient *redis.Client
	logger      log.Logger
	config      RateLimitConfig
}

func NewRateLimiter(redisClient *redis.Client, logger log.Logger, config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		logger:      logger,
		config:      config,
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := rl.config.KeyGenerator(c)
		
		// Check rate limit before processing request
		allowed, remaining, resetTime, err := rl.checkRateLimit(c, key)
		if err != nil {
			rl.logger.Errorf("Rate limit check failed: %v", err)
			// Allow request on Redis errors (fail open)
			c.Next()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(rl.config.RequestsPerMinute))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		if !allowed {
			rl.logger.Warnf("Rate limit exceeded for key: %s", key)
			response.Error(c, response.ErrorResponse{
				Code:    "RATE_LIMIT_EXCEEDED",
				Message: "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		// Store original writer to check response status later
		writer := &responseWriter{ResponseWriter: c.Writer, statusCode: 200}
		c.Writer = writer

		// Process request
		c.Next()

		// Update rate limit counter based on response
		shouldCount := true
		if rl.config.SkipSuccessful && writer.statusCode >= 200 && writer.statusCode < 300 {
			shouldCount = false
		}

		if shouldCount {
			if err := rl.incrementCounter(c, key); err != nil {
				rl.logger.Errorf("Failed to increment rate limit counter: %v", err)
			}
		}
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rl *RateLimiter) checkRateLimit(ctx context.Context, key string) (allowed bool, remaining int, resetTime time.Time, err error) {
	now := time.Now()
	windowStart := now.Truncate(rl.config.WindowSize)
	resetTime = windowStart.Add(rl.config.WindowSize)

	// Use sliding window with Redis
	currentCountStr, err := rl.redisClient.Get(ctx, key)
	if err != nil && err != redis.Nil {
		return false, 0, resetTime, err
	}

	currentCount := 0
	if currentCountStr != "" {
		currentCount, _ = strconv.Atoi(currentCountStr)
	}

	// Check burst limit
	burstKey := key + ":burst"
	burstCountStr, err := rl.redisClient.Get(ctx, burstKey)
	if err != nil && err != redis.Nil {
		return false, 0, resetTime, err
	}

	burstCount := 0
	if burstCountStr != "" {
		burstCount, _ = strconv.Atoi(burstCountStr)
	}

	// Check if within limits
	if currentCount >= rl.config.RequestsPerMinute || burstCount >= rl.config.BurstSize {
		remaining = max(0, rl.config.RequestsPerMinute-currentCount)
		return false, remaining, resetTime, nil
	}

	remaining = rl.config.RequestsPerMinute - currentCount - 1
	return true, remaining, resetTime, nil
}

func (rl *RateLimiter) incrementCounter(ctx context.Context, key string) error {
	now := time.Now()
	fmt.Println(now)
	
	// Increment main counter
	count, err := rl.redisClient.Incr(ctx, key)
	if err != nil {
		return err
	}

	// Set expiration on first increment
	if count == 1 {
		err = rl.redisClient.Expire(ctx, key, rl.config.WindowSize)
		if err != nil {
			return err
		}
	}

	// Increment burst counter (shorter window for burst protection)
	burstKey := key + ":burst"
	burstCount, err := rl.redisClient.Incr(ctx, burstKey)
	if err != nil {
		return err
	}

	// Set shorter expiration for burst counter (e.g., 10 seconds)
	if burstCount == 1 {
		err = rl.redisClient.Expire(ctx, burstKey, 10*time.Second)
		if err != nil {
			return err
		}
	}

	return nil
}

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Rate limit middleware factory functions
func RateLimitMiddleware(redisClient *redis.Client, logger log.Logger) gin.HandlerFunc {
	limiter := NewRateLimiter(redisClient, logger, DefaultRateLimit())
	return limiter.Middleware()
}

func StrictRateLimitMiddleware(redisClient *redis.Client, logger log.Logger) gin.HandlerFunc {
	limiter := NewRateLimiter(redisClient, logger, StrictRateLimit())
	return limiter.Middleware()
}

func AuthUserRateLimitMiddleware(redisClient *redis.Client, logger log.Logger) gin.HandlerFunc {
	limiter := NewRateLimiter(redisClient, logger, AuthenticatedUserRateLimit())
	return limiter.Middleware()
}

func ExpensiveOpRateLimitMiddleware(redisClient *redis.Client, logger log.Logger) gin.HandlerFunc {
	limiter := NewRateLimiter(redisClient, logger, ExpensiveOperationRateLimit())
	return limiter.Middleware()
}
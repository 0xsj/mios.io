package middleware

import (
	"strings"

	"github.com/0xsj/mios.io/log"
	"github.com/0xsj/mios.io/pkg/context"
	"github.com/0xsj/mios.io/pkg/errors"
	"github.com/0xsj/mios.io/pkg/response"
	"github.com/0xsj/mios.io/pkg/token"
	"github.com/0xsj/mios.io/service"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens from the Authorization header
func AuthMiddleware(authService service.AuthService, logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Authentication failed: missing Authorization header")
			response.Error(c, response.ErrUnauthorizedResponse)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("Authentication failed: invalid Authorization header format")
			response.Error(c, response.ErrUnauthorizedResponse)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := authService.ValidateToken(c, tokenString)
		if err != nil {
			logger.Warn("Authentication failed: invalid token:", err)
			response.HandleError(c, err, logger)
			c.Abort()
			return
		}

		logger.Debugf("User authenticated: %s", claims.UserID)
		context.SetUserID(c, claims.UserID)
		c.Set("claims", claims)
		c.Next()
	}
}

// AdminMiddleware ensures that the authenticated user has admin privileges
func AdminMiddleware(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsValue, exists := c.Get("claims")
		if !exists {
			logger.Warn("Admin access denied: no authentication claims found")
			response.Error(c, response.ErrForbiddenResponse)
			c.Abort()
			return
		}

		claims, ok := claimsValue.(*token.Claims)
		if !ok {
			logger.Warn("Admin access denied: invalid claims type")
			response.Error(c, response.ErrorResponse{
				Code:    "FORBIDDEN",
				Message: "Invalid authentication data",
			})
			c.Abort()
			return
		}

		if !claims.IsAdmin {
			logger.Warnf("Admin access denied for user: %s", claims.UserID)
			response.Error(c, response.ErrForbiddenResponse)
			c.Abort()
			return
		}

		logger.Debugf("Admin access granted for user: %s", claims.UserID)
		c.Next()
	}
}

// RequireVerifiedEmail ensures the user's email is verified before allowing access
func RequireVerifiedEmail(authService service.AuthService, logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := context.GetUserID(c)
		if err != nil {
			logger.Warn("Email verification check failed: user ID not found in context")
			response.Error(c, response.ErrUnauthorizedResponse)
			c.Abort()
			return
		}

		verified, err := authService.IsEmailVerified(c, userID)
		if err != nil {
			logger.Warnf("Email verification check failed for user %s: %v", userID, err)
			response.HandleError(c, err, logger)
			c.Abort()
			return
		}

		if !verified {
			logger.Warnf("Access denied for user %s: email not verified", userID)
			response.Error(c, response.ErrorResponse{
				Code:    "EMAIL_NOT_VERIFIED",
				Message: "Email verification is required for this action",
			})
			c.Abort()
			return
		}

		logger.Debugf("Email verified for user: %s", userID)
		c.Next()
	}
}

// RequestLogger logs information about incoming requests
func RequestLogger(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log before request
		logger.Infof("Request started: %s %s", c.Request.Method, c.Request.URL.Path)

		// Process request
		c.Next()

		// Log after request
		logger.Infof("Request completed: %s %s, status: %d",
			c.Request.Method, c.Request.URL.Path, c.Writer.Status())
	}
}

// Recovery middleware to handle panics
func Recovery(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("Panic recovered: %v", r)
				err := errors.NewInternalError("Internal server error", nil)
				response.HandleError(c, err, logger)
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

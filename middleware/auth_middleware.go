package middleware

import (
	"strings"

	"github.com/0xsj/gin-sqlc/api"
	"github.com/0xsj/gin-sqlc/pkg/context"
	"github.com/0xsj/gin-sqlc/pkg/token"
	"github.com/0xsj/gin-sqlc/service"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			api.RespondWithError(c, api.ErrUnauthorizedResponse)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			api.RespondWithError(c, api.ErrUnauthorizedResponse)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := authService.ValidateToken(c, tokenString)
		if err != nil {
			api.RespondWithError(c, api.ErrUnauthorizedResponse)
			c.Abort()
			return
		}

		context.SetUserID(c, claims.UserID)
		c.Set("claims", claims)

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsValue, exists := c.Get("claims")
		if !exists {
			api.RespondWithError(c, api.ErrorResponse{
				Status:  403,
				Code:    "FORBIDDEN",
				Message: "Admin access required",
			})
			c.Abort()
			return
		}

		claims, ok := claimsValue.(*token.Claims)
		if !ok {
			api.RespondWithError(c, api.ErrorResponse{
				Status:  403,
				Code:    "FORBIDDEN",
				Message: "Invalid authentication data",
			})
			c.Abort()
			return
		}

		if !claims.IsAdmin {
			api.RespondWithError(c, api.ErrorResponse{
				Status:  403,
				Code:    "FORBIDDEN",
				Message: "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireVerifiedEmail(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := context.GetUserID(c)
		if err != nil {
			api.RespondWithError(c, api.ErrUnauthorizedResponse)
			c.Abort()
			return
		}

		verified, err := authService.IsEmailVerified(c, userID)
		if err != nil || !verified {
			api.RespondWithError(c, api.ErrorResponse{
				Status:  403,
				Code:    "EMAIL_NOT_VERIFIED",
				Message: "email verification is required for this action",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

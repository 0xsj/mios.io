package auth

import (
	"net/http"

	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/pkg/response"
	"github.com/0xsj/gin-sqlc/service"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for authentication operations
type Handler struct {
	authService service.AuthService
	logger      log.Logger
}

// NewHandler creates a new auth handler
func NewHandler(authService service.AuthService, logger log.Logger) *Handler {
	return &Handler{
		authService: authService,
		logger:      logger,
	}
}

// RegisterRoutes registers auth routes on the given router
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
		authGroup.POST("/refresh", h.RefreshToken)
		authGroup.POST("/forgot-password", h.ForgotPassword)
		authGroup.POST("/reset-password", h.ResetPassword)
		authGroup.POST("/verify-email", h.VerifyEmail)
		authGroup.POST("/logout", h.Logout)
	}

	h.logger.Info("Auth routes registered successfully")
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	h.logger.Info("Register handler called")

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	h.logger.Debugf("Received registration request for email: %s", req.Email)

	input := service.RegisterInput{
		Username:        req.Username,
		Handle:          req.Handle,
		Email:           req.Email,
		Password:        req.Password,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Bio:             req.Bio,
		ProfileImageURL: req.ProfileImageURL,
		LayoutVersion:   req.LayoutVersion,
		CustomDomain:    req.CustomDomain,
	}

	user, err := h.authService.Register(c, input)
	if err != nil {
		h.logger.Errorf("Failed to register user: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("User registered successfully with ID: %s", user.ID)
	response.Success(c, user, "User registered successfully", http.StatusCreated)
}

// Login handles user authentication
func (h *Handler) Login(c *gin.Context) {
	h.logger.Info("Login handler called")

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	h.logger.Debugf("Received login request for email: %s", req.Email)

	input := service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	tokenResponse, err := h.authService.Login(c, input)
	if err != nil {
		h.logger.Warnf("Login failed for email %s: %v", req.Email, err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("Login successful for user: %s", tokenResponse.User.ID)
	response.Success(c, tokenResponse, "Login successful")
}

// RefreshToken refreshes an authentication token
func (h *Handler) RefreshToken(c *gin.Context) {
	h.logger.Info("RefreshToken handler called")

	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	// Don't log token details for security reasons
	h.logger.Debug("Received token refresh request")

	input := service.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}

	tokenResponse, err := h.authService.RefreshToken(c, input)
	if err != nil {
		h.logger.Warnf("Token refresh failed: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("Token refreshed successfully for user: %s", tokenResponse.User.ID)
	response.Success(c, tokenResponse, "Token refreshed successfully")
}

// ForgotPassword initiates the password reset process
func (h *Handler) ForgotPassword(c *gin.Context) {
	h.logger.Info("ForgotPassword handler called")

	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	h.logger.Debugf("Received forgot password request for email: %s", req.Email)

	err := h.authService.GenerateResetToken(c, req.Email)
	if err != nil {
		// Log but don't expose failure details to client for security
		h.logger.Warnf("Failed to generate reset token: %v", err)

		// For security, always return the same message regardless of whether the email exists
		response.Success(c, nil, "If the email exists, a reset link has been sent")
		return
	}

	h.logger.Infof("Password reset token generated for email: %s", req.Email)
	response.Success(c, nil, "If the email exists, a reset link has been sent")
}

// ResetPassword completes the password reset process
func (h *Handler) ResetPassword(c *gin.Context) {
	h.logger.Info("ResetPassword handler called")

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	h.logger.Debugf("Received password reset request for email: %s", req.Email)

	input := service.ResetPasswordInput{
		Token:           req.Token,
		Email:           req.Email,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	}

	err := h.authService.ResetPassword(c, input)
	if err != nil {
		h.logger.Errorf("Failed to reset password: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("Password reset successfully for email: %s", req.Email)
	response.Success(c, nil, "Password has been reset successfully")
}

// VerifyEmail validates a user's email address
func (h *Handler) VerifyEmail(c *gin.Context) {
	h.logger.Info("VerifyEmail handler called")

	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	// For security reasons, don't log the actual token
	h.logger.Debug("Received email verification request")

	err := h.authService.VerifyEmail(c, req.Token)
	if err != nil {
		h.logger.Errorf("Failed to verify email: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Info("Email verified successfully")
	response.Success(c, nil, "Email has been verified successfully")
}

// Logout ends a user's session
func (h *Handler) Logout(c *gin.Context) {
	h.logger.Info("Logout handler called")

	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	h.logger.Debugf("Received logout request for user ID: %s", req.UserID)

	err := h.authService.Logout(c, req.UserID)
	if err != nil {
		h.logger.Errorf("Failed to logout user: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("User logged out successfully: %s", req.UserID)
	response.Success(c, nil, "Logged out successfully")
}

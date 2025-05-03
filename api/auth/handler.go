package auth

import (
	"fmt"
	"net/http"

	"github.com/0xsj/gin-sqlc/api"
	"github.com/0xsj/gin-sqlc/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	authService service.AuthService
}

func NewHandler(authservice service.AuthService) *Handler {
	return &Handler{
		authService: authservice,
	}
}

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
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

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
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, user, "User registered successfully", http.StatusCreated)
}

// func (h *Handler) Login(c *gin.Context) {
// 	var req LoginRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		api.HandleError(c, api.ErrInvalidInput)
// 		return
// 	}

// 	input := service.LoginInput{
// 		Email:    req.Email,
// 		Password: req.Password,
// 	}

// 	response, err := h.authService.Login(c, input)
// 	if err != nil {
// 		api.HandleError(c, err)
// 		return
// 	}

// 	api.RespondWithSuccess(c, response, "Login successful")
// }

func (h *Handler) Login(c *gin.Context) {
	fmt.Println("Login handler called")
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Login binding error: %v\n", err)
		api.HandleError(c, api.ErrInvalidInput)
		return
	}
	fmt.Printf("Login request: %s\n", req.Email)

	input := service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	response, err := h.authService.Login(c, input)
	if err != nil {
		fmt.Printf("Login service error: %v\n", err)
		api.HandleError(c, err)
		return
	}

	fmt.Printf("Login successful for user: %s\n", response.User.ID)
	api.RespondWithSuccess(c, response, "Login successful")
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	input := service.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}

	response, err := h.authService.RefreshToken(c, input)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, response, "Token refreshed successfully")
}

func (h *Handler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	err := h.authService.GenerateResetToken(c, req.Email)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, nil, "If the email exists, a reset link has been sent")
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var req RestPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	input := service.ResetPasswordInput{
		Token:           req.Token,
		Email:           req.Email,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	}

	err := h.authService.ResetPassword(c, input)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, nil, "Password has been reset successfully")
}

func (h *Handler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	err := h.authService.VerifyEmail(c, req.Token)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, nil, "Email has been verified successfully")
}

func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	err := h.authService.Logout(c, req.UserID)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, nil, "Logged out successfully")
}

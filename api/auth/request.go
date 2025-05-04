package auth

// Request types

// RegisterRequest represents the payload for user registration
type RegisterRequest struct {
	Username        string `json:"username" binding:"required"`
	Handle          string `json:"handle"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Bio             string `json:"bio"`
	ProfileImageURL string `json:"profile_image_url"`
	LayoutVersion   string `json:"layout_version"`
	CustomDomain    string `json:"custom_domain"`
}

// LoginRequest represents the payload for user authentication
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest represents the payload for refreshing an authentication token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ForgotPasswordRequest represents the payload for requesting a password reset
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents the payload for resetting a password
type ResetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

// VerifyEmailRequest represents the payload for email verification
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// LogoutRequest represents the payload for ending a user session
type LogoutRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// Response types

// UserResponse represents a user in response payloads
type UserResponse struct {
	ID              string `json:"id"`
	Username        string `json:"username"`
	Handle          string `json:"handle"`
	Email           string `json:"email"`
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	Bio             string `json:"bio,omitempty"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
	LayoutVersion   string `json:"layout_version,omitempty"`
	CustomDomain    string `json:"custom_domain,omitempty"`
	IsPremium       bool   `json:"is_premium"`
	IsAdmin         bool   `json:"is_admin"`
	Onboarded       bool   `json:"onboarded"`
	IsEmailVerified bool   `json:"is_email_verified"`
}

// TokenResponse represents authentication tokens in response payloads
type TokenResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    int64        `json:"expires_at"`
	User         *UserResponse `json:"user"`
}
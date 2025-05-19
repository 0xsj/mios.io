package service

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/pkg/email"
	"github.com/0xsj/gin-sqlc/pkg/errors"
	"github.com/0xsj/gin-sqlc/pkg/password"
	"github.com/0xsj/gin-sqlc/pkg/token"
	"github.com/0xsj/gin-sqlc/repository"
	"github.com/google/uuid"
)

const (
	DefaultAccessTokenDuration  = 24 * time.Hour
	DefaultRefreshTokenDuration = 7 * 24 * time.Hour
	DefaultResetTokenDuration   = 1 * time.Hour
)

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*UserDTO, error)
	Login(ctx context.Context, input LoginInput) (*TokenResponse, error)
	RefreshToken(ctx context.Context, input RefreshTokenRequest) (*TokenResponse, error)
	GenerateResetToken(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, input ResetPasswordInput) error
	// VerifyEmail(ctx context.Context, token string) error
	Logout(ctx context.Context, userID string) error
	ValidateToken(ctx context.Context, tokenStr string) (*token.Claims, error)
	IsEmailVerified(ctx context.Context, userID string) (bool, error)
	SendVerificationEmail(ctx context.Context, email, username, token string) error
	SendPasswordResetEmail(ctx context.Context, email, username, token string) error
	SendPasswordChangedEmail(ctx context.Context, email, username string) error
	SendAccountLockedEmail(ctx context.Context, email, username, unlockTime string) error
}

type RegisterInput struct {
	Username        string `json:"username" binding:"required"`
	Handle          string `json:"handle" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Bio             string `json:"bio"`
	ProfileImageURL string `json:"profile_image_url"`
	LayoutVersion   string `json:"layout_version"`
	CustomDomain    string `json:"custom_domain"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ResetPasswordInput struct {
	Token           string `json:"token" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

type TokenResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresAt    int64    `json:"expires_at"`
	User         *UserDTO `json:"user"`
}

type authService struct {
	userRepo    repository.UserRepository
	authRepo    repository.AuthRepository
	emailClient *email.EmailClient
	jwtSecret   string
	tokenExpiry time.Duration
	logger      log.Logger
	baseURL     string
}

func NewAuthService(
	userRepo repository.UserRepository,
	authRepo repository.AuthRepository,
	emailClient *email.EmailClient,
	jwtSecret string,
	tokenExpiry time.Duration,
	logger log.Logger,
	baseURL string,
) AuthService {
	if tokenExpiry == 0 {
		tokenExpiry = DefaultAccessTokenDuration
	}

	return &authService{
		userRepo:    userRepo,
		authRepo:    authRepo,
		emailClient: emailClient,
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
		logger:      logger,
		baseURL:     baseURL,
	}
}

func (s *authService) Register(ctx context.Context, input RegisterInput) (*UserDTO, error) {
	s.logger.Infof("Registering new user with email: %s and username: %s", input.Email, input.Username)

	err := password.ValidatePassword(input.Password, password.DefaultPasswordConfig())
	if err != nil {
		s.logger.Warnf("Password validation failed for new user registration: %v", err)
		return nil, errors.NewValidationError("Invalid password format", err)
	}

	_, err = s.userRepo.GetUserByEmail(ctx, input.Email)
	if err == nil {
		s.logger.Warnf("Registration failed: email %s already exists", input.Email)
		return nil, errors.NewConflictError("Email already registered", nil)
	} else if !errors.IsNotFound(err) {
		s.logger.Errorf("Error checking existing email: %v", err)
		return nil, errors.Wrap(err, "Failed to check existing email")
	}

	_, err = s.userRepo.GetUserByUsername(ctx, input.Username)
	if err == nil {
		s.logger.Warnf("Registration failed: username %s already exists", input.Username)
		return nil, errors.NewConflictError("Username already taken", nil)
	} else if !errors.IsNotFound(err) {
		s.logger.Errorf("Error checking existing username: %v", err)
		return nil, errors.Wrap(err, "Failed to check existing username")
	}

	userParams := repository.CreateUserParams{
		Username:        input.Username,
		Handle:          input.Handle,
		Email:           input.Email,
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		Bio:             input.Bio,
		ProfileImageURL: input.ProfileImageURL,
		LayoutVersion:   input.LayoutVersion,
		CustomDomain:    input.CustomDomain,
		IsPremium:       false,
		IsAdmin:         false,
		Onboarded:       false,
	}

	user, err := s.userRepo.CreateUser(ctx, userParams)
	if err != nil {
		s.logger.Errorf("Failed to create user: %v", err)
		return nil, errors.Wrap(err, "Failed to create user")
	}

	hashedPassword, salt, err := password.HashPassword(input.Password)
	if err != nil {
		s.logger.Errorf("Failed to hash password: %v", err)
		_ = s.userRepo.DeleteUser(ctx, user.UserID)
		return nil, errors.NewInternalError("Failed to secure password", err)
	}

	verificationToken, err := token.GenerateVerificationToken()
	if err != nil {
		s.logger.Errorf("Failed to generate verification token: %v", err)
		_ = s.userRepo.DeleteUser(ctx, user.UserID)
		return nil, errors.NewInternalError("Failed to generate verification token", err)
	}

	authParams := repository.CreateAuthParams{
		UserID:            user.UserID,
		PasswordHash:      hashedPassword,
		Salt:              salt,
		IsEmailVerified:   false,
		VerificationToken: verificationToken,
	}

	err = s.authRepo.CreateAuth(ctx, authParams)
	if err != nil {
		s.logger.Errorf("Failed to create auth record: %v", err)
		_ = s.userRepo.DeleteUser(ctx, user.UserID)
		return nil, errors.Wrap(err, "Failed to create auth record")
	}

	// Send verification email
	err = s.SendVerificationEmail(ctx, user.UserID.String(), user.Username, verificationToken)
	if err != nil {
		s.logger.Warnf("Failed to send verification email: %v", err)
		// Non-critical error, continue with registration
	}

	s.logger.Infof("User successfully registered with ID: %s", user.UserID)
	return mapUserToDTO(user), nil
}

func (s *authService) Login(ctx context.Context, input LoginInput) (*TokenResponse, error) {
	s.logger.Infof("Login attempt for email: %s", input.Email)

	// Get user by email
	user, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		s.logger.Warnf("Login failed: email lookup error: %v", err)
		if errors.IsNotFound(err) {
			return nil, errors.NewUnauthorizedError("Invalid credentials", nil)
		}
		return nil, errors.Wrap(err, "Failed to retrieve user")
	}

	// Get auth record for user
	auth, err := s.authRepo.GetAuthByUserID(ctx, user.UserID)
	if err != nil {
		s.logger.Errorf("Error retrieving auth for user %s: %v", user.UserID, err)
		return nil, errors.Wrap(err, "Failed to retrieve authentication information")
	}

	// Check if account is locked
	if auth.LockedUntil != nil && time.Now().Before(*auth.LockedUntil) {
		s.logger.Warnf("Login attempt for locked account: %s until %v", user.UserID, *auth.LockedUntil)
		return nil, errors.NewForbiddenError("Account is temporarily locked", nil)
	}

	// Verify password
	err = password.VerifyPassword(input.Password, auth.PasswordHash, auth.Salt)
	if err != nil {
		s.logger.Warnf("Login failed: invalid password for user %s", user.UserID)

		// Increment failed login attempts
		errIncrement := s.authRepo.IncrementFailedLoginAttempts(ctx, user.UserID)
		if errIncrement != nil {
			s.logger.Errorf("Failed to increment login attempts: %v", errIncrement)
		}

		// Check if account should be locked
		if auth.FailedLoginAttempts != nil && *auth.FailedLoginAttempts >= 5 {
			lockUntil := time.Now().Add(15 * time.Minute)
			s.logger.Warnf("Locking account %s until %v due to multiple failed attempts", user.UserID, lockUntil)

			errLock := s.authRepo.SetAccountLockout(ctx, user.UserID, lockUntil)
			if errLock != nil {
				s.logger.Errorf("Failed to lock account: %v", errLock)
			}

			// Send account locked notification
			_ = s.SendAccountLockedEmail(ctx, user.UserID.String(), user.Username, lockUntil.Format(time.RFC1123))
		}

		return nil, errors.NewUnauthorizedError("Invalid credentials", nil)
	}

	// Update last login time
	err = s.authRepo.UpdateLastLogin(ctx, user.UserID)
	if err != nil {
		s.logger.Warnf("Failed to update last login: %v", err)
		// Non-critical error, continue with login
	}

	// Create token pair
	isAdmin := user.IsAdmin != nil && *user.IsAdmin
	isPremium := user.IsPremium != nil && *user.IsPremium

	jwtMaker := token.NewJWTMaker(s.jwtSecret)
	tokenPair, err := jwtMaker.CreateTokenPair(
		user.UserID.String(),
		user.Username,
		user.Email,
		isAdmin,
		isPremium,
		s.tokenExpiry,
		DefaultRefreshTokenDuration,
	)

	if err != nil {
		s.logger.Errorf("Failed to create token pair: %v", err)
		return nil, errors.NewInternalError("Failed to generate authentication tokens", err)
	}

	// Store refresh token
	err = s.authRepo.StoreRefreshToken(ctx, user.UserID, tokenPair.RefreshToken)
	if err != nil {
		s.logger.Errorf("Failed to store refresh token: %v", err)
		return nil, errors.Wrap(err, "Failed to store refresh token")
	}

	s.logger.Infof("User %s logged in successfully", user.UserID)
	return &TokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		User:         mapUserToDTO(user),
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, input RefreshTokenRequest) (*TokenResponse, error) {
	s.logger.Debugf("Processing refresh token request")

	// Verify refresh token
	jwtMaker := token.NewJWTMaker(s.jwtSecret)
	claims, err := jwtMaker.VerifyToken(input.RefreshToken)
	if err != nil {
		s.logger.Warnf("Invalid refresh token: %v", err)
		return nil, errors.NewUnauthorizedError("Invalid refresh token", err)
	}

	if claims.TokenType != token.RefreshToken {
		s.logger.Warnf("Wrong token type provided for refresh: %s", claims.TokenType)
		return nil, errors.NewUnauthorizedError("Invalid token type", nil)
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		s.logger.Errorf("Invalid user ID in token: %v", err)
		return nil, errors.NewInternalError("Invalid user identifier in token", err)
	}

	// Get user and auth info
	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		s.logger.Warnf("User from refresh token not found: %s", userID)
		return nil, errors.NewUnauthorizedError("Invalid refresh token", nil)
	}

	auth, err := s.authRepo.GetAuthByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get auth for refresh token: %v", err)
		return nil, errors.NewUnauthorizedError("Invalid refresh token", nil)
	}

	// Verify refresh token matches stored token
	if auth.RefreshToken == nil || *auth.RefreshToken != input.RefreshToken {
		s.logger.Warnf("Refresh token doesn't match stored token for user %s", userID)
		return nil, errors.NewUnauthorizedError("Refresh token has been invalidated", nil)
	}

	// Create new token pair
	isAdmin := user.IsAdmin != nil && *user.IsAdmin
	isPremium := user.IsPremium != nil && *user.IsPremium

	tokenPair, err := jwtMaker.CreateTokenPair(
		user.UserID.String(),
		user.Username,
		user.Email,
		isAdmin,
		isPremium,
		s.tokenExpiry,
		DefaultRefreshTokenDuration,
	)
	if err != nil {
		s.logger.Errorf("Failed to create new token pair: %v", err)
		return nil, errors.NewInternalError("Failed to generate new tokens", err)
	}

	// Store new refresh token
	err = s.authRepo.StoreRefreshToken(ctx, userID, tokenPair.RefreshToken)
	if err != nil {
		s.logger.Errorf("Failed to store new refresh token: %v", err)
		return nil, errors.Wrap(err, "Failed to store new refresh token")
	}

	s.logger.Infof("Token refreshed successfully for user %s", userID)
	return &TokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		User:         mapUserToDTO(user),
	}, nil
}

func (s *authService) GenerateResetToken(ctx context.Context, email string) error {
	s.logger.Infof("Generating password reset token for email: %s", email)

	// Find user by email
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.IsNotFound(err) {
			// Don't reveal that email doesn't exist for security reasons
			s.logger.Infof("Reset token requested for non-existent email: %s", email)
			return nil
		}
		s.logger.Errorf("Error looking up user by email: %v", err)
		return errors.Wrap(err, "Failed to lookup user email")
	}

	// Generate reset token
	resetToken, err := token.GenerateResetToken()
	if err != nil {
		s.logger.Errorf("Failed to generate reset token: %v", err)
		return errors.NewInternalError("Failed to generate reset token", err)
	}

	// Store reset token with expiration
	expiresAt := time.Now().Add(DefaultResetTokenDuration)
	err = s.authRepo.SetResetToken(ctx, user.UserID, resetToken, expiresAt)
	if err != nil {
		s.logger.Errorf("Failed to store reset token: %v", err)
		return errors.Wrap(err, "Failed to store reset token")
	}

	// Send password reset email
	err = s.SendPasswordResetEmail(ctx, user.UserID.String(), user.Username, resetToken)
	if err != nil {
		s.logger.Errorf("Failed to send password reset email: %v", err)
		return errors.Wrap(err, "Failed to send password reset email")
	}

	s.logger.Infof("Reset token generated successfully for user ID: %s", user.UserID)
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, input ResetPasswordInput) error {
	s.logger.Infof("Processing password reset for email: %s", input.Email)

	// Validate passwords match
	if input.NewPassword != input.ConfirmPassword {
		s.logger.Warnf("Password reset failed: passwords don't match")
		return errors.NewValidationError("Passwords do not match", nil)
	}

	// Validate password strength
	err := password.ValidatePassword(input.NewPassword, password.DefaultPasswordConfig())
	if err != nil {
		s.logger.Warnf("Password reset failed: password validation failed: %v", err)
		return errors.NewValidationError("Password does not meet requirements", err)
	}

	// Find user by email
	user, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		s.logger.Warnf("Password reset failed: email not found: %s", input.Email)
		return errors.NewNotFoundError("Invalid email address", err)
	}

	// Get auth record to verify token
	auth, err := s.authRepo.GetAuthByUserID(ctx, user.UserID)
	if err != nil {
		s.logger.Errorf("Failed to get auth record: %v", err)
		return errors.Wrap(err, "Failed to verify reset token")
	}

	// Verify reset token is valid
	if auth.ResetToken == nil || *auth.ResetToken != input.Token {
		s.logger.Warnf("Password reset failed: invalid token for user %s", user.UserID)
		return errors.NewUnauthorizedError("Invalid reset token", nil)
	}

	// Verify token hasn't expired
	if auth.ResetTokenExpiresAt == nil || time.Now().After(*auth.ResetTokenExpiresAt) {
		s.logger.Warnf("Password reset failed: expired token for user %s", user.UserID)
		return errors.NewUnauthorizedError("Reset token has expired", nil)
	}

	// Hash new password
	newHash, newSalt, err := password.HashPassword(input.NewPassword)
	if err != nil {
		s.logger.Errorf("Failed to hash new password: %v", err)
		return errors.NewInternalError("Failed to secure new password", err)
	}

	// Update password
	err = s.authRepo.UpdatePassword(ctx, user.UserID, newHash, newSalt)
	if err != nil {
		s.logger.Errorf("Failed to update password: %v", err)
		return errors.Wrap(err, "Failed to update password")
	}

	// Clear reset token
	err = s.authRepo.ClearResetToken(ctx, user.UserID)
	if err != nil {
		s.logger.Warnf("Failed to clear reset token: %v", err)
		// Non-critical error, password was updated successfully
	}

	// Send password changed confirmation email
	err = s.SendPasswordChangedEmail(ctx, user.UserID.String(), user.Username)
	if err != nil {
		s.logger.Warnf("Failed to send password changed notification: %v", err)
		// Non-critical error, password was reset successfully
	}

	s.logger.Infof("Password reset successfully for user ID: %s", user.UserID)
	return nil
}

// func (s *authService) VerifyEmail(ctx context.Context, verificationToken string) error {
//     s.logger.Infof("Processing email verification for token: %s", verificationToken)

//     // Find user with this verification token
//     // Note: You'll need to add a repository method to find auth by verification token
//     // Below is pseudo-code for the repository method you would need to add
//     /*
//     GetAuthByVerificationToken(ctx context.Context, token string) (*db.Auth, error)
//     */

//     // For now, let's assume you need to manually look up auth records
//     // This is inefficient and would be better with a direct DB query
//     // Get all users and find the one with the matching verification token

//     // Example implementation (assuming you add the necessary repository method):
//     auth, err := s.authRepo.GetAuthByVerificationToken(ctx, verificationToken)
//     if err != nil {
//         if errors.IsNotFound(err) {
//             s.logger.Warnf("Email verification failed: invalid token")
//             return errors.NewUnauthorizedError("Invalid verification token", nil)
//         }
//         s.logger.Errorf("Error retrieving auth record: %v", err)
//         return errors.Wrap(err, "Failed to verify email")
//     }

//     // Check if token is valid
//     if auth.VerificationToken == nil || *auth.VerificationToken != verificationToken {
//         s.logger.Warnf("Email verification failed: invalid token")
//         return errors.NewUnauthorizedError("Invalid verification token", nil)
//     }

//     // Check if email is already verified
//     if auth.IsEmailVerified != nil && *auth.IsEmailVerified {
//         s.logger.Infof("Email already verified for user ID: %s", auth.UserID)
//         return nil // Already verified, no error
//     }

//     // Verify the email
//     err = s.authRepo.VerifyEmail(ctx, auth.UserID)
//     if err != nil {
//         s.logger.Errorf("Failed to verify email: %v", err)
//         return errors.Wrap(err, "Failed to verify email")
//     }

//     // Get user information for confirmation email
//     user, err := s.userRepo.GetUser(ctx, auth.UserID)
//     if err != nil {
//         s.logger.Errorf("Failed to retrieve user: %v", err)
//         return errors.Wrap(err, "Failed to retrieve user")
//     }

//     data := map[string]interface{}{
//         "Username": user.Username,
//         "AppName":  "Your App Name",
//         "Year":     time.Now().Year(),
//     }

//     emailErr := s.emailClient.SendTemplate(
//         []string{user.Email},
//         "Email Verification Successful",
//         "verification_success.html",
//         data,
//     )
//     if emailErr != nil {
//         s.logger.Warnf("Failed to send verification confirmation email: %v", emailErr)
//     }

//     s.logger.Infof("Email verified successfully for user ID: %s", auth.UserID)
//     return nil
// }

func (s *authService) Logout(ctx context.Context, userIDStr string) error {
	s.logger.Infof("Processing logout for user ID: %s", userIDStr)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		s.logger.Warnf("Logout failed: invalid user ID format: %v", err)
		return errors.NewValidationError("Invalid user ID format", err)
	}

	err = s.authRepo.InvalidateRefreshToken(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to invalidate refresh token: %v", err)
		return errors.Wrap(err, "Failed to complete logout")
	}

	s.logger.Infof("User %s logged out successfully", userID)
	return nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenStr string) (*token.Claims, error) {
	s.logger.Debugf("Validating token")

	// Verify token signature and expiration
	jwtMaker := token.NewJWTMaker(s.jwtSecret)
	claims, err := jwtMaker.VerifyToken(tokenStr)
	if err != nil {
		s.logger.Warnf("Token validation failed: %v", err)
		return nil, errors.NewUnauthorizedError("Invalid token", err)
	}

	// Check token type
	if claims.TokenType != token.AccessToken {
		s.logger.Warnf("Wrong token type: %s", claims.TokenType)
		return nil, errors.NewUnauthorizedError("Invalid token type", nil)
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		s.logger.Warnf("Invalid user ID in token: %v", err)
		return nil, errors.NewUnauthorizedError("Invalid token", err)
	}

	// Verify user exists
	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Warnf("Token validation failed: user not found: %s", userID)
			return nil, errors.NewUnauthorizedError("User not found", nil)
		}
		s.logger.Errorf("Error retrieving user for token validation: %v", err)
		return nil, errors.Wrap(err, "Failed to validate user")
	}

	s.logger.Debugf("Token validated successfully for user %s", userID)
	return claims, nil
}

func (s *authService) IsEmailVerified(ctx context.Context, userIDStr string) (bool, error) {
	s.logger.Debugf("Checking email verification status for user: %s", userIDStr)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		s.logger.Warnf("Invalid user ID format: %v", err)
		return false, errors.NewValidationError("Invalid user ID format", err)
	}

	auth, err := s.authRepo.GetAuthByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get auth record: %v", err)
		return false, errors.Wrap(err, "Failed to check email verification status")
	}

	isVerified := auth.IsEmailVerified != nil && *auth.IsEmailVerified
	s.logger.Debugf("Email verification status for user %s: %v", userID, isVerified)
	return isVerified, nil
}

func (s *authService) SendVerificationEmail(ctx context.Context, email, username, token string) error {
	s.logger.Infof("Sending verification email to: %s", email)

	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", s.baseURL, token)

	data := map[string]interface{}{
		"Username": username,
		"Link":     verificationLink,
		"AppName":  "Your App Name",
		"Year":     time.Now().Year(),
	}

	err := s.emailClient.SendTemplate([]string{email}, "Verify Your Email", "verification.html", data)
	if err != nil {
		s.logger.Errorf("Failed to send verification email: %v", err)
		return errors.Wrap(err, "Failed to send verification email")
	}

	s.logger.Infof("Verification email sent successfully to: %s", email)
	return nil
}

func (s *authService) SendPasswordResetEmail(ctx context.Context, email, username, token string) error {
	s.logger.Infof("Sending password reset email to: %s", email)

	resetLink := fmt.Sprintf("%s/reset-password?token=%s&email=%s", s.baseURL, token, url.QueryEscape(email))

	data := map[string]interface{}{
		"Username": username,
		"Link":     resetLink,
		"AppName":  "Your App Name",
		"Year":     time.Now().Year(),
	}

	err := s.emailClient.SendTemplate([]string{email}, "Reset Your Password", "password_reset.html", data)
	if err != nil {
		s.logger.Errorf("Failed to send password reset email: %v", err)
		return errors.Wrap(err, "Failed to send password reset email")
	}

	s.logger.Infof("Password reset email sent successfully to: %s", email)
	return nil
}

func (s *authService) SendPasswordChangedEmail(ctx context.Context, email, username string) error {
	s.logger.Infof("Sending password changed notification to: %s", email)

	data := map[string]interface{}{
		"Username": username,
		"AppName":  "Your App Name",
		"Year":     time.Now().Year(),
	}

	err := s.emailClient.SendTemplate([]string{email}, "Your Password Has Been Changed", "password_changed.html", data)
	if err != nil {
		s.logger.Errorf("Failed to send password changed email: %v", err)
		return errors.Wrap(err, "Failed to send password changed notification")
	}

	s.logger.Infof("Password changed notification sent successfully to: %s", email)
	return nil
}

func (s *authService) SendAccountLockedEmail(ctx context.Context, email, username, unlockTime string) error {
	s.logger.Infof("Sending account locked notification to: %s", email)

	data := map[string]interface{}{
		"Username":   username,
		"AppName":    "Your App Name",
		"Year":       time.Now().Year(),
		"UnlockTime": unlockTime,
	}

	err := s.emailClient.SendTemplate([]string{email}, "Your Account Has Been Temporarily Locked", "account_locked.html", data)
	if err != nil {
		s.logger.Errorf("Failed to send account locked email: %v", err)
		return errors.Wrap(err, "Failed to send account locked notification")
	}

	s.logger.Infof("Account locked notification sent successfully to: %s", email)
	return nil
}

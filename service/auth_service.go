package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xsj/gin-sqlc/api"
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
	VerifyEmail(ctx context.Context, token string) error
	Logout(ctx context.Context, userID string) error
	ValidateToken(ctx context.Context, tokenStr string) (*token.Claims, error)
	IsEmailVerified(ctx context.Context, userID string) (bool, error)
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
	jwtSecret   string
	tokenExpiry time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	authRepo repository.AuthRepository,
	jwtSecret string,
	tokenExpiry time.Duration,
) AuthService {
	if tokenExpiry == 0 {
		tokenExpiry = DefaultAccessTokenDuration
	}

	return &authService{
		userRepo:    userRepo,
		authRepo:    authRepo,
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
	}
}

func (s *authService) Register(ctx context.Context, input RegisterInput) (*UserDTO, error) {
	err := password.ValidatePassword(input.Password, password.DefaultPasswordConfig())
	if err != nil {
		return nil, api.ErrInvalidInput
	}

	_, err = s.userRepo.GetUserByEmail(ctx, input.Email)
	if err == nil {
		return nil, api.ErrDuplicateEntry
	} else if !errors.Is(err, repository.ErrRecordNotFound) {
		return nil, api.ErrInternalServer
	}

	_, err = s.userRepo.GetUserByUsername(ctx, input.Username)
	if err == nil {
		return nil, api.ErrDuplicateEntry
	} else if !errors.Is(err, repository.ErrRecordNotFound) {
		return nil, api.ErrInternalServer
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
		if errors.Is(err, repository.ErrDuplicateKey) {
			return nil, api.ErrDuplicateEntry
		}
		return nil, api.ErrInternalServer
	}

	hashedPassword, salt, err := password.HashPassword(input.Password)
	if err != nil {
		_ = s.userRepo.DeleteUser(ctx, user.UserID)
		return nil, api.ErrInternalServer
	}

	verificationToken, err := token.GenerateVerificationToken()
	if err != nil {
		_ = s.userRepo.DeleteUser(ctx, user.UserID)
		return nil, api.ErrInternalServer
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
		_ = s.userRepo.DeleteUser(ctx, user.UserID)
		return nil, api.ErrInternalServer
	}

	return mapUserToDTO(user), nil
}

func (s *authService) Login(ctx context.Context, input LoginInput) (*TokenResponse, error) {
	fmt.Println("Login service called with email:", input.Email)
	user, err := s.userRepo.GetUserByEmail(ctx, input.Email)

	if err != nil {
		fmt.Printf("Error getting user by email: %v\n", err)
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, api.ErrUnauthorized
		}
		return nil, api.ErrInternalServer
	}
	fmt.Printf("User found with ID: %s\n", user.UserID)

	auth, err := s.authRepo.GetAuthByUserID(ctx, user.UserID)
	if err != nil {
		fmt.Printf("Error in authRepo.GetAuthByUserID: %v\n", err)
		return nil, api.ErrInternalServer
	}

	fmt.Printf("Auth found for user: %s, Salt length: %d\n", user.UserID, len(auth.Salt))

	if auth.LockedUntil != nil && time.Now().Before(*auth.LockedUntil) {
		return nil, api.ErrForbidden
	}

	err = password.VerifyPassword(input.Password, auth.PasswordHash, auth.Salt)
	if err != nil {
		_ = s.authRepo.IncrementFailedLoginAttempts(ctx, user.UserID)

		if auth.FailedLoginAttempts != nil && *auth.FailedLoginAttempts >= 5 {
			lockUntil := time.Now().Add(15 * time.Minute)
			_ = s.authRepo.SetAccountLockout(ctx, user.UserID, lockUntil)
		}
		return nil, api.ErrUnauthorized
	}

	_ = s.authRepo.UpdateLastLogin(ctx, user.UserID)

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
		return nil, api.ErrInternalServer
	}

	err = s.authRepo.StoreRefreshToken(ctx, user.UserID, tokenPair.RefreshToken)
	if err != nil {
		return nil, api.ErrInternalServer
	}

	return &TokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		User:         mapUserToDTO(user),
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, input RefreshTokenRequest) (*TokenResponse, error) {
	jwtMaker := token.NewJWTMaker(s.jwtSecret)
	claims, err := jwtMaker.VerifyToken(input.RefreshToken)
	if err != nil {
		return nil, api.ErrUnauthorized
	}

	if claims.TokenType != token.RefreshToken {
		return nil, api.ErrUnauthorized
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, api.ErrInternalServer
	}

	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, api.ErrUnauthorized
	}

	auth, err := s.authRepo.GetAuthByUserID(ctx, userID)
	if err != nil {
		return nil, api.ErrUnauthorized
	}

	if auth.RefreshToken == nil || *auth.RefreshToken != input.RefreshToken {
		return nil, api.ErrUnauthorized
	}

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
		return nil, api.ErrInternalServer
	}

	err = s.authRepo.StoreRefreshToken(ctx, userID, tokenPair.RefreshToken)
	if err != nil {
		return nil, api.ErrInternalServer
	}

	return &TokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		User:         mapUserToDTO(user),
	}, nil
}

func (s *authService) GenerateResetToken(ctx context.Context, email string) error {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil
		}
		return api.ErrInternalServer
	}

	resetToken, err := token.GenerateResetToken()
	if err != nil {
		return api.ErrInternalServer
	}

	expiresAt := time.Now().Add(DefaultRefreshTokenDuration)

	err = s.authRepo.SetResetToken(ctx, user.UserID, resetToken, expiresAt)
	if err != nil {
		return api.ErrInternalServer
	}
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, input ResetPasswordInput) error {
	if input.NewPassword != input.ConfirmPassword {
		return api.ErrInvalidInput
	}

	err := password.ValidatePassword(input.NewPassword, password.DefaultPasswordConfig())
	if err != nil {
		return api.ErrInvalidInput
	}

	user, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return api.ErrNotFound
	}

	auth, err := s.authRepo.GetAuthByUserID(ctx, user.UserID)
	if err != nil {
		return api.ErrInternalServer
	}

	if auth.ResetToken == nil || *auth.ResetToken != input.Token {
		return api.ErrUnauthorized
	}

	if auth.ResetTokenExpiresAt == nil || time.Now().After(*auth.ResetTokenExpiresAt) {
		return api.ErrUnauthorized
	}

	newHash, newSalt, err := password.HashPassword(input.NewPassword)
	if err != nil {
		return api.ErrInternalServer
	}

	err = s.authRepo.UpdatePassword(ctx, user.UserID, newHash, newSalt)
	if err != nil {
		return api.ErrInternalServer
	}

	err = s.authRepo.ClearResetToken(ctx, user.UserID)
	if err != nil {
		return api.ErrInternalServer
	}
	return nil
}

func (s *authService) VerifyEmail(ctx context.Context, verificationToken string) error {
	// This is a simplified implementation
	// You would typically need to lookup the user by verification token
	// Since our schema doesn't have a direct query for this, we would need to
	// iterate through users or add a new query

	// For demonstration, assume we can get the user from the token
	// In a real implementation, you would need to query the database

	// Verifying the token is valid and not expired
	// ...

	// Update user's email verification status
	// userID := ... // From token lookup
	// err := s.authRepo.VerifyEmail(ctx, userID)
	// if err != nil {
	//     return api.ErrInternalServer
	// }

	return fmt.Errorf("not implemented: requires additional query for token lookup")
}

func (s *authService) Logout(ctx context.Context, userIDStr string) error {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return api.ErrInvalidInput
	}

	err = s.authRepo.InvalidateRefreshToken(ctx, userID)
	if err != nil {
		return api.ErrInternalServer
	}
	return nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenStr string) (*token.Claims, error) {
	jwtMaker := token.NewJWTMaker(s.jwtSecret)
	claims, err := jwtMaker.VerifyToken(tokenStr)
	if err != nil {
		return nil, api.ErrUnauthorized
	}

	if claims.TokenType != token.AccessToken {
		return nil, api.ErrUnauthorized
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, api.ErrUnauthorized
		}
		return nil, api.ErrInternalServer
	}

	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, api.ErrUnauthorized
		}

		return nil, api.ErrInternalServer
	}

	return claims, nil
}

func (s *authService) IsEmailVerified(ctx context.Context, userIDStr string) (bool, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return false, api.ErrInvalidInput
	}

	auth, err := s.authRepo.GetAuthByUserID(ctx, userID)
	if err != nil {
		return false, api.ErrInternalServer
	}

	return auth.IsEmailVerified != nil && *auth.IsEmailVerified, nil
}

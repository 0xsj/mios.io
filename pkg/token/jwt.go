package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Claims struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	IsPremium bool      `json:"is_premium"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
	fmt.Printf("Creating JWT Maker with secret key length: %d\n", len(secretKey))
	return &JWTMaker{
		secretKey: secretKey,
	}
}

func (maker *JWTMaker) CreateToken(
	userID string,
	username string,
	email string,
	isAdmin bool,
	isPremium bool,
	tokenType TokenType,
	duration time.Duration,
) (string, time.Time, error) {
	fmt.Printf("Creating token for user: %s with token type: %s\n", userID, tokenType)

	expiresAt := time.Now().Add(duration)

	claims := Claims{
		UserID:    userID,
		Username:  username,
		Email:     email,
		IsAdmin:   isAdmin,
		IsPremium: isPremium,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}

	fmt.Printf("Claims created, about to create token with signing method HS256\n")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	fmt.Printf("Token created, attempting to sign with secret key length: %d\n", len(maker.secretKey))

	tokenString, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		fmt.Printf("Error signing token: %v\n", err)
		return "", time.Time{}, err
	}

	fmt.Printf("Token signed successfully, length: %d\n", len(tokenString))
	return tokenString, expiresAt, nil
}

func (maker *JWTMaker) CreateTokenPair(
	userID string,
	username string,
	email string,
	isAdmin bool,
	isPremium bool,
	accessDuration time.Duration,
	refreshDuration time.Duration,
) (*TokenPair, error) {
	fmt.Printf("Creating token pair for user: %s\n", userID)

	accessToken, expiresAt, err := maker.CreateToken(
		userID, username, email, isAdmin, isPremium, AccessToken, accessDuration,
	)
	if err != nil {
		fmt.Printf("Failed to create access token: %v\n", err)
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	fmt.Printf("Access token created successfully\n")

	refreshToken, _, err := maker.CreateToken(
		userID, username, email, isAdmin, isPremium, RefreshToken, refreshDuration,
	)
	if err != nil {
		fmt.Printf("Failed to create refresh token: %v\n", err)
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	fmt.Printf("Refresh token created successfully\n")

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Unix(),
	}, nil
}

func (maker *JWTMaker) VerifyToken(tokenString string) (*Claims, error) {
	fmt.Printf("Verifying token\n")

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			fmt.Printf("Unexpected signing method: %v\n", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(maker.secretKey), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, keyFunc)
	if err != nil {
		fmt.Printf("Error parsing token: %v\n", err)
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		fmt.Printf("Invalid token claims\n")
		return nil, errors.New("invalid token claims")
	}

	fmt.Printf("Token verified successfully\n")
	return claims, nil
}

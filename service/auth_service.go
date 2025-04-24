package service

import (
	"time"

	"github.com/0xsj/gin-sqlc/repository"
)

type AuthService interface {}

type RegisterInput struct {}

type LoginInput struct {}

type ResetPasswordInput struct {}

type TokenResponse struct {}

type authService struct {
	userRepo repository.UserRepository
	authRepo repository.AuthRepository
	jwtSecret string
	tokenExpiry time.Duration
}

func NewAuthService(){}

func Register() {}

func Login() {}

func GenerateResetToken(){}

func ResetPassword(){}

func VerifyEmail(){}

func RefreshToken(){}

func generateTokens(){}

func generateRandomString(){}

func hashPassword(){}

func verifyPassword(){}

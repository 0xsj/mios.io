package service

import (
	"context"
	"regexp"
	"strings"

	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/0xsj/gin-sqlc/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, input CreateUserInput) (*UserDTO, error)
	GetUser(ctx context.Context, id string) (*UserDTO, error)
	GetUserByUsername(ctx context.Context, username string) (*UserDTO, error)
	GetUserByHandle(ctx context.Context, handle string) (*UserDTO, error)
	GetUserByEmail(ctx context.Context, email string) (*UserDTO, error)
	UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*UserDTO, error)
	UpdateHandle(ctx context.Context, id string, handle string) (*UserDTO, error)
	UpdatePremiumStatus(ctx context.Context, id string, isPremium bool) (*UserDTO, error)
	UpdateAdminStatus(ctx context.Context, id string, isAdmin bool) (*UserDTO, error)
	UpdateOnboardedStatus(ctx context.Context, id string, onboarded bool) (*UserDTO, error)
	DeleteUser(ctx context.Context, id string) error
}

type CreateUserInput struct {
	Username        string `json:"username" binding:"required"`
	Handle          string `json:"handle" binding:"required"`
	Email           string `json:"email" binding:"required"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Bio             string `json:"bio"`
	ProfileImageURL string `json:"profile_image_url"`
	LayoutVersion   string `json:"layout_version"`
	CustomDomain    string `json:"custom_domain"`
	IsPremium       bool   `json:"is_premium"`
	IsAdmin         bool   `json:"is_admin"`
	Onboarded       bool   `json:"onboarded"`
}

type UpdateUserInput struct {
	Username        *string `json:"username"`
	Email           *string `json:"email"`
	FirstName       *string `json:"first_name"`
	LastName        *string `json:"last_name"`
	Bio             *string `json:"bio"`
	ProfileImageURL *string `json:"profile_image_url"`
	LayoutVersion   *string `json:"layout_version"`
	CustomDomain    *string `json:"custom_domain"`
}

type UserDTO struct {
	ID              string `json:"id"`
	Username        string `json:"username"`
	Handle          string `json:"handle,omitempty"`
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
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
}


type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}


func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (*UserDTO, error) {}


func (s *userService) GetUser(ctx context.Context, id string) (*UserDTO, error) {}


func (s *userService) GetUserByUsername(ctx context.Context, username string) (*UserDTO, error) {}


func (s *userService) GetUserByHandle(ctx context.Context, handle string) (*UserDTO, error) {}


func (s *userService) GetUserByEmail(ctx context.Context, email string) (*UserDTO, error) {}


func (s *userService) UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*UserDTO, error) {}


func (s *userService) UpdateHandle(ctxt context.Context, id string, handle string) (*UserDTO, error) {}


func (s *userService) UpdatePremiumStatus(ctx context.Context, id string, isPremium bool) (*UserDTO, error) {}


func (s *userService) UpdateAdminStatus(ctx context.Context, id string, isAdmin bool) (*UserDTO, error) {}


func (s *userService) UpdateOnboardedStatus(ctx context.Context, id string, onboarded bool) (*UserDTO, error) {}


func (s *userService) DeleteUser(ctx context.Context, id string) error {}


func mapUserToDTO(user *db.User) *UserDTO {
	dto := &UserDTO{
		ID:        user.UserID.String(),
		Username:  user.Username,
		Email:     user.Email,
		IsPremium: user.IsPremium != nil && *user.IsPremium,
		IsAdmin:   user.IsAdmin != nil && *user.IsAdmin,
		Onboarded: user.Onboarded != nil && *user.Onboarded,
	}

	if user.Handle != nil {
		dto.Handle = *user.Handle
	}
	
	if user.FirstName != nil {
		dto.FirstName = *user.FirstName
	}
	
	if user.LastName != nil {
		dto.LastName = *user.LastName
	}
	
	if user.Bio != nil {
		dto.Bio = *user.Bio
	}
	
	if user.ProfileImageUrl != nil {
		dto.ProfileImageURL = *user.ProfileImageUrl
	}
	
	if user.LayoutVersion != nil {
		dto.LayoutVersion = *user.LayoutVersion
	}
	
	if user.CustomDomain != nil {
		dto.CustomDomain = *user.CustomDomain
	}
	
	if user.CreatedAt != nil {
		dto.CreatedAt = user.CreatedAt.String()
	}
	
	if user.UpdatedAt != nil {
		dto.UpdatedAt = user.UpdatedAt.String()
	}

	return dto
}

func isValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 30 {
		return false
	}

	pattern := `^[a-zA-Z0-9_]+$`
	match, _ := regexp.MatchString(pattern, username)
	return match
}

func isValidHandle(handle string) bool {
	if len(handle) < 2 || len(handle) > 30 {
		return false
	}

	pattern := `^[a-zA-Z0-9_-]+$`
	match, _ := regexp.MatchString(pattern, handle)
	return match
}

func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}
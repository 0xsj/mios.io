package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/0xsj/gin-sqlc/api"
	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/0xsj/gin-sqlc/repository"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, input CreateUserInput) (*UserDTO, error)
	GetUser(ctx context.Context, id string) (*UserDTO, error)
	GetUserByUsername(ctx context.Context, username string) (*UserDTO, error)
	GetUserByEmail(ctx context.Context, email string) (*UserDTO, error)
}

type CreateUserInput struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UpdateUserInput struct {
    Username  *string `json:"username"`
    Email     *string `json:"email"`
    FirstName *string `json:"first_name"`
    LastName  *string `json:"last_name"`
}

type UserDTO struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	IsPremium bool   `json:"is_premium"`
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (*UserDTO, error) {
	fmt.Println("UserService.CreateUser called")

	if !isValidUsername(input.Username) {
		return nil, api.ErrInvalidInput
	}

	if !isValidEmail(input.Email) {
		return nil, api.ErrInvalidInput
	}

	params := repository.CreateUserParams{
		Username:  input.Username,
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	fmt.Printf("Calling repository with params: %+v\n", params)

	user, err := s.userRepo.CreateUser(ctx, params)
	if err != nil {
		fmt.Printf("Repository error: %v\n", err)
		if err == repository.ErrDuplicateKey {
			return nil, api.ErrDuplicateEntry
		}
		return nil, api.ErrInternalServer
	}

	return mapUserToDTO(user), nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*UserDTO, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, api.ErrInvalidInput
	}

	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, api.ErrNotFound
		}
		return nil, api.ErrInternalServer
	}

	return mapUserToDTO(user), nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*UserDTO, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, api.ErrNotFound
		}
		return nil, api.ErrInternalServer
	}

	return mapUserToDTO(user), nil
}


func (s *userService) GetUserByEmail(ctx context.Context, email string) (*UserDTO, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, api.ErrNotFound
		}
		return nil, api.ErrInternalServer
	}

	return mapUserToDTO(user), nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*UserDTO, error) {
    userID, err := uuid.Parse(id)
    if err != nil {
        return nil, api.ErrInvalidInput
    }
    
    if input.Username != nil && !isValidUsername(*input.Username) {
        return nil, api.ErrInvalidInput
    }
    
    if input.Email != nil && !isValidEmail(*input.Email) {
        return nil, api.ErrInvalidInput
    }
    
    params := repository.UpdateUserParams{
        UserID: userID,
    }
    
    if input.Username != nil {
        params.Username = *input.Username
    }
    
    if input.Email != nil {
        params.Email = *input.Email
    }
    
    if input.FirstName != nil {
        params.FirstName = *input.FirstName
    }
    
    if input.LastName != nil {
        params.LastName = *input.LastName
    }
    
    
    err = s.userRepo.UpdateUser(ctx, params)
    if err != nil {
        if err == repository.ErrDuplicateKey {
            return nil, api.ErrDuplicateEntry
        }
        return nil, api.ErrInternalServer
    }
    
    updatedUser, err := s.userRepo.GetUser(ctx, userID)
    if err != nil {
        return nil, api.ErrInternalServer
    }
    
    return mapUserToDTO(updatedUser), nil
}


func mapUserToDTO(user *db.User) *UserDTO {
	dto := &UserDTO{
		ID:        user.UserID.String(),
		Username:  user.Username,
		Email:     user.Email,
		IsPremium: user.IsPremium != nil && *user.IsPremium,
	}

	if user.FirstName != nil {
		dto.FirstName = *user.FirstName
	}
	if user.LastName != nil {
		dto.LastName = *user.LastName
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

func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}

package service

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/0xsj/gin-sqlc/log"
	apperror "github.com/0xsj/gin-sqlc/pkg/errors"

	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/0xsj/gin-sqlc/repository"
	"github.com/google/uuid"
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
	logger   log.Logger
}

func NewUserService(userRepo repository.UserRepository, logger log.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func handleValidationError(message string, err error) error {
	return apperror.NewValidationError(message, err)
}

func parseUUID(id string) (uuid.UUID, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, apperror.NewBadRequestError("Invalid user ID format", err)
	}
	return userID, nil
}

func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (*UserDTO, error) {
	s.logger.Infof("Creating new user with username: %s, email: %s", input.Username, input.Email)

	if !isValidUsername(input.Username) {
		return nil, handleValidationError("Invalid username format", nil)
	}

	if !isValidEmail(input.Email) {
		return nil, handleValidationError("Invalid email format", nil)
	}

	if !isValidHandle(input.Handle) {
		return nil, handleValidationError("Invalid handle format", nil)
	}

	params := repository.CreateUserParams{
		Username:        input.Username,
		Handle:          input.Handle,
		Email:           input.Email,
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		Bio:             input.Bio,
		ProfileImageURL: input.ProfileImageURL,
		LayoutVersion:   input.LayoutVersion,
		CustomDomain:    input.CustomDomain,
		IsPremium:       input.IsPremium,
		IsAdmin:         input.IsAdmin,
		Onboarded:       input.Onboarded,
	}

	start := time.Now()
	user, err := s.userRepo.CreateUser(ctx, params)
	duration := time.Since(start)

	if err != nil {
		s.logger.Errorf("Failed to create user: %v", err)
		return nil, err 
	}

	s.logger.Infof("User created successfully in %v: %s", duration, user.UserID)
	return mapUserToDTO(user), nil
}


func (s *userService) GetUser(ctx context.Context, id string) (*UserDTO, error) {
	s.logger.Debugf("Getting user by ID: %s", id)

	userID, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	user, err := s.userRepo.GetUser(ctx, userID)
	duration := time.Since(start)

	if err != nil {
		s.logger.Errorf("Failed to get user by ID %s: %v", id, err)
		return nil, err
	}

	s.logger.Debugf("Retrieved user by ID %s in %v", id, duration)
	return mapUserToDTO(user), nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*UserDTO, error) {
	s.logger.Debugf("Getting user by username: %s", username)

	start := time.Now()
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	duration := time.Since(start)

	if err != nil {
		s.logger.Errorf("Failed to get user by username %s: %v", username, err)
		return nil, err 
	}

	s.logger.Debugf("Retrieved user by username %s in %v", username, duration)
	return mapUserToDTO(user), nil
}

func (s *userService) GetUserByHandle(ctx context.Context, handle string) (*UserDTO, error) {
	s.logger.Debugf("Getting user by handle: %s", handle)

	start := time.Now()
	user, err := s.userRepo.GetUserByHandle(ctx, handle)
	duration := time.Since(start)

	if err != nil {
		s.logger.Errorf("Failed to get user by handle %s: %v", handle, err)
		return nil, err
	}

	s.logger.Debugf("Retrieved user by handle %s in %v", handle, duration)
	return mapUserToDTO(user), nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*UserDTO, error) {
	s.logger.Debugf("Getting user by email: %s", email)

	start := time.Now()
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	duration := time.Since(start)

	if err != nil {
		s.logger.Errorf("Failed to get user by email %s: %v", email, err)
		return nil, err 
	}

	s.logger.Debugf("Retrieved user by email %s in %v", email, duration)
	return mapUserToDTO(user), nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*UserDTO, error) {
	s.logger.Infof("Updating user with ID: %s", id)

	userID, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	// First get the current user to check for existence
	start := time.Now()
	currentUser, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get user for update with ID %s: %v", id, err)
		return nil, err // Repository already returns appropriate app errors
	}

	// Update user's basic details
	params := repository.UpdateUserParams{
		UserID:          userID,
		FirstName:       getValueOrEmpty(input.FirstName),
		LastName:        getValueOrEmpty(input.LastName),
		Bio:             getValueOrEmpty(input.Bio),
		ProfileImageURL: getValueOrEmpty(input.ProfileImageURL),
		LayoutVersion:   getValueOrEmpty(input.LayoutVersion),
		CustomDomain:    getValueOrEmpty(input.CustomDomain),
	}

	err = s.userRepo.UpdateUser(ctx, params)
	if err != nil {
		s.logger.Errorf("Failed to update user details with ID %s: %v", id, err)
		return nil, err // Repository already returns appropriate app errors
	}

	// Update username if provided and different
	if input.Username != nil && *input.Username != currentUser.Username {
		if !isValidUsername(*input.Username) {
			return nil, handleValidationError("Invalid username format", nil)
		}

		err = s.userRepo.UpdateUsername(ctx, userID, *input.Username)
		if err != nil {
			s.logger.Errorf("Failed to update username for user ID %s: %v", id, err)
			return nil, err // Repository already returns appropriate app errors
		}
	}

	// Update email if provided and different
	if input.Email != nil && *input.Email != currentUser.Email {
		if !isValidEmail(*input.Email) {
			return nil, handleValidationError("Invalid email format", nil)
		}

		err = s.userRepo.UpdateEmail(ctx, userID, *input.Email)
		if err != nil {
			s.logger.Errorf("Failed to update email for user ID %s: %v", id, err)
			return nil, err // Repository already returns appropriate app errors
		}
	}

	// Get the updated user
	updatedUser, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get updated user with ID %s: %v", id, err)
		return nil, apperror.NewInternalError("Failed to retrieve updated user", err)
	}

	duration := time.Since(start)
	s.logger.Infof("User with ID %s updated successfully in %v", id, duration)
	return mapUserToDTO(updatedUser), nil
}

func (s *userService) UpdateHandle(ctx context.Context, id string, handle string) (*UserDTO, error) {
	s.logger.Infof("Updating handle for user ID: %s to: %s", id, handle)

	userID, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	if !isValidHandle(handle) {
		return nil, handleValidationError("Invalid handle format", nil)
	}

	start := time.Now()
	err = s.userRepo.UpdateHandle(ctx, userID, handle)
	if err != nil {
		s.logger.Errorf("Failed to update handle for user ID %s: %v", id, err)
		return nil, err // Repository already returns appropriate app errors
	}

	// Get the updated user
	updatedUser, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get updated user with ID %s: %v", id, err)
		return nil, apperror.NewInternalError("Failed to retrieve updated user", err)
	}

	duration := time.Since(start)
	s.logger.Infof("Handle for user ID %s updated successfully in %v", id, duration)
	return mapUserToDTO(updatedUser), nil
}

func (s *userService) UpdatePremiumStatus(ctx context.Context, id string, isPremium bool) (*UserDTO, error) {
	s.logger.Infof("Updating premium status for user ID: %s to: %v", id, isPremium)

	userID, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	err = s.userRepo.UpdatePremiumStatus(ctx, userID, isPremium)
	if err != nil {
		s.logger.Errorf("Failed to update premium status for user ID %s: %v", id, err)
		return nil, err // Repository already returns appropriate app errors
	}

	// Get the updated user
	updatedUser, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get updated user with ID %s: %v", id, err)
		return nil, apperror.NewInternalError("Failed to retrieve updated user", err)
	}

	duration := time.Since(start)
	s.logger.Infof("Premium status for user ID %s updated successfully in %v", id, duration)
	return mapUserToDTO(updatedUser), nil
}

func (s *userService) UpdateAdminStatus(ctx context.Context, id string, isAdmin bool) (*UserDTO, error) {
	s.logger.Infof("Updating admin status for user ID: %s to: %v", id, isAdmin)

	userID, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	err = s.userRepo.UpdateAdminStatus(ctx, userID, isAdmin)
	if err != nil {
		s.logger.Errorf("Failed to update admin status for user ID %s: %v", id, err)
		return nil, err // Repository already returns appropriate app errors
	}

	// Get the updated user
	updatedUser, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get updated user with ID %s: %v", id, err)
		return nil, apperror.NewInternalError("Failed to retrieve updated user", err)
	}

	duration := time.Since(start)
	s.logger.Infof("Admin status for user ID %s updated successfully in %v", id, duration)
	return mapUserToDTO(updatedUser), nil
}

func (s *userService) UpdateOnboardedStatus(ctx context.Context, id string, onboarded bool) (*UserDTO, error) {
	s.logger.Infof("Updating onboarded status for user ID: %s to: %v", id, onboarded)

	userID, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	err = s.userRepo.UpdateOnboardedStatus(ctx, userID, onboarded)
	if err != nil {
		s.logger.Errorf("Failed to update onboarded status for user ID %s: %v", id, err)
		return nil, err // Repository already returns appropriate app errors
	}

	// Get the updated user
	updatedUser, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get updated user with ID %s: %v", id, err)
		return nil, apperror.NewInternalError("Failed to retrieve updated user", err)
	}

	duration := time.Since(start)
	s.logger.Infof("Onboarded status for user ID %s updated successfully in %v", id, duration)
	return mapUserToDTO(updatedUser), nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	s.logger.Warnf("Deleting user with ID: %s", id)

	userID, err := parseUUID(id)
	if err != nil {
		return err
	}

	// First check if the user exists
	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to find user for deletion with ID %s: %v", id, err)
		return err // Repository already returns appropriate app errors
	}

	start := time.Now()
	err = s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to delete user with ID %s: %v", id, err)
		return err // Repository already returns appropriate app errors
	}

	duration := time.Since(start)
	s.logger.Warnf("User with ID %s deleted successfully in %v", id, duration)
	return nil
}

func mapUserToDTO(user *db.User) *UserDTO {
	dto := &UserDTO{
		ID:        user.UserID.String(),
		Username:  user.Username,
		Handle:    user.Handle,
		Email:     user.Email,
		IsPremium: user.IsPremium != nil && *user.IsPremium,
		IsAdmin:   user.IsAdmin != nil && *user.IsAdmin,
		Onboarded: user.Onboarded != nil && *user.Onboarded,
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

func getValueOrEmpty(ptr *string) string {
    if ptr == nil {
        return ""
    }
    return *ptr
}
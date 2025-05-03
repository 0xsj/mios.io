package repository

import (
	"context"
	"time"

	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/0xsj/gin-sqlc/log"
	apperror "github.com/0xsj/gin-sqlc/pkg/errors"
	"github.com/0xsj/gin-sqlc/pkg/ptr"
	"github.com/google/uuid"
)



type UserRepository interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (*db.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*db.User, error)
	GetUserByUsername(ctx context.Context, username string) (*db.User, error)
	GetUserByHandle(ctx context.Context, handle string) (*db.User, error)
	GetUserByEmail(ctx context.Context, email string) (*db.User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) error
	UpdateUsername(ctx context.Context, userID uuid.UUID, username string) error
	UpdateHandle(ctx context.Context, userID uuid.UUID, handle string) error
	UpdateEmail(ctx context.Context, userID uuid.UUID, email string) error
	UpdatePremiumStatus(ctx context.Context, userID uuid.UUID, isPremium bool) error
	UpdateAdminStatus(ctx context.Context, userID uuid.UUID, isAdmin bool) error
	UpdateOnboardedStatus(ctx context.Context, userID uuid.UUID, onboarded bool) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type CreateUserParams struct {
	Username        string
	Handle          string
	Email           string
	FirstName       string
	LastName        string
	ProfileImageURL string
	Bio             string
	LayoutVersion   string
	CustomDomain    string
	IsPremium       bool
	IsAdmin         bool
	Onboarded       bool
}

type UpdateUserParams struct {
	UserID          uuid.UUID
	FirstName       string
	LastName        string
	ProfileImageURL string
	Bio             string
	LayoutVersion   string
	CustomDomain    string
}

const (
	PgErrUniqueViolation     = "23505"
	PgErrForeignKeyViolation = "23503"
	PgErrCheckViolation      = "23514"
)


type SQLCUserRepository struct {
	db *db.Queries
	logger log.Logger
}

func NewUserRepository(db *db.Queries, logger log.Logger) UserRepository {
	return &SQLCUserRepository{
		db: db,
		logger: logger,
	}
}

func (r *SQLCUserRepository) CreateUser(ctx context.Context, arg CreateUserParams) (*db.User, error) {
	r.logger.Infof("Creating user with username: %s, email: %s", arg.Username, arg.Email)
	params := db.CreateUserParams{
		Username:        arg.Username,
		Handle:          arg.Handle,
		Email:           arg.Email,
		FirstName:       ptr.String(arg.FirstName),
		LastName:        ptr.String(arg.LastName),
		Bio:             ptr.String(arg.Bio),
		ProfileImageUrl: ptr.String(arg.ProfileImageURL),
		LayoutVersion:   ptr.String(arg.LayoutVersion),
		CustomDomain:    ptr.String(arg.CustomDomain),
		IsPremium:       ptr.Bool(arg.IsPremium),
		IsAdmin:         ptr.Bool(arg.IsAdmin),
		Onboarded:       ptr.Bool(arg.Onboarded),
	}

	if arg.Username == "" || arg.Email == "" {
		return nil, apperror.NewValidationError("username and email are required", nil)
	}

	start := time.Now()
    user, err := r.db.CreateUser(ctx, params)
    duration := time.Since(start)

	if err != nil {
        appErr := apperror.HandleDBError(err, "user")        
        appErr.Log(r.logger)
        
        return nil, appErr
    }
	r.logger.Infof("User created successfully in %v: %s", duration, user.UserID)
    return user, nil
}

func (r *SQLCUserRepository) GetUser(ctx context.Context, userID uuid.UUID) (*db.User, error) {
    r.logger.Debugf("Getting user with ID: %s", userID)
    
    start := time.Now()
    user, err := r.db.GetUser(ctx, userID)
    duration := time.Since(start)
    
    if err != nil {
        appErr := apperror.HandleDBError(err, "user")
        appErr.Log(r.logger)
        return nil, appErr
    }
    
    r.logger.Debugf("Retrieved user with ID: %s in %v", userID, duration)
    return user, nil
}


func (r *SQLCUserRepository) GetUserByUsername(ctx context.Context, username string) (*db.User, error) {
    r.logger.Debugf("Getting user by username: %s", username)
    
    start := time.Now()
    user, err := r.db.GetUserByUsername(ctx, username)
    duration := time.Since(start)
    
    if err != nil {
        appErr := apperror.HandleDBError(err, "user")
        appErr.Log(r.logger)
        return nil, appErr
    }
    
    r.logger.Debugf("Retrieved user by username: %s in %v", username, duration)
    return user, nil
}

func (r *SQLCUserRepository) GetUserByHandle(ctx context.Context, handle string) (*db.User, error) {
    r.logger.Debugf("Getting user by handle: %s", handle)
    
    start := time.Now()
    user, err := r.db.GetUserByHandle(ctx, handle)
    duration := time.Since(start)
    
    if err != nil {
        appErr := apperror.HandleDBError(err, "user")
        appErr.Log(r.logger)
        return nil, appErr
    }
    
    r.logger.Debugf("Retrieved user by handle: %s in %v", handle, duration)
    return user, nil
}

func (r *SQLCUserRepository) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
    r.logger.Debugf("Getting user by email: %s", email)
    
    start := time.Now()
    user, err := r.db.GetUserByEmail(ctx, email)
    duration := time.Since(start)
    
    if err != nil {
        appErr := apperror.HandleDBError(err, "user")
        appErr.Log(r.logger)
        return nil, appErr
    }
    
    r.logger.Debugf("Retrieved user by email: %s in %v", email, duration)
    return user, nil
}

func (r *SQLCUserRepository) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
    r.logger.Infof("Updating user with ID: %s", arg.UserID)
    
    params := db.UpdateUserParams{
        UserID:          arg.UserID,
        FirstName:       ptr.String(arg.FirstName),
        LastName:        ptr.String(arg.LastName),
        ProfileImageUrl: ptr.String(arg.ProfileImageURL),
        Bio:             ptr.String(arg.Bio),
        LayoutVersion:   ptr.String(arg.LayoutVersion),
        CustomDomain:    ptr.String(arg.CustomDomain),
    }

    start := time.Now()
    err := r.db.UpdateUser(ctx, params)
    duration := time.Since(start)
    
    if err != nil {
        appErr := apperror.HandleDBError(err, "user")
        appErr.Log(r.logger)
        return appErr
    }
    
    r.logger.Infof("Updated user with ID: %s in %v", arg.UserID, duration)
    return nil
}

func (r *SQLCUserRepository) UpdateUsername(ctx context.Context, userID uuid.UUID, username string) error {
    r.logger.Infof("Updating username for user ID: %s to: %s", userID, username)
    
    params := db.UpdateUsernameParams{
        UserID:   userID,
        Username: username,
    }

    start := time.Now()
    err := r.db.UpdateUsername(ctx, params)
    duration := time.Since(start)
    
    if err != nil {
        appErr := apperror.HandleDBError(err, "username")
        appErr.Log(r.logger)
        return appErr
    }
    
    r.logger.Infof("Updated username for user ID: %s in %v", userID, duration)
    return nil
}


func (r *SQLCUserRepository) UpdateHandle(ctx context.Context, userID uuid.UUID, handle string) error {
    r.logger.Infof("Updating handle for user ID: %s to: %s", userID, handle)
    
    params := db.UpdateHandleParams{
        UserID: userID,
        Handle: handle,
    }

    start := time.Now()
    err := r.db.UpdateHandle(ctx, params)
    duration := time.Since(start)
    
    if err != nil {
        appErr := apperror.HandleDBError(err, "handle")
        appErr.Log(r.logger)
        return appErr
    }
    
    r.logger.Infof("Updated handle for user ID: %s in %v", userID, duration)
    return nil
}

func (r *SQLCUserRepository) UpdateEmail(ctx context.Context, userID uuid.UUID, email string) error {
    r.logger.Infof("Updating email for user ID: %s to: %s", userID, email)
    
    params := db.UpdateEmailParams{
        UserID: userID,
        Email:  email,
    }

    start := time.Now()
    err := r.db.UpdateEmail(ctx, params)
    duration := time.Since(start)
    
    if err != nil {
        appErr := apperror.HandleDBError(err, "email")
        appErr.Log(r.logger)
        return appErr
    }
    
    r.logger.Infof("Updated email for user ID: %s in %v", userID, duration)
    return nil
}

func (r *SQLCUserRepository) UpdatePremiumStatus(ctx context.Context, userID uuid.UUID, isPremium bool) error {
    r.logger.Infof("Updating premium status for user ID: %s to: %v", userID, isPremium)
    
    params := db.UpdateUserPremiumStatusParams{
        UserID:    userID,
        IsPremium: ptr.Bool(isPremium),
    }

    start := time.Now()
    err := r.db.UpdateUserPremiumStatus(ctx, params)
    duration := time.Since(start)
    
    if err != nil {
        appErr := apperror.HandleDBError(err, "user")
        appErr.Log(r.logger)
        return appErr
    }
    
    r.logger.Infof("Updated premium status for user ID: %s in %v", userID, duration)
    return nil
}

func (r *SQLCUserRepository) UpdateAdminStatus(ctx context.Context, userID uuid.UUID, isAdmin bool) error {
    r.logger.Infof("Updating admin status for user ID: %s to: %v", userID, isAdmin)
    
    params := db.UpdateUserAdminStatusParams{
        UserID:  userID,
        IsAdmin: ptr.Bool(isAdmin),
    }

    start := time.Now()
    err := r.db.UpdateUserAdminStatus(ctx, params)
    duration := time.Since(start)
    
    if err != nil {
        appErr := apperror.HandleDBError(err, "user")
        appErr.Log(r.logger)
        return appErr
    }
    
    r.logger.Infof("Updated admin status for user ID: %s in %v", userID, duration)
    return nil
}

func (r *SQLCUserRepository) UpdateOnboardedStatus(ctx context.Context, userID uuid.UUID, onboarded bool) error {
    r.logger.Infof("Updating onboarded status for user ID: %s to: %v", userID, onboarded)
    
    params := db.UpdateUserOnboardedStatusParams{
        UserID:    userID,
        Onboarded: ptr.Bool(onboarded),
    }

    start := time.Now()
    err := r.db.UpdateUserOnboardedStatus(ctx, params)
    duration := time.Since(start)
    
    if err != nil {
        appErr := apperror.HandleDBError(err, "user")
        appErr.Log(r.logger)
        return appErr
    }
    
    r.logger.Infof("Updated onboarded status for user ID: %s in %v", userID, duration)
    return nil
}

func (r *SQLCUserRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
    r.logger.Warnf("Deleting user with ID: %s", userID)
    
    start := time.Now()
    err := r.db.DeleteUser(ctx, userID)
    duration := time.Since(start)
    
    if err != nil {
        appErr := apperror.HandleDBError(err, "user")
        appErr.Log(r.logger)
        return appErr
    }
    
    r.logger.Warnf("Deleted user with ID: %s in %v", userID, duration)
    return nil
}
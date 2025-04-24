package repository

import (
	"context"
	"errors"

	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateKey   = errors.New("duplicate key violation")
	ErrDatabase       = errors.New("database error")

	ErrInvalidInput        = errors.New("invalid input parameters")
	ErrPermissionDenied    = errors.New("permission denied for operation")
	ErrForeignKeyViolation = errors.New("foreign key constraint violation")
	ErrTransactionFailed   = errors.New("database transaction failed")
	ErrConnectionFailed    = errors.New("database connection failed")
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

type SQLCUserRepository struct {
	db *db.Queries
}

func NewUserRepository(db *db.Queries) UserRepository {
	return &SQLCUserRepository{
		db: db,
	}
}

func (r *SQLCUserRepository) CreateUser(ctx context.Context, arg CreateUserParams) (*db.User, error) {
	var firstNamePtr, lastNamePtr, bioPtr, profileImageURLPtr *string
	var layoutVersionPtr, customDomainPtr *string
	var isPremiumPtr, isAdminPtr, onboardedPtr *bool

	if arg.FirstName != "" {
		firstNameCopy := arg.FirstName
		firstNamePtr = &firstNameCopy
	}

	if arg.LastName != "" {
		lastNameCopy := arg.LastName
		lastNamePtr = &lastNameCopy
	}

	if arg.Bio != "" {
		bioCopy := arg.Bio
		bioPtr = &bioCopy
	}

	if arg.ProfileImageURL != "" {
		profileCopy := arg.ProfileImageURL
		profileImageURLPtr = &profileCopy
	}

	if arg.LayoutVersion != "" {
		layoutCopy := arg.LayoutVersion
		layoutVersionPtr = &layoutCopy
	}

	if arg.CustomDomain != "" {
		domainCopy := arg.CustomDomain
		customDomainPtr = &domainCopy
	}

	isPremiumCopy := arg.IsPremium
	isPremiumPtr = &isPremiumCopy

	isAdminCopy := arg.IsAdmin
	isAdminPtr = &isAdminCopy

	onboardedCopy := arg.Onboarded
	onboardedPtr = &onboardedCopy

	params := db.CreateUserParams{
		Username:        arg.Username,
		Handle:          arg.Handle,
		Email:           arg.Email,
		FirstName:       firstNamePtr,
		LastName:        lastNamePtr,
		Bio:             bioPtr,
		ProfileImageUrl: profileImageURLPtr,
		LayoutVersion:   layoutVersionPtr,
		CustomDomain:    customDomainPtr,
		IsPremium:       isPremiumPtr,
		IsAdmin:         isAdminPtr,
		Onboarded:       onboardedPtr,
	}

	user, err := r.db.CreateUser(ctx, params)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok {
			if pgErr.Code == "23505" {
				return nil, ErrDuplicateKey
			}
		}
		return nil, ErrDatabase
	}

	return user, nil
}

func (r *SQLCUserRepository) GetUser(ctx context.Context, userID uuid.UUID) (*db.User, error) {
	user, err := r.db.GetUser(ctx, userID)
	if err != nil {
		return nil, ErrRecordNotFound
	}

	return user, nil
}

func (r *SQLCUserRepository) GetUserByUsername(ctx context.Context, username string) (*db.User, error) {
	user, err := r.db.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, ErrRecordNotFound
	}
	return user, nil
}

func (r *SQLCUserRepository) GetUserByHandle(ctx context.Context, handle string) (*db.User, error) {
	user, err := r.db.GetUserByHandle(ctx, handle)
	if err != nil {
		return nil, ErrRecordNotFound
	}

	return user, nil
}

func (r *SQLCUserRepository) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	user, err := r.db.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrRecordNotFound
	}
	return user, nil
}

func (r *SQLCUserRepository) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	var firstNamePtr, lastNamePtr, profileImageURLPtr, bioPtr *string
	var layoutVersionPtr, customDomainPtr *string

	if arg.FirstName != "" {
		firstNamePtr = &arg.FirstName
	}

	if arg.LastName != "" {
		lastNamePtr = &arg.LastName
	}

	if arg.ProfileImageURL != "" {
		profileImageURLPtr = &arg.ProfileImageURL
	}

	if arg.Bio != "" {
		bioPtr = &arg.Bio
	}

	if arg.LayoutVersion != "" {
		layoutVersionPtr = &arg.LayoutVersion
	}

	if arg.CustomDomain != "" {
		customDomainPtr = &arg.CustomDomain
	}

	params := db.UpdateUserParams{
		UserID:          arg.UserID,
		FirstName:       firstNamePtr,
		LastName:        lastNamePtr,
		ProfileImageUrl: profileImageURLPtr,
		Bio:             bioPtr,
		LayoutVersion:   layoutVersionPtr,
		CustomDomain:    customDomainPtr,
	}

	err := r.db.UpdateUser(ctx, params)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok {
			if pgErr.Code == "23505" {
				return ErrDuplicateKey
			}
		}
		return ErrDatabase
	}

	return nil
}

func (r *SQLCUserRepository) UpdateUsername(ctx context.Context, userID uuid.UUID, username string) error {
	params := db.UpdateUsernameParams{
		UserID:   userID,
		Username: username,
	}

	err := r.db.UpdateUsername(ctx, params)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" {
			return ErrDuplicateKey
		}
		return ErrDatabase
	}
	return nil
}

func (r *SQLCUserRepository) UpdateHandle(ctx context.Context, userID uuid.UUID, handle string) error {
	params := db.UpdateHandleParams{
		UserID: userID,
		Handle: handle,
	}

	err := r.db.UpdateHandle(ctx, params)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" {
			return ErrDuplicateKey
		}
		return ErrDatabase
	}
	return nil
}

func (r *SQLCUserRepository) UpdateEmail(ctx context.Context, userID uuid.UUID, email string) error {
	params := db.UpdateEmailParams{
		UserID: userID,
		Email:  email,
	}

	err := r.db.UpdateEmail(ctx, params)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" {
			return ErrDuplicateKey
		}
		return ErrDatabase
	}
	return nil
}

func (r *SQLCUserRepository) UpdatePremiumStatus(ctx context.Context, userID uuid.UUID, isPremium bool) error {
	params := db.UpdateUserPremiumStatusParams{
		UserID:    userID,
		IsPremium: &isPremium,
	}

	err := r.db.UpdateUserPremiumStatus(ctx, params)
	if err != nil {
		return ErrDatabase
	}
	return nil
}

func (r *SQLCUserRepository) UpdateAdminStatus(ctx context.Context, userID uuid.UUID, isAdmin bool) error {
	params := db.UpdateUserAdminStatusParams{
		UserID:  userID,
		IsAdmin: &isAdmin,
	}

	err := r.db.UpdateUserAdminStatus(ctx, params)
	if err != nil {
		return ErrDatabase
	}
	return nil
}

func (r *SQLCUserRepository) UpdateOnboardedStatus(ctx context.Context, userID uuid.UUID, onboarded bool) error {
	params := db.UpdateUserOnboardedStatusParams{
		UserID:    userID,
		Onboarded: &onboarded,
	}

	err := r.db.UpdateUserOnboardedStatus(ctx, params)
	if err != nil {
		return ErrDatabase
	}

	return nil
}

func (r *SQLCUserRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	err := r.db.DeleteUser(ctx, userID)
	if err != nil {
		return ErrDatabase
	}

	return nil
}

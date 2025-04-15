package repository

import (
	"context"
	"errors"
	"fmt"

	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateKey   = errors.New("duplicate key violation")
	ErrDatabase       = errors.New("database error")
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (*db.User, error)
	// GetUser(ctx context.Context, userID uuid.UUID) (db.User, error)
	// GetUserByUsername(ctx context.Context, username string) (db.User, error)
	// GetUserByEmail(ctx context.Context, email string) (db.User, error)
	// UpdateUser(ctx context.Context, arg UpdateUserParams) error
	// UpdateUsername(ctx context.Context, userID uuid.UUID, username string) error
	// UpdateEmail(ctx context.Context, userID uuid.UUID, email string) error
	// UpdatePremiumStatus(ctx context.Context, userID uuid.UUID, isPremium bool) error
	// DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type CreateUserParams struct {
	Username        string
	Email           string
	FirstName       string
	LastName        string
	ProfileImageURL string
	Bio             string
	Theme           string
}

type UpdateUserParams struct {
	UserID          uuid.UUID
	FirstName       string
	LastName        string
	ProfileImageURL string
	Bio             string
	Theme           string
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
	fmt.Println("UserRepository.CreateUser called")

	var firstNamePtr, lastNamePtr, profileImagePtr, bioPtr, themePtr *string

	if arg.FirstName != "" {
		firstNamePtr = &arg.FirstName
	}
	if arg.LastName != "" {
		lastNamePtr = &arg.LastName
	}
	if arg.ProfileImageURL != "" {
		profileImagePtr = &arg.ProfileImageURL
	}
	if arg.Bio != "" {
		bioPtr = &arg.Bio
	}
	if arg.Theme != "" {
		themePtr = &arg.Theme
	}

	params := db.CreateUserParams{
		Username:        arg.Username,
		Email:           arg.Email,
		FirstName:       firstNamePtr,
		LastName:        lastNamePtr,
		ProfileImageUrl: profileImagePtr,
		Bio:             bioPtr,
		Theme:           themePtr,
	}
	fmt.Printf("Executing DB query with params: %+v\n", params)

	user, err := r.db.CreateUser(ctx, params)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
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

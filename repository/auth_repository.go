package repository

import (
	"context"
	"fmt"
	"time"

	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
)

type AuthRepository interface {
	CreateAuth(ctx context.Context, params CreateAuthParams) error
	GetAuthByUserID(ctx context.Context, userID uuid.UUID) (*db.Auth, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash, salt string) error
	SetResetToken(ctx context.Context, userID uuid.UUID, resetToken string, expiresAt time.Time) error
	ClearResetToken(ctx context.Context, userID uuid.UUID) error
	VerifyEmail(ctx context.Context, userID uuid.UUID) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	IncrementFailedLoginAttempts(ctx context.Context, userID uuid.UUID) error
	SetAccountLockout(ctx context.Context, userID uuid.UUID, lockedUntil time.Time) error
	StoreRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) error
	InvalidateRefreshToken(ctx context.Context, userID uuid.UUID) error
}

type CreateAuthParams struct {
	UserID            uuid.UUID
	PasswordHash      string
	Salt              string
	IsEmailVerified   bool
	VerificationToken string
}

type SQLCAuthRepository struct {
	db *db.Queries
}

func NewAuthRepository(db *db.Queries) AuthRepository {
	return &SQLCAuthRepository{
		db: db,
	}
}

func (r *SQLCAuthRepository) CreateAuth(ctx context.Context, params CreateAuthParams) error {
	isEmailVerified := &params.IsEmailVerified
	var verificationToken *string
	if params.VerificationToken != "" {
		verificationToken = &params.VerificationToken
	}

	dbParams := db.CreateAuthParams{
		UserID:              params.UserID,
		PasswordHash:        params.PasswordHash,
		Salt:                params.Salt,
		IsEmailVerified:     isEmailVerified,
		VerificationToken:   verificationToken,
		ResetToken:          nil,
		ResetTokenExpiresAt: nil,
	}
	err := r.db.CreateAuth(ctx, dbParams)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok {
			if pgErr.Code == "23505" {
				return ErrDuplicateKey
			}
			if pgErr.Code == "23503" {
				return ErrForeignKeyViolation
			}
		}
		return ErrDatabase
	}

	return nil
}

func (r *SQLCAuthRepository) GetAuthByUserID(ctx context.Context, userID uuid.UUID) (*db.Auth, error) {
    fmt.Printf("Getting auth record for user ID: %s\n", userID)
    
    auth, err := r.db.GetAuthByUserID(ctx, userID)
    if err != nil {
        fmt.Printf("Error retrieving auth by user ID: %v\n", err)
        return nil, ErrRecordNotFound
    }
    
    fmt.Printf("Auth record found successfully\n")
    return auth, nil
}

func (r *SQLCAuthRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash, salt string) error {
	params := db.UpdatePasswordHashParams{
		UserID:       userID,
		PasswordHash: passwordHash,
		Salt:         salt,
	}
	err := r.db.UpdatePasswordHash(ctx, params)
	if err != nil {
		return ErrDatabase
	}
	return nil
}

func (r *SQLCAuthRepository) SetResetToken(ctx context.Context, userID uuid.UUID, resetToken string, expiresAt time.Time) error {
	params := db.SetResetTokenParams{
		UserID:              userID,
		ResetToken:          &resetToken,
		ResetTokenExpiresAt: &expiresAt,
	}
	err := r.db.SetResetToken(ctx, params)
	if err != nil {
		return ErrDatabase
	}
	return nil
}

func (r *SQLCAuthRepository) ClearResetToken(ctx context.Context, userID uuid.UUID) error {
	err := r.db.ClearResetToken(ctx, userID)
	if err != nil {
		return ErrDatabase
	}

	return nil
}

func (r *SQLCAuthRepository) VerifyEmail(ctx context.Context, userID uuid.UUID) error {
	err := r.db.VerifyEmail(ctx, userID)
	if err != nil {
		return ErrDatabase
	}
	return nil
}

func (r *SQLCAuthRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	err := r.db.UpdateLastLogin(ctx, userID)
	if err != nil {
		return ErrDatabase
	}
	return nil
}

func (r *SQLCAuthRepository) IncrementFailedLoginAttempts(ctx context.Context, userID uuid.UUID) error {
	err := r.db.IncrementFailedLoginAttempts(ctx, userID)
	if err != nil {
		return ErrDatabase
	}
	return nil
}

func (r *SQLCAuthRepository) SetAccountLockout(ctx context.Context, userID uuid.UUID, lockedUntil time.Time) error {
	params := db.SetAccountLockoutParams{
		UserID:      userID,
		LockedUntil: &lockedUntil,
	}
	err := r.db.SetAccountLockout(ctx, params)
	if err != nil {
		return ErrDatabase
	}
	return nil
}

func (r *SQLCAuthRepository) StoreRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) error {
    fmt.Printf("Storing refresh token for user %s, token length: %d\n", userID, len(refreshToken))
    
    params := db.StoreRefreshTokenParams{
        UserID:       userID,
        RefreshToken: &refreshToken,
    }
    
    err := r.db.StoreRefreshToken(ctx, params)
    if err != nil {
        fmt.Printf("Error storing refresh token: %v\n", err)
        return ErrDatabase
    }
    
    fmt.Println("Refresh token stored successfully")
    return nil
}

func (r *SQLCAuthRepository) InvalidateRefreshToken(ctx context.Context, userID uuid.UUID) error {
	err := r.db.InvalidateRefreshToken(ctx, userID)
	if err != nil {
		return ErrDatabase
	}
	return nil
}

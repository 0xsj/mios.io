package repository

import (
	"context"
	"fmt"
	"time"

	db "github.com/0xsj/mios.io/db/sqlc"
	"github.com/0xsj/mios.io/log"
	"github.com/0xsj/mios.io/pkg/errors"
	"github.com/google/uuid"
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
	GetAuthByVerificationToken(ctx context.Context, verificationToken string) (*db.Auth, error) // Add this if not present
	UpdateEmailVerificationStatus(ctx context.Context, userID uuid.UUID, isVerified bool) error // Add this
	ClearVerificationToken(ctx context.Context, userID uuid.UUID) error // Add this
}

type CreateAuthParams struct {
	UserID            uuid.UUID
	PasswordHash      string
	Salt              string
	IsEmailVerified   bool
	VerificationToken string
}

type SQLCAuthRepository struct {
	db     *db.Queries
	logger log.Logger
}

func NewAuthRepository(db *db.Queries, logger log.Logger) AuthRepository {
	return &SQLCAuthRepository{
		db:     db,
		logger: logger,
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

	start := time.Now()
	err := r.db.CreateAuth(ctx, dbParams)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "auth record")
		appErr.Log(r.logger)
		return appErr
	}

	r.logger.Infof("Auth record created successfully for user ID: %s in %v", params.UserID, duration)
	return nil
}

func (r *SQLCAuthRepository) GetAuthByUserID(ctx context.Context, userID uuid.UUID) (*db.Auth, error) {
	fmt.Printf("Getting auth record for user ID: %s\n", userID)

	start := time.Now()
	auth, err := r.db.GetAuthByUserID(ctx, userID)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "auth record")
		appErr.Log(r.logger)
		return nil, appErr
	}

	r.logger.Debugf("Auth record retrieved successfully for user ID: %s in %v", userID, duration)
	return auth, nil
}

func (r *SQLCAuthRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash, salt string) error {
	r.logger.Infof("Updating password for user ID: %s", userID)

	params := db.UpdatePasswordHashParams{
		UserID:       userID,
		PasswordHash: passwordHash,
		Salt:         salt,
	}

	start := time.Now()
	err := r.db.UpdatePasswordHash(ctx, params)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "password update")
		appErr.Log(r.logger)
		return appErr
	}

	r.logger.Infof("Password updated successfully for user ID: %s in %v", userID, duration)
	return nil
}

func (r *SQLCAuthRepository) SetResetToken(ctx context.Context, userID uuid.UUID, resetToken string, expiresAt time.Time) error {
	r.logger.Infof("Setting reset token for user ID: %s", userID)

	params := db.SetResetTokenParams{
		UserID:              userID,
		ResetToken:          &resetToken,
		ResetTokenExpiresAt: &expiresAt,
	}

	start := time.Now()
	err := r.db.SetResetToken(ctx, params)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "reset token")
		appErr.Log(r.logger)
		return appErr
	}

	r.logger.Infof("Reset token set successfully for user ID: %s in %v", userID, duration)
	return nil
}

func (r *SQLCAuthRepository) ClearResetToken(ctx context.Context, userID uuid.UUID) error {
	r.logger.Infof("Clearing reset token for user ID: %s", userID)

	start := time.Now()
	err := r.db.ClearResetToken(ctx, userID)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "reset token clearing")
		appErr.Log(r.logger)
		return appErr
	}

	r.logger.Infof("Reset token cleared successfully for user ID: %s in %v", userID, duration)
	return nil
}

func (r *SQLCAuthRepository) VerifyEmail(ctx context.Context, userID uuid.UUID) error {
	r.logger.Infof("Verifying email for user ID: %s", userID)

	start := time.Now()
	err := r.db.VerifyEmail(ctx, userID)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "email verification")
		appErr.Log(r.logger)
		return appErr
	}

	r.logger.Infof("Email verified successfully for user ID: %s in %v", userID, duration)
	return nil
}

func (r *SQLCAuthRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	r.logger.Debugf("Updating last login for user ID: %s", userID)

	start := time.Now()
	err := r.db.UpdateLastLogin(ctx, userID)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "last login update")
		appErr.Log(r.logger)
		return appErr
	}

	r.logger.Debugf("Last login updated successfully for user ID: %s in %v", userID, duration)
	return nil
}

func (r *SQLCAuthRepository) IncrementFailedLoginAttempts(ctx context.Context, userID uuid.UUID) error {
	r.logger.Warnf("Incrementing failed login attempts for user ID: %s", userID)

	start := time.Now()
	err := r.db.IncrementFailedLoginAttempts(ctx, userID)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "failed login attempts update")
		appErr.Log(r.logger)
		return appErr
	}

	r.logger.Warnf("Failed login attempts incremented for user ID: %s in %v", userID, duration)
	return nil
}

func (r *SQLCAuthRepository) SetAccountLockout(ctx context.Context, userID uuid.UUID, lockedUntil time.Time) error {
	r.logger.Warnf("Setting account lockout for user ID: %s until %v", userID, lockedUntil)

	params := db.SetAccountLockoutParams{
		UserID:      userID,
		LockedUntil: &lockedUntil,
	}

	start := time.Now()
	err := r.db.SetAccountLockout(ctx, params)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "account lockout")
		appErr.Log(r.logger)
		return appErr
	}

	r.logger.Warnf("Account lockout set successfully for user ID: %s in %v", userID, duration)
	return nil
}

func (r *SQLCAuthRepository) StoreRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	r.logger.Debugf("Storing refresh token for user ID: %s, token length: %d", userID, len(refreshToken))

	params := db.StoreRefreshTokenParams{
		UserID:       userID,
		RefreshToken: &refreshToken,
	}

	start := time.Now()
	err := r.db.StoreRefreshToken(ctx, params)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "refresh token storage")
		appErr.Log(r.logger)
		return appErr
	}

	r.logger.Debugf("Refresh token stored successfully for user ID: %s in %v", userID, duration)
	return nil
}

func (r *SQLCAuthRepository) InvalidateRefreshToken(ctx context.Context, userID uuid.UUID) error {
	r.logger.Infof("Invalidating refresh token for user ID: %s", userID)

	start := time.Now()
	err := r.db.InvalidateRefreshToken(ctx, userID)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "refresh token invalidation")
		appErr.Log(r.logger)
		return appErr
	}

	r.logger.Infof("Refresh token invalidated successfully for user ID: %s in %v", userID, duration)
	return nil
}

func (r *SQLCAuthRepository) GetAuthByVerificationToken(ctx context.Context, verificationToken string) (*db.Auth, error) {
	r.logger.Debugf("Getting auth by verification token")

	tokenPtr := verificationToken

	auth, err := r.db.GetAuthByVerificationToken(ctx, &tokenPtr)
	if err != nil {
		appErr := errors.HandleDBError(err, "auth record")
		appErr.Log(r.logger)
		return nil, appErr
	}

	r.logger.Debugf("Auth record retrieved successfully by verification token")
	return auth, nil
}

func (r *SQLCAuthRepository) UpdateEmailVerificationStatus(ctx context.Context, userID uuid.UUID, isVerified bool) error {
	r.logger.Infof("Updating email verification status for user ID: %s to %v", userID, isVerified)

	// This is the same as VerifyEmail, but more flexible
	if isVerified {
		return r.VerifyEmail(ctx, userID)
	}

	// If we need to unverify (rare case), we'd need a new SQLC query
	// For now, we'll just handle the verify case
	return errors.NewInternalError("Unverifying email is not supported", nil)
}

func (r *SQLCAuthRepository) ClearVerificationToken(ctx context.Context, userID uuid.UUID) error {
	r.logger.Infof("Clearing verification token for user ID: %s", userID)

	start := time.Now()
	err := r.db.ClearVerificationToken(ctx, userID)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "verification token clearing")
		appErr.Log(r.logger)
		return appErr
	}

	r.logger.Infof("Verification token cleared successfully for user ID: %s in %v", userID, duration)
	return nil
}
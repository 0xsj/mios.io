-- name: CreateAuth :exec
INSERT INTO auth (
    user_id, password_hash, salt, is_email_verified, verification_token, reset_token, reset_token_expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: GetAuthByUserID :one
SELECT * FROM auth
WHERE user_id = $1 LIMIT 1;

-- name: UpdatePasswordHash :exec
UPDATE auth
SET
    password_hash = $2,
    salt = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: SetResetToken :exec
UPDATE auth
SET
    reset_token = $2,
    reset_token_expires_at = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: ClearResetToken :exec
UPDATE auth
SET
    reset_token = NULL,
    reset_token_expires_at = NULL,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: VerifyEmail :exec
UPDATE auth
SET
    is_email_verified = true,
    verification_token = NULL,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: UpdateLastLogin :exec
UPDATE auth
SET
    last_login = CURRENT_TIMESTAMP,
    failed_login_attempts = 0,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: IncrementFailedLoginAttempts :exec
UPDATE auth
SET
    failed_login_attempts = failed_login_attempts + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: SetAccountLockout :exec
UPDATE auth
SET
    locked_until = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: StoreRefreshToken :exec
UPDATE auth
SET
    refresh_token = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: InvalidateRefreshToken :exec
UPDATE auth
SET
    refresh_token = NULL,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: GetAuthByVerificationToken :one
SELECT * FROM auth
WHERE verification_token = $1
LIMIT 1;

-- name: ClearVerificationToken :exec
UPDATE auth
SET verification_token = NULL
WHERE user_id = $1;
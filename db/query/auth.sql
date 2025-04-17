-- name: CreateAuth :one
INSERT INTO auth (
    user_id, password_hash, salt, is_email_verified, verification_token
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetAuthByUserID :one
SELECT * FROM auth
WHERE user_id = $1 LIMIT 1;

-- name: UpdatePassword :exec
UPDATE auth
SET
    password_hash = $2,
    salt = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: UpdateVerificationStatus :exec
UPDATE auth
SET
    is_email_verified = $2,
    verification_token = NULL,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: UpdateResetToken :exec
UPDATE auth
SET
    reset_token = $2,
    reset_token_expires_at = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: GetAuthByResetToken :one
SELECT * FROM auth
WHERE reset_token = $1 
AND reset_token_expires_at > CURRENT_TIMESTAMP
LIMIT 1;

-- name: UpdateLastLogin :exec
UPDATE auth
SET
    last_login = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: DeleteAuth :exec
DELETE FROM auth
WHERE user_id = $1;
-- name: CreateUser :one
INSERT INTO users (
    username, handle, email, first_name, last_name, 
    bio, profile_image_url, layout_version, custom_domain, 
    is_premium, is_admin, onboarded
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserByHandle :one
SELECT * FROM users
WHERE handle = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUser :exec
UPDATE users
SET 
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    bio = COALESCE($4, bio),
    profile_image_url = COALESCE($5, profile_image_url),
    layout_version = COALESCE($6, layout_version),
    custom_domain = COALESCE($7, custom_domain),
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: UpdateUsername :exec
UPDATE users
SET
    username = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: UpdateHandle :exec
UPDATE users
SET
    handle = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: UpdateEmail :exec
UPDATE users
SET
    email = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: UpdateUserPremiumStatus :exec
UPDATE users
SET
    is_premium = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: UpdateUserAdminStatus :exec
UPDATE users
SET
    is_admin = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: UpdateUserOnboardedStatus :exec
UPDATE users
SET
    onboarded = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;
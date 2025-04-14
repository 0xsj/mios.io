-- name: CreateUser :one
INSERT INTO users (
    username, email, first_name, last_name, 
    profile_image_url, bio, theme
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUser :exec
UPDATE users
SET 
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    profile_image_url = COALESCE($4, profile_image_url),
    bio = COALESCE($5, bio),
    theme = COALESCE($6, theme),
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: UpdateUserUsername :exec
UPDATE users
SET
    username = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: UpdateUserEmail :exec
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

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;
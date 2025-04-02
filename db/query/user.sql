-- name: CreateUser :one
INSERT INTO users (
  username,
  email,
  first_name,
  last_name
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE user_id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET 
  username = COALESCE($2, username),
  email = COALESCE($3, email),
  first_name = COALESCE($4, first_name),
  last_name = COALESCE($5, last_name),
  profile_image_url = COALESCE($6, profile_image_url),
  bio = COALESCE($7, bio),
  updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;
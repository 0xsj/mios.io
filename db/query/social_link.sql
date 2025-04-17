-- name: CreateSocialLink :one
INSERT INTO social_links (
    user_id, platform, username, url, position, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetSocialLink :one
SELECT * FROM social_links
WHERE social_link_id = $1 LIMIT 1;

-- name: GetSocialLinksByUserID :many
SELECT * FROM social_links
WHERE user_id = $1
ORDER BY position ASC;

-- name: GetActiveSocialLinksByUserID :many
SELECT * FROM social_links
WHERE user_id = $1 AND is_active = true
ORDER BY position ASC;

-- name: GetSocialLinkByPlatform :one
SELECT * FROM social_links
WHERE user_id = $1 AND platform = $2
LIMIT 1;

-- name: UpdateSocialLink :exec
UPDATE social_links
SET
    username = COALESCE($2, username),
    url = COALESCE($3, url),
    is_active = COALESCE($4, is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE social_link_id = $1;

-- name: UpdateSocialLinkPosition :exec
UPDATE social_links
SET
    position = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE social_link_id = $1;

-- name: DeleteSocialLink :exec
DELETE FROM social_links
WHERE social_link_id = $1;

-- name: DeleteUserSocialLinks :exec
DELETE FROM social_links
WHERE user_id = $1;
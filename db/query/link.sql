-- name: CreateLink :one
INSERT INTO links (
    user_id, title, url, description, icon, custom_thumbnail_url, position, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetLink :one
SELECT * FROM links
WHERE link_id = $1 LIMIT 1;

-- name: GetLinksByUserID :many
SELECT * FROM links
WHERE user_id = $1
ORDER BY position ASC;

-- name: GetActiveLinksByUserID :many
SELECT * FROM links
WHERE user_id = $1 AND is_active = true
ORDER BY position ASC;

-- name: UpdateLink :exec
UPDATE links
SET
    title = COALESCE($2, title),
    url = COALESCE($3, url),
    description = COALESCE($4, description),
    icon = COALESCE($5, icon),
    custom_thumbnail_url = COALESCE($6, custom_thumbnail_url),
    is_active = COALESCE($7, is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE link_id = $1;

-- name: UpdateLinkPosition :exec
UPDATE links
SET
    position = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE link_id = $1;

-- name: IncrementLinkClickCount :exec
UPDATE links
SET
    click_count = click_count + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE link_id = $1;

-- name: DeleteLink :exec
DELETE FROM links
WHERE link_id = $1;

-- name: DeleteUserLinks :exec
DELETE FROM links
WHERE user_id = $1;
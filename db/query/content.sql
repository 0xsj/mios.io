-- name: CreateContentItem :one
INSERT INTO content_items (
    user_id, content_id, content_type, title, href, url, media_type,
    desktop_x, desktop_y, desktop_style, mobile_x, mobile_y, mobile_style,
    halign, valign, content_data, overrides, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
) RETURNING *;

-- name: GetContentItem :one
SELECT * FROM content_items
WHERE item_id = $1 LIMIT 1;

-- name: GetUserContentItems :many
SELECT * FROM content_items
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateContentItem :exec
UPDATE content_items
SET
    title = COALESCE($2, title),
    href = COALESCE($3, href),
    url = COALESCE($4, url),
    media_type = COALESCE($5, media_type),
    desktop_style = COALESCE($6, desktop_style),
    mobile_style = COALESCE($7, mobile_style),
    halign = COALESCE($8, halign),
    valign = COALESCE($9, valign),
    content_data = COALESCE($10, content_data),
    overrides = COALESCE($11, overrides),
    is_active = COALESCE($12, is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE item_id = $1;

-- name: UpdateContentItemPosition :exec
UPDATE content_items
SET
    desktop_x = COALESCE($2, desktop_x),
    desktop_y = COALESCE($3, desktop_y),
    mobile_x = COALESCE($4, mobile_x),
    mobile_y = COALESCE($5, mobile_y),
    updated_at = CURRENT_TIMESTAMP
WHERE item_id = $1;

-- name: DeleteContentItem :exec
DELETE FROM content_items
WHERE item_id = $1;
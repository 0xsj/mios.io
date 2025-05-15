-- name: CreateLinkMetadata :one
INSERT INTO link_metadata (
    domain, url, title, description, favicon_url, image_url,
    platform_name, platform_type, platform_color, is_verified
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetLinkMetadataByURL :one
SELECT * FROM link_metadata
WHERE url = $1 LIMIT 1;

-- name: GetLinkMetadataByDomain :many
SELECT * FROM link_metadata
WHERE domain = $1
ORDER BY created_at DESC;

-- name: UpdateLinkMetadata :one
UPDATE link_metadata
SET
    title = COALESCE($2, title),
    description = COALESCE($3, description),
    favicon_url = COALESCE($4, favicon_url),
    image_url = COALESCE($5, image_url),
    platform_name = COALESCE($6, platform_name),
    platform_type = COALESCE($7, platform_type),
    platform_color = COALESCE($8, platform_color),
    is_verified = COALESCE($9, is_verified),
    updated_at = CURRENT_TIMESTAMP
WHERE url = $1
RETURNING *;

-- name: DeleteLinkMetadata :exec
DELETE FROM link_metadata
WHERE metadata_id = $1;
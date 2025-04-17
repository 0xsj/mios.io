-- name: CreateSection :one
INSERT INTO sections (
    user_id, title, description, position, is_active
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetSection :one
SELECT * FROM sections
WHERE section_id = $1 LIMIT 1;

-- name: GetSectionsByUserID :many
SELECT * FROM sections
WHERE user_id = $1
ORDER BY position ASC;

-- name: GetActiveSectionsByUserID :many
SELECT * FROM sections
WHERE user_id = $1 AND is_active = true
ORDER BY position ASC;

-- name: UpdateSection :exec
UPDATE sections
SET
    title = COALESCE($2, title),
    description = COALESCE($3, description),
    is_active = COALESCE($4, is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE section_id = $1;

-- name: UpdateSectionPosition :exec
UPDATE sections
SET
    position = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE section_id = $1;

-- name: DeleteSection :exec
DELETE FROM sections
WHERE section_id = $1;

-- name: DeleteUserSections :exec
DELETE FROM sections
WHERE user_id = $1;

-- name: AddLinkToSection :one
INSERT INTO section_links (
    section_id, link_id, position
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetSectionLinks :many
SELECT l.* FROM links l
JOIN section_links sl ON l.link_id = sl.link_id
WHERE sl.section_id = $1
ORDER BY sl.position ASC;

-- name: RemoveLinkFromSection :exec
DELETE FROM section_links
WHERE section_id = $1 AND link_id = $2;

-- name: UpdateSectionLinkPosition :exec
UPDATE section_links
SET
    position = $3
WHERE section_id = $1 AND link_id = $2;
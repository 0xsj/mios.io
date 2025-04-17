-- name: CreateTheme :one
INSERT INTO themes (
    name, background_color, text_color, button_style, font_family, is_premium, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetTheme :one
SELECT * FROM themes
WHERE theme_id = $1 LIMIT 1;

-- name: GetThemeByName :one
SELECT * FROM themes
WHERE name = $1 LIMIT 1;

-- name: ListThemes :many
SELECT * FROM themes
WHERE is_active = true
ORDER BY name ASC;

-- name: ListPremiumThemes :many
SELECT * FROM themes
WHERE is_premium = true AND is_active = true
ORDER BY name ASC;

-- name: ListFreeThemes :many
SELECT * FROM themes
WHERE is_premium = false AND is_active = true
ORDER BY name ASC;

-- name: UpdateTheme :exec
UPDATE themes
SET
    background_color = COALESCE($2, background_color),
    text_color = COALESCE($3, text_color),
    button_style = COALESCE($4, button_style),
    font_family = COALESCE($5, font_family),
    is_premium = COALESCE($6, is_premium),
    is_active = COALESCE($7, is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE theme_id = $1;

-- name: DeleteTheme :exec
DELETE FROM themes
WHERE theme_id = $1;

-- name: CreateUserTheme :one
INSERT INTO user_themes (
    user_id, name, background_color, background_image_url, text_color, 
    button_style, font_family, custom_css, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetUserTheme :one
SELECT * FROM user_themes
WHERE user_theme_id = $1 LIMIT 1;

-- name: GetUserThemeByName :one
SELECT * FROM user_themes
WHERE user_id = $1 AND name = $2 LIMIT 1;

-- name: GetUserThemes :many
SELECT * FROM user_themes
WHERE user_id = $1
ORDER BY name ASC;

-- name: UpdateUserTheme :exec
UPDATE user_themes
SET
    name = COALESCE($2, name),
    background_color = COALESCE($3, background_color),
    background_image_url = COALESCE($4, background_image_url),
    text_color = COALESCE($5, text_color),
    button_style = COALESCE($6, button_style),
    font_family = COALESCE($7, font_family),
    custom_css = COALESCE($8, custom_css),
    is_active = COALESCE($9, is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE user_theme_id = $1;

-- name: DeleteUserTheme :exec
DELETE FROM user_themes
WHERE user_theme_id = $1;
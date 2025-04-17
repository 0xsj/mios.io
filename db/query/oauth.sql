-- name: CreateOAuthAccount :one
INSERT INTO oauth_accounts (
    user_id, provider, provider_user_id, email, name,
    access_token, refresh_token, token_expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetOAuthAccountByProviderID :one
SELECT * FROM oauth_accounts
WHERE provider = $1 AND provider_user_id = $2
LIMIT 1;

-- name: GetOAuthAccountsByUserID :many
SELECT * FROM oauth_accounts
WHERE user_id = $1;

-- name: UpdateOAuthTokens :exec
UPDATE oauth_accounts
SET 
    access_token = $3,
    refresh_token = $4,
    token_expires_at = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE provider = $1 AND provider_user_id = $2;

-- name: DeleteOAuthAccount :exec
DELETE FROM oauth_accounts
WHERE provider = $1 AND provider_user_id = $2;

-- name: DeleteAllUserOAuthAccounts :exec
DELETE FROM oauth_accounts
WHERE user_id = $1;
-- Drop indexes
DROP INDEX IF EXISTS idx_oauth_accounts_user_id;
DROP INDEX IF EXISTS idx_oauth_provider_id;

-- Drop table
DROP TABLE IF EXISTS oauth_accounts;
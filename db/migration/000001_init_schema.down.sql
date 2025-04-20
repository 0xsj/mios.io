-- Drop indexes
DROP INDEX IF EXISTS idx_oauth_provider_id;
DROP INDEX IF EXISTS idx_oauth_accounts_user_id;
DROP INDEX IF EXISTS idx_user_themes_user_id;
DROP INDEX IF EXISTS idx_auth_user_id;
DROP INDEX IF EXISTS idx_analytics_item_id;
DROP INDEX IF EXISTS idx_analytics_user_id;
DROP INDEX IF EXISTS idx_content_items_user_id;

-- Drop tables (in reverse order of creation to handle dependencies)
DROP TABLE IF EXISTS oauth_accounts;
DROP TABLE IF EXISTS user_themes;
DROP TABLE IF EXISTS themes;
DROP TABLE IF EXISTS analytics;
DROP TABLE IF EXISTS content_items;
DROP TABLE IF EXISTS auth;
DROP TABLE IF EXISTS users;

-- Drop extensions
DROP EXTENSION IF EXISTS "uuid-ossp";
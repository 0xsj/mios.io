-- Drop indexes
DROP INDEX IF EXISTS idx_user_themes_user_id;
DROP INDEX IF EXISTS idx_analytics_social_link_id;
DROP INDEX IF EXISTS idx_analytics_link_id;
DROP INDEX IF EXISTS idx_analytics_user_id;
DROP INDEX IF EXISTS idx_sections_user_id;
DROP INDEX IF EXISTS idx_social_links_user_id;
DROP INDEX IF EXISTS idx_links_user_id;
DROP INDEX IF EXISTS idx_auth_user_id;

-- Drop tables (in reverse order of creation to handle dependencies)
DROP TABLE IF EXISTS user_themes;
DROP TABLE IF EXISTS themes;
DROP TABLE IF EXISTS analytics;
DROP TABLE IF EXISTS section_links;
DROP TABLE IF EXISTS sections;
DROP TABLE IF EXISTS social_links;
DROP TABLE IF EXISTS links;
DROP TABLE IF EXISTS auth;
DROP TABLE IF EXISTS users;

DROP extension IF EXISTS "uuid-ossp";
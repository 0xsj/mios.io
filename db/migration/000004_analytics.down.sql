-- db/migration/000004_add_analytics_page_view.down.sql
DROP INDEX IF EXISTS idx_analytics_page_view;
ALTER TABLE analytics
DROP COLUMN page_view;
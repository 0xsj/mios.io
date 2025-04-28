-- db/migration/000004_add_analytics_page_view.up.sql
ALTER TABLE analytics
ADD COLUMN page_view BOOLEAN DEFAULT false;

-- Add an index for faster querying
CREATE INDEX idx_analytics_page_view ON analytics(user_id, page_view);
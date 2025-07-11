-- Remove indexes
DROP INDEX IF EXISTS idx_embed_configs_platform;
DROP INDEX IF EXISTS idx_conversions_analytics_id;
DROP INDEX IF EXISTS idx_analytics_device;
DROP INDEX IF EXISTS idx_analytics_country;
DROP INDEX IF EXISTS idx_users_theme_id;

-- Drop new tables
DROP TABLE IF EXISTS embed_configs;
DROP TABLE IF EXISTS conversions;

-- Remove content item columns
ALTER TABLE content_items DROP COLUMN IF EXISTS auto_embed;
ALTER TABLE content_items DROP COLUMN IF EXISTS embed_data;
ALTER TABLE content_items DROP COLUMN IF EXISTS custom_styling;

-- Remove analytics columns
ALTER TABLE analytics DROP COLUMN IF EXISTS utm_campaign;
ALTER TABLE analytics DROP COLUMN IF EXISTS utm_medium;
ALTER TABLE analytics DROP COLUMN IF EXISTS utm_source;
ALTER TABLE analytics DROP COLUMN IF EXISTS browser;
ALTER TABLE analytics DROP COLUMN IF EXISTS device_type;
ALTER TABLE analytics DROP COLUMN IF EXISTS country;

-- Remove users columns
ALTER TABLE users DROP COLUMN IF EXISTS theme_customization;
ALTER TABLE users DROP COLUMN IF EXISTS theme_id;

-- Restore original themes table
DROP TABLE IF EXISTS themes CASCADE;
CREATE TABLE themes (
    theme_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    background_color VARCHAR(20),
    text_color VARCHAR(20),
    button_style VARCHAR(50),
    font_family VARCHAR(100),
    is_premium BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
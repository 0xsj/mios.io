-- Enhanced Theme System
DROP TABLE IF EXISTS themes CASCADE;
CREATE TABLE themes (
    theme_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    preview_image_url VARCHAR(500),
    is_premium BOOLEAN DEFAULT FALSE,
    config JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add theme references to users table
ALTER TABLE users ADD COLUMN theme_id UUID REFERENCES themes(theme_id);
ALTER TABLE users ADD COLUMN theme_customization JSONB;

-- Add custom styling to content items
ALTER TABLE content_items ADD COLUMN custom_styling JSONB;

-- Enhanced Analytics
ALTER TABLE analytics ADD COLUMN country VARCHAR(5);
ALTER TABLE analytics ADD COLUMN device_type VARCHAR(20);
ALTER TABLE analytics ADD COLUMN browser VARCHAR(50);
ALTER TABLE analytics ADD COLUMN utm_source VARCHAR(100);
ALTER TABLE analytics ADD COLUMN utm_medium VARCHAR(100);
ALTER TABLE analytics ADD COLUMN utm_campaign VARCHAR(100);

-- Conversion tracking
CREATE TABLE conversions (
    conversion_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    analytics_id UUID REFERENCES analytics(analytics_id),
    conversion_type VARCHAR(50) NOT NULL,
    conversion_value DECIMAL(10,2),
    conversion_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Enhanced Content and Embeds
ALTER TABLE content_items ADD COLUMN embed_data JSONB;
ALTER TABLE content_items ADD COLUMN auto_embed BOOLEAN DEFAULT FALSE;

-- Platform-specific embed configurations
CREATE TABLE embed_configs (
    config_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    platform VARCHAR(50) NOT NULL,
    embed_template TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_users_theme_id ON users(theme_id);
CREATE INDEX idx_analytics_country ON analytics(country) WHERE country IS NOT NULL;
CREATE INDEX idx_analytics_device ON analytics(device_type) WHERE device_type IS NOT NULL;
CREATE INDEX idx_conversions_analytics_id ON conversions(analytics_id);
CREATE INDEX idx_embed_configs_platform ON embed_configs(platform);
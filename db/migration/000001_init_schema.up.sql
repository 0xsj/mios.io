CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Core user information
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    handle VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    bio TEXT,
    profile_image_url TEXT,
    layout_version VARCHAR(20) DEFAULT 'v1',
    custom_domain VARCHAR(255) UNIQUE,
    is_premium BOOLEAN DEFAULT FALSE,
    is_admin BOOLEAN DEFAULT FALSE,
    onboarded BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Authentication information
CREATE TABLE auth (
    auth_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    password_hash VARCHAR(255) NOT NULL,
    salt VARCHAR(255) NOT NULL,
    is_email_verified BOOLEAN DEFAULT FALSE,
    verification_token VARCHAR(255),
    reset_token VARCHAR(255),
    reset_token_expires_at TIMESTAMP WITH TIME ZONE,
    last_login TIMESTAMP WITH TIME ZONE,
    UNIQUE(user_id)
);

-- Content items (grid layout)
CREATE TABLE content_items (
    item_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    content_id VARCHAR(100) NOT NULL, -- External reference ID
    content_type VARCHAR(50) NOT NULL, -- 'link', 'media', 'rich-text', etc.
    
    -- Link data
    title VARCHAR(100),
    href TEXT,
    
    -- Media data
    url TEXT,
    media_type VARCHAR(50),
    
    -- Layout information
    desktop_x INTEGER,
    desktop_y INTEGER,
    desktop_style VARCHAR(20), -- e.g., "4x2"
    mobile_x INTEGER,
    mobile_y INTEGER, 
    mobile_style VARCHAR(20),
    
    -- Styling and alignment
    halign VARCHAR(20), -- 'left', 'center', 'right'
    valign VARCHAR(20), -- 'top', 'middle', 'bottom'
    
    -- Additional data
    content_data JSONB, -- For storing type-specific data
    overrides JSONB, -- For customizations
    
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Analytics for tracking clicks and views
CREATE TABLE analytics (
    analytics_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    item_id UUID NOT NULL REFERENCES content_items(item_id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    ip_address VARCHAR(45),
    user_agent TEXT,
    referrer TEXT,
    clicked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Themes for styling profiles
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

-- User-specific theme customizations
CREATE TABLE user_themes (
    user_theme_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    background_color VARCHAR(20),
    background_image_url TEXT,
    text_color VARCHAR(20),
    button_style VARCHAR(50),
    font_family VARCHAR(100),
    custom_css TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, name)
);

-- Social connections (OAuth)
CREATE TABLE oauth_accounts (
    oauth_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL, -- 'google', 'twitter', etc.
    provider_user_id VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    name VARCHAR(255),
    access_token TEXT,
    refresh_token TEXT,
    token_expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(provider, provider_user_id)
);

-- Create indexes for performance
CREATE INDEX idx_content_items_user_id ON content_items(user_id);
CREATE INDEX idx_analytics_user_id ON analytics(user_id);
CREATE INDEX idx_analytics_item_id ON analytics(item_id);
CREATE INDEX idx_auth_user_id ON auth(user_id);
CREATE INDEX idx_user_themes_user_id ON user_themes(user_id);
CREATE INDEX idx_oauth_accounts_user_id ON oauth_accounts(user_id);
CREATE INDEX idx_oauth_provider_id ON oauth_accounts(provider, provider_user_id);
-- Create OAuth accounts table
CREATE TABLE IF NOT EXISTS oauth_accounts (
    oauth_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL, -- e.g., 'google', 'twitter', 'facebook'
    provider_user_id VARCHAR(255) NOT NULL, -- ID from the provider
    email VARCHAR(255),
    name VARCHAR(255),
    access_token TEXT,
    refresh_token TEXT,
    token_expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(provider, provider_user_id)
);

-- Create index on user_id for efficient lookups
CREATE INDEX idx_oauth_accounts_user_id ON oauth_accounts(user_id);

-- Create composite index on provider and provider_user_id
CREATE INDEX idx_oauth_provider_id ON oauth_accounts(provider, provider_user_id);
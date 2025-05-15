CREATE TABLE link_metadata (
    metadata_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    domain TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    title TEXT,
    description TEXT,
    favicon_url TEXT,
    image_url TEXT,
    platform_name TEXT,
    platform_type TEXT,
    platform_color TEXT,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_link_metadata_domain ON link_metadata(domain);
CREATE INDEX idx_link_metadata_url ON link_metadata(url);
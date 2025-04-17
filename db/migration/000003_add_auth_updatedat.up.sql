-- Add updated_at column to auth table
ALTER TABLE auth ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
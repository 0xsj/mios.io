-- Down migration
ALTER TABLE auth
ALTER COLUMN refresh_token TYPE VARCHAR(255);
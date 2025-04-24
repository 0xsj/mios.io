ALTER TABLE auth
DROP COLUMN refresh_token,
DROP COLUMN failed_login_attempts,
DROP COLUMN locked_until,
DROP COLUMN created_at,
DROP COLUMN updated_at;
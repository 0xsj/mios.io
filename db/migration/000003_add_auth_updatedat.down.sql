-- Remove updated_at column from auth table
ALTER TABLE auth DROP COLUMN IF EXISTS updated_at;
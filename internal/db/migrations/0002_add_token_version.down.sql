-- Remove token_version column from user_account table
DROP INDEX IF EXISTS idx_user_account_token_version;
ALTER TABLE user_account DROP COLUMN token_version;

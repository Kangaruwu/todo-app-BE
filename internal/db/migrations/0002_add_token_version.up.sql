-- Add token_version column to user_account table
ALTER TABLE user_account ADD COLUMN token_version INTEGER DEFAULT 1;

-- Create index for performance
CREATE INDEX idx_user_account_token_version ON user_account(user_id, token_version);

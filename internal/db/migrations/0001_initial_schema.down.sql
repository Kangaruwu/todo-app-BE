BEGIN;


DROP INDEX IF EXISTS idx_users_id;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_todos_created_at;
DROP INDEX IF EXISTS idx_todos_completed;
DROP INDEX IF EXISTS idx_todos_user_id;

DROP TABLE IF EXISTS todos;

DROP TABLE IF EXISTS user_account_external;

DROP TABLE IF EXISTS external_providers;

DROP TABLE IF EXISTS user_account;

DROP TYPE IF EXISTS email_validation_status_enum;

DROP TYPE IF EXISTS user_role_enum;


COMMIT;
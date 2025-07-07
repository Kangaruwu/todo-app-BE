-- down.sql
BEGIN;

DROP TABLE IF EXISTS todos;

DROP TABLE IF EXISTS user_login_data_external;

DROP TABLE IF EXISTS external_providers;

DROP TABLE IF EXISTS user_login_data;

DROP TABLE IF EXISTS user_account;

DROP TYPE IF EXISTS user_role_enum;

DROP EXTENSION IF EXISTS "uuid-ossp";

COMMIT;

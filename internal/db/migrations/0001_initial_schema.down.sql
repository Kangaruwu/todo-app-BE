-- Drop indexes 
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_todos_created_at;
DROP INDEX IF EXISTS idx_todos_completed;
DROP INDEX IF EXISTS idx_todos_user_id;

-- Drop tables (todos first since having FK from users)
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS users;

-- Drop extension 
DROP EXTENSION IF EXISTS "uuid-ossp";



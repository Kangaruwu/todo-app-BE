-- 0.1 User role type
CREATE TYPE user_role_enum AS ENUM ('user', 'admin');

-- 1. user_account
CREATE TABLE user_account (
    user_id              INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_role             user_role_enum          NOT NULL DEFAULT 'user'
);

-- 2. user_login_data (1-1 với user_account qua user_id)
CREATE TABLE user_login_data (
    user_id                   			INTEGER               	PRIMARY KEY,
    user_name                 			VARCHAR(20)           	NOT NULL,
    password_hash             			VARCHAR(250)          	NOT NULL,
    password_salt             			VARCHAR(100)          	NOT NULL,
    hash_algorithm            			VARCHAR(10)           	NOT NULL,
    email_address             			VARCHAR(100)          	NOT NULL,

    confirmation_token        			VARCHAR(100),
    confirmation_token_generation_time  TIMESTAMP,
    email_validation_status   			VARCHAR(20)	      	    NOT NULL,

    password_recovery_token   			VARCHAR(100),
    recovery_token_time       			TIMESTAMP,

    FOREIGN KEY (user_id)     			REFERENCES user_account	(user_id)
);

-- 3. external_providers
CREATE TABLE external_providers (
    external_provider_id     INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    provider_name            VARCHAR(50)           NOT NULL,
    ws_endpoint              VARCHAR(200)          NOT NULL
);

-- 4. user_login_data_external (liên kết 1-n giữa user_login_data và external_providers)
CREATE TABLE user_login_data_external (
    user_id                   INTEGER               NOT NULL,
    external_provider_id      INTEGER               NOT NULL,
    external_provider_token   VARCHAR(100)          NOT NULL,
    PRIMARY KEY (user_id, external_provider_id),
    FOREIGN KEY (user_id)               REFERENCES user_login_data   (user_id),
    FOREIGN KEY (external_provider_id)  REFERENCES external_providers (external_provider_id)
);

-- 5. Create todos table
CREATE TABLE todos (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    deadline TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    user_id INTEGER NOT NULL REFERENCES user_account(user_id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX idx_todos_user_id ON todos(user_id);
CREATE INDEX idx_todos_completed ON todos(completed);
CREATE INDEX idx_todos_created_at ON todos(created_at);
CREATE INDEX idx_users_email ON user_login_data(email_address);
CREATE INDEX idx_users_username ON user_login_data(user_name);


-- Create sample data for users
INSERT INTO user_account (user_role) VALUES ('admin');
INSERT INTO user_login_data (user_id, user_name, password_hash, password_salt, hash_algorithm, email_address, email_validation_status)
VALUES (1, 'admin', 'hashed_password', 'salt_value', 'bcrypt', 'admin@example.com', 'confirmed');

-- Create sample data for external providers
INSERT INTO external_providers (provider_name, ws_endpoint)
VALUES ('Google', 'https://accounts.google.com/o/oauth2/auth'),
       ('GitHub', 'https://github.com/login/oauth/authorize');
-- Create sample data for user_login_data_external
INSERT INTO user_login_data_external (user_id, external_provider_id, external_provider_token)
VALUES (1, 1, 'google_token'),
       (1, 2, 'github_token');
-- Create sample data for todos
INSERT INTO todos (title, deadline, completed, user_id)
VALUES
('Buy groceries', NOW() + INTERVAL '1 day', FALSE, 1),
('Complete project report', NOW() + INTERVAL '2 days', FALSE, 1),
('Call mom', NOW() + INTERVAL '3 days', TRUE, 1);
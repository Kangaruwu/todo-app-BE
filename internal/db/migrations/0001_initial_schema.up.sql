-- 0.1 uuid extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 0.2 User role type
CREATE TYPE user_role_enum AS ENUM ('user', 'admin');

-- 0.3 Email validation status type
CREATE TYPE email_validation_status_enum AS ENUM ('unconfirmed', 'confirmed', 'pending');

-- 1. user_account
CREATE TABLE
    user_account (
        user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        user_name VARCHAR(20) NOT NULL UNIQUE,
        user_role user_role_enum NOT NULL DEFAULT 'user',
        password_hash VARCHAR(250) NOT NULL,
        email_address VARCHAR(100) NOT NULL,
        verification_token VARCHAR(250),
        verification_token_generation_time TIMESTAMP,
        email_validation_status email_validation_status_enum NOT NULL DEFAULT 'unconfirmed',
        password_recovery_token VARCHAR(250),
        password_recovery_token_generation_time TIMESTAMP,
        created_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT NOW (),
        updated_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT NOW ()
    );

-- 2. external_providers
CREATE TABLE
    external_providers (
        external_provider_id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        provider_name VARCHAR(50) NOT NULL,
        ws_endpoint VARCHAR(200) NOT NULL
    );

-- 3. user_account_external (liên kết 1-n giữa user_account và external_providers)
CREATE TABLE
    user_account_external (
        user_id UUID NOT NULL,
        external_provider_id UUID NOT NULL,
        external_provider_token VARCHAR(100) NOT NULL,
        PRIMARY KEY (user_id, external_provider_id),
        FOREIGN KEY (user_id) REFERENCES user_account (user_id),
        FOREIGN KEY (external_provider_id) REFERENCES external_providers (external_provider_id)
    );

-- 4. Create todos table
CREATE TABLE
    todos (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        title VARCHAR(255) NOT NULL,
        deadline TIMESTAMP
        WITH
            TIME ZONE DEFAULT NOW (),
            completed BOOLEAN DEFAULT FALSE,
            created_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT NOW (),
            updated_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT NOW (),
            user_id UUID NOT NULL REFERENCES user_account (user_id) ON DELETE CASCADE
    );

-- Create indexes for better performance
CREATE INDEX idx_todos_user_id ON todos (user_id);

CREATE INDEX idx_todos_completed ON todos (completed);

CREATE INDEX idx_todos_created_at ON todos (created_at);

CREATE INDEX idx_users_email ON user_account (email_address);

CREATE INDEX idx_users_username ON user_account (user_name);

CREATE INDEX idx_users_id ON user_account (user_id);

-- Create sample data for users
INSERT INTO
    user_account (
        user_name,
        user_role,
        password_hash,
        email_address,
        email_validation_status
    )
VALUES
    (
        'BaoNguyxn',
        'admin',
        'hashed_password',
        'admin@example.com',
        'confirmed'
    );

-- Create sample data for external providers
INSERT INTO
    external_providers (provider_name, ws_endpoint)
VALUES
    (
        'Google',
        'https://accounts.google.com/o/oauth2/auth'
    ),
    (
        'GitHub',
        'https://github.com/login/oauth/authorize'
    );

-- Create sample data for user_account_external
INSERT INTO
    user_account_external (
        user_id,
        external_provider_id,
        external_provider_token
    )
VALUES
    (
        (
            SELECT
                user_id
            FROM
                user_account
            WHERE
                user_role = 'admin'
        ),
        (
            SELECT
                external_provider_id
            FROM
                external_providers
            WHERE
                provider_name = 'Google'
        ),
        'google_token_12345'
    ),
    (
        (
            SELECT
                user_id
            FROM
                user_account
            WHERE
                user_role = 'admin'
        ),
        (
            SELECT
                external_provider_id
            FROM
                external_providers
            WHERE
                provider_name = 'GitHub'
        ),
        'github_token_67890'
    );

-- Create sample data for todos
INSERT INTO
    todos (title, deadline, completed, user_id)
VALUES
    (
        'Buy groceries',
        NOW () + INTERVAL '1 day',
        FALSE,
        (
            SELECT
                user_id
            FROM
                user_account
            WHERE
                user_role = 'admin'
        )
    ),
    (
        'Complete project report',
        NOW () + INTERVAL '2 days',
        FALSE,
        (
            SELECT
                user_id
            FROM
                user_account
            WHERE
                user_role = 'admin'
        )
    ),
    (
        'Schedule dentist appointment',
        NOW () + INTERVAL '5 days',
        TRUE,
        (
            SELECT
                user_id
            FROM
                user_account
            WHERE
                user_role = 'admin'
        )
    ),
    (
        'Call mom',
        NOW () + INTERVAL '3 days',
        TRUE,
        (
            SELECT
                user_id
            FROM
                user_account
            WHERE
                user_role = 'admin'
        )
    );
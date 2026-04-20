CREATE EXTENSION
IF NOT EXISTS "pgcrypto";

CREATE TYPE users_role_enum as ENUM
('admin', 'user');

CREATE TABLE users
(
    internal_id     BIGSERIAL       PRIMARY KEY,
    public_id       UUID            NOT NULL DEFAULT gen_random_uuid(),
    username        VARCHAR(255)    NOT NULL,
    email           VARCHAR(255)    NOT NULL UNIQUE,
    password        text            NOT NULL,
    role            users_role_enum NOT NULL DEFAULT 'user',
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ     NULL,
    CONSTRAINT user_public_id_unique UNIQUE(public_id)
);
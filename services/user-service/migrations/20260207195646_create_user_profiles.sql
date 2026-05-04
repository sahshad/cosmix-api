-- +goose Up
-- +goose StatementBegin

CREATE TABLE user_profiles (
    id BIGSERIAL PRIMARY KEY,
    auth_user_id BIGINT NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE,
    display_name VARCHAR(255),
    bio TEXT,
    avatar_url TEXT,
    is_private BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    date_of_birth DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS user_profiles;

-- +goose StatementEnd
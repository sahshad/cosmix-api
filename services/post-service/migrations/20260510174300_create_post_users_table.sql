-- +goose Up
-- +goose StatementBegin

CREATE TABLE post_users (
    user_id BIGINT PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    display_name VARCHAR(150) NOT NULL,
    avatar_url TEXT NULL,
    is_verified BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_post_users_username ON post_users(username);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS post_users;

-- +goose StatementEnd
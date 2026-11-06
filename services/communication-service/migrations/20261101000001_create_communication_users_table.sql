-- +goose Up
-- +goose StatementBegin

CREATE TABLE communication_users (
    user_id UUID PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    display_name VARCHAR(150) NOT NULL,
    avatar_url TEXT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_communication_users_username ON communication_users(username);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS communication_users;

-- +goose StatementEnd

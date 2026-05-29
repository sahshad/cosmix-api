-- +goose Up
-- +goose StatementBegin

CREATE TABLE auth_users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    display_name VARCHAR(100) NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    last_login_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL,

    CONSTRAINT uq_auth_users_email UNIQUE(email)
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE auth_users;

-- +goose StatementEnd

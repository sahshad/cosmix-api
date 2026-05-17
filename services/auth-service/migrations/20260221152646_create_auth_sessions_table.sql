-- +goose Up
-- +goose StatementBegin

CREATE TABLE auth_sessions (
    id BIGSERIAL PRIMARY KEY,
    auth_user_id BIGINT NOT NULL,
    refresh_token_hash VARCHAR(500) NOT NULL,
    device VARCHAR(255) NULL,
    ip_address VARCHAR(100) NULL,
    user_agent VARCHAR(500) NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL,

    CONSTRAINT fk_auth_sessions_auth_user_id 
    FOREIGN KEY (auth_user_id)
    REFERENCES auth_users(id)
    ON DELETE CASCADE
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE auth_sessions;

-- +goose StatementEnd
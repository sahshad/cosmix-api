-- +goose Up
-- +goose StatementBegin

CREATE TABLE auth_sessions (
    id BIGSERIAL PRIMARY KEY,
    auth_user_id INT NOT NULL,
    refresh_token_hash VARCHAR(500) NOT NULL,
    device VARCHAR(255) NULL,
    ip_address VARCHAR(100) NULL,
    user_agent VARCHAR(500) NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Foreign Key to auth_users
ALTER TABLE auth_sessions
ADD CONSTRAINT FK_auth_sessions_user
FOREIGN KEY (auth_user_id)
REFERENCES auth_users(id)
ON DELETE CASCADE;

-- Index for faster lookup
CREATE INDEX IX_auth_sessions_user_id
ON auth_sessions(auth_user_id);

CREATE INDEX IX_auth_sessions_token_hash
ON auth_sessions(refresh_token_hash);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE auth_sessions;

-- +goose StatementEnd
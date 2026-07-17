-- +goose Up
-- +goose StatementBegin

CREATE TABLE auth_tokens (
    id UUID PRIMARY KEY,
    auth_user_id UUID NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_auth_tokens_type
    CHECK (
        type IN (
            'email_verification',
            'password_reset'
        )
    ),

    CONSTRAINT fk_auth_tokens_auth_user_id
    FOREIGN KEY (auth_user_id)
    REFERENCES auth_users(id)
    ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE auth_tokens;

-- +goose StatementEnd
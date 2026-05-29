-- +goose Up
-- +goose StatementBegin
CREATE TABLE email_verification_tokens (
    id BIGSERIAL PRIMARY KEY,
    auth_user_id BIGINT NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_email_verification_tokens_auth_user_id FOREIGN KEY (auth_user_id) REFERENCES auth_users(id) ON DELETE CASCADE
);

CREATE INDEX idx_email_verification_tokens_auth_user_id ON email_verification_tokens(auth_user_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE email_verification_tokens;

-- +goose StatementEnd
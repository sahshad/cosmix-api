-- +goose Up
-- +goose StatementBegin
CREATE TABLE email_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NULL,
    recipient VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    subject TEXT NOT NULL,
    template VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    provider VARCHAR(100) NULL,
    error_message TEXT NULL,
    sent_at TIMESTAMPTZ NULL,
    failed_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_email_logs_user_id ON email_logs(user_id);

CREATE INDEX idx_email_logs_status ON email_logs(status);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS email_logs;

-- +goose StatementEnd
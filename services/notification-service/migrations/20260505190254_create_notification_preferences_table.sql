-- +goose Up
-- +goose StatementBegin
CREATE TABLE notification_preferences (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    email_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    push_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    internal_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    like_notifications BOOLEAN NOT NULL DEFAULT TRUE,
    comment_notifications BOOLEAN NOT NULL DEFAULT TRUE,
    follow_notifications BOOLEAN NOT NULL DEFAULT TRUE,
    mention_notifications BOOLEAN NOT NULL DEFAULT TRUE,
    message_notifications BOOLEAN NOT NULL DEFAULT TRUE,
    marketing_emails_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_notification_preferences_user_id ON notification_preferences(user_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notification_preferences;

-- +goose StatementEnd
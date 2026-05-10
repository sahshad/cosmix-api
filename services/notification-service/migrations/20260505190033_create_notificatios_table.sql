-- +goose Up
-- +goose StatementBegin
CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    actor_id BIGINT NULL,
    actor_username VARCHAR(255) NULL,
    actor_display_name VARCHAR(255) NULL,
    actor_avatar_url TEXT NULL,
    type VARCHAR(100) NOT NULL,
    entity_id BIGINT NULL,
    entity_type VARCHAR(100) NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    image_url TEXT NULL,
    action_url TEXT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    read_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_read ON notifications(user_id, is_read);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notifications;

-- +goose StatementEnd
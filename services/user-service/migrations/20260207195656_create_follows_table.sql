-- +goose Up
-- +goose StatementBegin

CREATE TABLE follows (
    id UUID PRIMARY KEY,
    follower_id UUID NOT NULL,
    following_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_follows_follower_following UNIQUE (follower_id, following_id),

    CONSTRAINT fk_follows_follower_id
    FOREIGN KEY (follower_id)
    REFERENCES users(user_id)
    ON DELETE CASCADE,

    CONSTRAINT fk_follows_following_id
    FOREIGN KEY (following_id)
    REFERENCES users(user_id)
    ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE follows;

-- +goose StatementEnd

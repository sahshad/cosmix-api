-- +goose Up
-- +goose StatementBegin

CREATE TABLE likes (
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT fk_likes_post_id
        FOREIGN KEY (post_id)
        REFERENCES posts(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_likes_user_id
        FOREIGN KEY (user_id)
        REFERENCES post_users(user_id),

    CONSTRAINT uq_likes_post_user
        UNIQUE(post_id, user_id)
);

CREATE INDEX idx_likes_post_id ON likes(post_id);
CREATE INDEX idx_likes_user_id ON likes(user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS likes;

-- +goose StatementEnd
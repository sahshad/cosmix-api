-- +goose Up
-- +goose StatementBegin

CREATE TABLE comment_likes (
    id UUID PRIMARY KEY,
    comment_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT fk_comment_likes_comment_id
        FOREIGN KEY (comment_id)
        REFERENCES comments(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_comment_likes_user_id
        FOREIGN KEY (user_id)
        REFERENCES post_users(user_id),

    CONSTRAINT uq_comment_likes_comment_user
        UNIQUE(comment_id, user_id)
);

CREATE INDEX idx_comment_likes_comment_id ON comment_likes(comment_id);
CREATE INDEX idx_comment_likes_user_id ON comment_likes(user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS comment_likes;

-- +goose StatementEnd

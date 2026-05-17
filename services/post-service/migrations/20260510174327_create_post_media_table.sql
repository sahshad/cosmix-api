-- +goose Up
-- +goose StatementBegin

CREATE TABLE post_media (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL,
    public_id VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    duration INTEGER NULL,
    created_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT fk_post_media_post_id
        FOREIGN KEY (post_id)
        REFERENCES posts(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_post_media_post_id ON post_media(post_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS post_media;

-- +goose StatementEnd
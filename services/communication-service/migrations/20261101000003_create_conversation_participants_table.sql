-- +goose Up
-- +goose StatementBegin

CREATE TABLE conversation_participants (
    conversation_id UUID NOT NULL,
    user_id UUID NOT NULL,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_read_at TIMESTAMPTZ NULL,

    PRIMARY KEY (conversation_id, user_id),

    CONSTRAINT fk_conversation_participants_conversation_id
        FOREIGN KEY (conversation_id)
        REFERENCES conversations(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_conversation_participants_user_id
        FOREIGN KEY (user_id)
        REFERENCES communication_users(user_id)
);

CREATE INDEX idx_conversation_participants_user_id ON conversation_participants(user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS conversation_participants;

-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin

CREATE TABLE message_read_receipts (
    message_id UUID NOT NULL,
    conversation_id UUID NOT NULL,
    user_id UUID NOT NULL,
    read_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (message_id, user_id),

    CONSTRAINT fk_message_read_receipts_message_id
        FOREIGN KEY (message_id)
        REFERENCES messages(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_message_read_receipts_conversation_id
        FOREIGN KEY (conversation_id)
        REFERENCES conversations(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_message_read_receipts_user_id
        FOREIGN KEY (user_id)
        REFERENCES communication_users(user_id)
);

CREATE INDEX idx_message_read_receipts_conversation_user
    ON message_read_receipts(conversation_id, user_id);

-- Preserve read state that was tracked by last_read_at before receipts existed.
INSERT INTO message_read_receipts (message_id, conversation_id, user_id, read_at)
SELECT m.id, m.conversation_id, cp.user_id, cp.last_read_at
FROM messages m
JOIN conversation_participants cp ON cp.conversation_id = m.conversation_id
WHERE cp.last_read_at IS NOT NULL
  AND m.sender_id <> cp.user_id
  AND m.created_at <= cp.last_read_at
ON CONFLICT (message_id, user_id) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS message_read_receipts;

-- +goose StatementEnd

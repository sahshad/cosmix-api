-- +goose Up
-- +goose StatementBegin

CREATE TABLE conversations (
    id UUID PRIMARY KEY,
    type VARCHAR(20) NOT NULL DEFAULT 'direct',
    direct_key VARCHAR(80) NULL,
    created_by UUID NOT NULL,
    last_message_id UUID NULL,
    last_message_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL,

    CONSTRAINT uq_conversations_direct_key UNIQUE(direct_key),

    CONSTRAINT fk_conversations_created_by
        FOREIGN KEY (created_by)
        REFERENCES communication_users(user_id)
);

CREATE INDEX idx_conversations_last_message_at ON conversations(last_message_at DESC NULLS LAST);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS conversations;

-- +goose StatementEnd

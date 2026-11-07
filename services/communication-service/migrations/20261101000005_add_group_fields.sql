-- +goose Up
-- +goose StatementBegin

ALTER TABLE conversations
    ADD COLUMN name VARCHAR(100) NULL,
    ADD COLUMN avatar_url TEXT NULL;

ALTER TABLE conversation_participants
    ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'member';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE conversation_participants
    DROP COLUMN IF EXISTS role;

ALTER TABLE conversations
    DROP COLUMN IF EXISTS avatar_url,
    DROP COLUMN IF EXISTS name;

-- +goose StatementEnd

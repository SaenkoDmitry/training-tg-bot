-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id                BIGSERIAL PRIMARY KEY,
    username          TEXT,
    chat_id           BIGINT,
    first_name        TEXT,
    last_name         TEXT,
    language_code     TEXT,
    active_program_id BIGINT,
    created_at        TIMESTAMP WITH TIME ZONE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN icon TEXT NOT NULL DEFAULT 'Smile';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN icon;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS training;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS training;
-- +goose StatementEnd

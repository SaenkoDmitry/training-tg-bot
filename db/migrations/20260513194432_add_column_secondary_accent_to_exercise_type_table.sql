-- +goose Up
-- +goose StatementBegin
ALTER TABLE exercise_types ADD COLUMN secondary_accent TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE exercise_types DROP COLUMN secondary_accent;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE training.exercise_group_types
(
    code TEXT UNIQUE,
    name TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE training.exercise_group_types;
-- +goose StatementEnd

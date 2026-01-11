-- +goose Up
-- +goose StatementBegin
CREATE TABLE exercise_group_types
(
    code TEXT UNIQUE,
    name TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE .exercise_group_types;
-- +goose StatementEnd

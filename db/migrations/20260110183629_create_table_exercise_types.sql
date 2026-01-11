-- +goose Up
-- +goose StatementBegin
CREATE TABLE training.exercise_types
(
    id                       BIGSERIAL PRIMARY KEY,
    name                     TEXT,
    url                      TEXT,
    exercise_group_type_code TEXT NOT NULL REFERENCES training.exercise_group_types (code),
    rest_in_seconds          INT,
    accent                   TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE training.exercise_types;
-- +goose StatementEnd

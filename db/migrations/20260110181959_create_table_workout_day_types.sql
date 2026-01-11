-- +goose Up
-- +goose StatementBegin
CREATE TABLE workout_day_types
(
    id                 BIGSERIAL PRIMARY KEY,
    workout_program_id BIGINT NOT NULL REFERENCES workout_programs (id),
    name               TEXT,
    preset             TEXT,
    created_at         TIMESTAMP WITH TIME ZONE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workout_day_types;
-- +goose StatementEnd

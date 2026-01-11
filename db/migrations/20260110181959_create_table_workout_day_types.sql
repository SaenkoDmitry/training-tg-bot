-- +goose Up
-- +goose StatementBegin
CREATE TABLE training.workout_day_types
(
    id                 BIGSERIAL PRIMARY KEY,
    workout_program_id BIGINT NOT NULL REFERENCES training.workout_programs (id),
    name               TEXT,
    preset             TEXT,
    created_at         TIMESTAMP WITH TIME ZONE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE training.workout_day_types;
-- +goose StatementEnd

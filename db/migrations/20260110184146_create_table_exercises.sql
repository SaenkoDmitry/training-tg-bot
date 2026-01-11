-- +goose Up
-- +goose StatementBegin
CREATE TABLE training.exercises
(
    id               BIGSERIAL PRIMARY KEY,
    workout_day_id   BIGINT NOT NULL REFERENCES training.workout_days (id) ON DELETE CASCADE,
    exercise_type_id BIGINT NOT NULL REFERENCES training.exercise_types (id),
    rest_in_seconds  INT,
    index            BIGINT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE training.exercises;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE exercises
(
    id               BIGSERIAL PRIMARY KEY,
    workout_day_id   BIGINT NOT NULL REFERENCES workout_days (id) ON DELETE CASCADE,
    exercise_type_id BIGINT NOT NULL REFERENCES exercise_types (id),
    rest_in_seconds  INT,
    index            BIGINT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE exercises;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE training.workout_sessions
(
    id                     BIGSERIAL PRIMARY KEY,
    workout_day_id         BIGINT NOT NULL REFERENCES training.workout_days (id) ON DELETE CASCADE,
    current_exercise_index BIGINT,
    started_at             TIMESTAMP WITH TIME ZONE,
    is_active              BOOLEAN
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE training.workout_sessions;
-- +goose StatementEnd

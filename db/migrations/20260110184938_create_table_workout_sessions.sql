-- +goose Up
-- +goose StatementBegin
CREATE TABLE workout_sessions
(
    id                     BIGSERIAL PRIMARY KEY,
    workout_day_id         BIGINT NOT NULL REFERENCES workout_days (id) ON DELETE CASCADE,
    current_exercise_index BIGINT,
    started_at             TIMESTAMP WITH TIME ZONE,
    is_active              BOOLEAN
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workout_sessions;
-- +goose StatementEnd

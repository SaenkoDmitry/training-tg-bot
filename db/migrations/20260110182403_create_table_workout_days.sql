-- +goose Up
-- +goose StatementBegin
CREATE TABLE workout_days
(
    id                   BIGSERIAL PRIMARY KEY,
    user_id              BIGINT,
    workout_day_type_id BIGINT NOT NULL REFERENCES workout_day_types (id),
    started_at           TIMESTAMP WITH TIME ZONE,
    ended_at             TIMESTAMP WITH TIME ZONE,
    completed            BOOLEAN
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workout_days;
-- +goose StatementEnd

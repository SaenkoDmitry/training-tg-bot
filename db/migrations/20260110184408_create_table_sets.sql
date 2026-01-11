-- +goose Up
-- +goose StatementBegin
CREATE TABLE training.sets
(
    id           BIGSERIAL PRIMARY KEY,
    exercise_id  BIGINT NOT NULL REFERENCES training.exercises (id) ON DELETE CASCADE,
    reps         BIGINT,
    fact_reps    BIGINT,
    weight       NUMERIC,
    fact_weight  NUMERIC,
    minutes      INT,
    fact_minutes INT,
    completed    BOOLEAN,
    completed_at TIMESTAMP WITH TIME ZONE,
    index        INT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE training.sets;
-- +goose StatementEnd

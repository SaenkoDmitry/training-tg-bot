-- +goose Up
-- +goose StatementBegin
CREATE TABLE workout_shares
(
    id             SERIAL PRIMARY KEY,
    workout_day_id BIGINT      NOT NULL REFERENCES workout_days (id) ON DELETE CASCADE,
    token          VARCHAR(64) NOT NULL UNIQUE,
    created_at     TIMESTAMP   NOT NULL DEFAULT NOW(),
    expires_at     TIMESTAMP,
    view_count     INT         NOT NULL DEFAULT 0
);

CREATE INDEX idx_workout_shares_token ON workout_shares (token);
CREATE INDEX idx_workout_shares_workout_day_id ON workout_shares (workout_day_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS workout_shares;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE rest_timers (
                             id BIGSERIAL PRIMARY KEY,
                             user_id BIGINT NOT NULL REFERENCES users (id),
                             workout_id BIGINT NOT NULL REFERENCES workout_days (id),
                             ends_at TIMESTAMP NOT NULL,
                             canceled BOOLEAN DEFAULT FALSE,
                             sent BOOLEAN DEFAULT FALSE,
                             created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_rest_timers_user ON rest_timers(user_id);
CREATE INDEX idx_rest_timers_ends_at ON rest_timers(ends_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rest_timers;
-- +goose StatementEnd

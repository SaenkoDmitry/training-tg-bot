-- +goose Up
-- +goose StatementBegin
CREATE TABLE push_subscriptions
(
    id         SERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users (id),
    endpoint   TEXT   NOT NULL UNIQUE,
    p256dh     TEXT   NOT NULL,
    auth       TEXT   NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE push_subscriptions;
-- +goose StatementEnd

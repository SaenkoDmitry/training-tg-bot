-- +goose Up
-- +goose StatementBegin
CREATE TABLE training.workout_programs
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES training.users (id) ON DELETE CASCADE,
    name       TEXT,
    created_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE training.users
    ADD CONSTRAINT users_workout_programs_foreign_key
        FOREIGN KEY (active_program_id) REFERENCES training.workout_programs (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE training.users DROP CONSTRAINT users_workout_programs_foreign_key;

DROP TABLE training.workout_programs;
-- +goose StatementEnd

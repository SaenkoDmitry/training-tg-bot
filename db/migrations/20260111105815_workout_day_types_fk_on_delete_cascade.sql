-- +goose Up
-- +goose StatementBegin
ALTER TABLE workout_day_types DROP CONSTRAINT workout_day_types_workout_program_id_fkey;
ALTER TABLE workout_day_types
    ADD CONSTRAINT workout_day_types_workout_program_id_fkey
        FOREIGN KEY (workout_program_id) REFERENCES workout_programs (id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workout_day_types DROP CONSTRAINT workout_day_types_workout_program_id_fkey;
ALTER TABLE workout_day_types
    ADD CONSTRAINT workout_day_types_workout_program_id_fkey
        FOREIGN KEY (workout_program_id) REFERENCES workout_programs (id);
-- +goose StatementEnd

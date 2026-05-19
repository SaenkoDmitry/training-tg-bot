-- +goose Up
-- +goose StatementBegin
ALTER TABLE workout_programs
    ADD COLUMN summary TEXT;
ALTER TABLE workout_programs
    ADD COLUMN warnings TEXT[];
ALTER TABLE workout_programs
    ADD COLUMN validation_notes TEXT[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workout_programs DROP COLUMN validation_notes;
ALTER TABLE workout_programs DROP COLUMN warnings;
ALTER TABLE workout_programs DROP COLUMN summary;
-- +goose StatementEnd

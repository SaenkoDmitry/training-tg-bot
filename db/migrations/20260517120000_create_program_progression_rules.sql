-- +goose Up
-- +goose StatementBegin
CREATE TABLE workout_program_progression_rules
(
    id                 BIGSERIAL PRIMARY KEY,
    workout_program_id BIGINT NOT NULL REFERENCES workout_programs (id) ON DELETE CASCADE,
    workout_day_type_id BIGINT REFERENCES workout_day_types (id) ON DELETE CASCADE,
    exercise_type_id   BIGINT REFERENCES exercise_types (id),
    rule               TEXT NOT NULL,
    reason             TEXT,
    source             VARCHAR(32) NOT NULL DEFAULT 'ai',
    created_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workout_program_progression_rules_program ON workout_program_progression_rules (workout_program_id);
CREATE INDEX idx_workout_program_progression_rules_day ON workout_program_progression_rules (workout_day_type_id);
CREATE INDEX idx_workout_program_progression_rules_exercise_type ON workout_program_progression_rules (exercise_type_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS workout_program_progression_rules;
-- +goose StatementEnd

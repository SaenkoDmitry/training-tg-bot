-- +goose Up
-- +goose StatementBegin
ALTER TABLE workout_days
    ADD COLUMN estimated_calories NUMERIC(8, 2),
    ADD COLUMN estimated_duration_minutes INTEGER,  -- полезно для статистики
    ADD COLUMN user_weight_kg NUMERIC(5,2); -- вес на момент тренировки

alter table exercises
    add column estimated_calories numeric(8, 2),
    add column estimated_duration_seconds integer;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workout_days DROP COLUMN estimated_calories;
ALTER TABLE workout_days DROM COLUMN estimated_duration_minutes;
ALTER TABLE workout_days DROM COLUMN user_weight_kg;

ALTER TABLE exercises DROM COLUMN estimated_calories;
ALTER TABLE exercises DROM COLUMN estimated_duration_seconds;
-- +goose StatementEnd

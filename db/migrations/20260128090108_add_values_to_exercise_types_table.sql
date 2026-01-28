-- +goose Up
-- +goose StatementBegin
INSERT INTO exercise_types (name, url, exercise_group_type_code, rest_in_seconds, accent, units)
VALUES ('Вертикальная тяга параллельным узким хватом', 'https://disk.yandex.ru/i/Sal03-DEMfRKAQ', 'back', 120,
        'нижняя часть широчайших мышц и середины спины', 'reps,weight');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM exercise_types where name = 'Вертикальная тяга параллельным узким хватом';
-- +goose StatementEnd

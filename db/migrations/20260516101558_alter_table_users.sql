-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN birth_date DATE,
    ADD COLUMN gender VARCHAR(10) CHECK (gender IN ('male', 'female')),
    ADD COLUMN weight_kg NUMERIC(5,2),
    ADD COLUMN height_cm INTEGER;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN birth_date;
ALTER TABLE users DROM COLUMN gender;
ALTER TABLE users DROM COLUMN weight_kg;
ALTER TABLE users DROM COLUMN height_cm;
-- +goose StatementEnd

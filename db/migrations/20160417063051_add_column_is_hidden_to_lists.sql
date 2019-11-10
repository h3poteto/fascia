
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE lists ADD COLUMN is_hidden boolean NOT NULL DEFAULT FALSE;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE lists DROP COLUMN is_hidden;



-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE tasks ADD COLUMN display_index int NOT NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE tasks DROP COLUMN display_index;

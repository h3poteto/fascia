
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE lists ADD COLUMN is_archived boolean NOT NULL DEFAULT FALSE AFTER list_option_id;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE lists DROP COLUMN is_archived;

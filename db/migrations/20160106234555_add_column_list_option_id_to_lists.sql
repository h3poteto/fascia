
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE lists ADD COLUMN list_option_id int DEFAULT NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE lists DROP COLUMN list_option_id;

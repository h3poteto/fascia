
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE lists ADD COLUMN user_id int NOT NULL;
ALTER TABLE tasks ADD COLUMN user_id int NOT NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE lists DROP COLUMN user_id;
ALTER TABLE tasks DROP COLUMN user_id;

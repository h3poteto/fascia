
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE lists ADD COLUMN user_id int(11) NOT NULL AFTER project_id;
ALTER TABLE tasks ADD COLUMN user_id int(11) NOT NULL AFTER list_id;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE lists DROP COLUMN user_id;
ALTER TABLE tasks DROP COLUMN user_id;


-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE tasks ADD COLUMN issue_number int(11) DEFAULT NULL AFTER user_id;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE tasks DROP COLUMN issue_number;

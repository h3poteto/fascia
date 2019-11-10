
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE repositories ADD COLUMN webhook_key varchar(255) NOT NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE repositories DROP COLUMN webhook_key;

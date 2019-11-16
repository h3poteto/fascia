
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE tasks ADD COLUMN pull_request boolean NOT NULL DEFAULT FALSE;
ALTER TABLE tasks ADD COLUMN html_url varchar(255) DEFAULT NULL;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE tasks DROP COLUMN pull_request;
ALTER TABLE tasks DROP COLUMN html_url;

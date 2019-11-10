
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE projects ADD COLUMN show_issues boolean NOT NULL DEFAULT TRUE;
ALTER TABLE projects ADD COLUMN show_pull_requests boolean NOT NULL DEFAULT TRUE;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE projects DROP COLUMN show_issues;
ALTER TABLE projects DROP COLUMN show_pull_requests;

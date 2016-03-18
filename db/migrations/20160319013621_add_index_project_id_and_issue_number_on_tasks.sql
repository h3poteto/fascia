
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE tasks ADD UNIQUE INDEX index_on_project_id_and_issue_number(project_id, issue_number);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE tasks DROP INDEX index_on_project_id_and_issue_number;

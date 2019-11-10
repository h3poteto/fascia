
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE UNIQUE INDEX project_id_and_issue_number_on_tasks on tasks (project_id, issue_number);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX project_id_and_issue_number_on_tasks;

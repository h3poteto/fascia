
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE repositories ADD UNIQUE index_on_repository_id(repository_id);

ALTER TABLE projects ADD UNIQUE index_on_title_and_user_id(title, user_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE repositories DROP INDEX index_on_repository_id;

ALTER TABLE projects DROP INDEX index_on_title_and_user_id;

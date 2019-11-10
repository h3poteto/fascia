
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE UNIQUE INDEX repository_id_on_repositories on repositories (repository_id);

CREATE UNIQUE INDEX title_and_user_id_on_projects on projects (title, user_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX repository_id_on_repositories;
DROP INDEX title_and_user_id_on_projects;

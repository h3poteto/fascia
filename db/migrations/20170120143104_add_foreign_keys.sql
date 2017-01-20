
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE projects ADD CONSTRAINT projects_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE projects ADD CONSTRAINT projects_repository_id_fkey FOREIGN KEY (repository_id) REFERENCES repositories(id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE reset_passwords ADD CONSTRAINT reset_passwords_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE lists ADD CONSTRAINT lists_project_id_fkey FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE lists ADD CONSTRAINT lists_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE lists ADD CONSTRAINT lists_list_option_id_fkey FOREIGN KEY (list_option_id) REFERENCES list_options(id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE tasks ADD CONSTRAINT tasks_list_id_fkey FOREIGN KEY (list_id) REFERENCES lists(id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE tasks ADD CONSTRAINT tasks_project_id_fkey FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE tasks ADD CONSTRAINT tasks_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE RESTRICT;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE projects DROP FOREIGN KEY projects_user_id_fkey;
ALTER TABLE projects DROP FOREIGN KEY projects_repository_id_fkey;
ALTER TABLE reset_passwords DROP FOREIGN KEY reset_passwords_user_id_fkey;
ALTER TABLE lists DROP FOREIGN KEY lists_project_id_fkey;
ALTER TABLE lists DROP FOREIGN KEY lists_user_id_fkey;
ALTER TABLE lists DROP FOREIGN KEY lists_list_option_id_fkey;
ALTER TABLE tasks DROP FOREIGN KEY tasks_list_id_fkey;
ALTER TABLE tasks DROP FOREIGN KEY tasks_project_id_fkey;
ALTER TABLE tasks DROP FOREIGN KEY tasks_user_id_fkey;

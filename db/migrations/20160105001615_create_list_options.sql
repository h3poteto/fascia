
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS list_options (
id SERIAL PRIMARY KEY,
action varchar(255) NOT NULL,
created_at timestamp NOT NULL DEFAULT current_timestamp,
updated_at timestamp NOT NULL DEFAULT current_timestamp);

CREATE UNIQUE INDEX action_on_list_options on list_options (action);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX action_on_list_options;
DROP TABLE list_options;

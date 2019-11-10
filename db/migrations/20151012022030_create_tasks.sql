
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS tasks (
id SERIAL PRIMARY KEY,
list_id int NOT NULL,
title varchar(255) NOT NULL,
created_at timestamp NOT NULL DEFAULT current_timestamp,
updated_at timestamp NOT NULL DEFAULT current_timestamp);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE tasks;

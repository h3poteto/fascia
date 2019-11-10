
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS lists (
id SERIAL PRIMARY KEY,
project_id int NOT NULL,
title varchar(255) DEFAULT NULL,
created_at timestamp NOT NULL DEFAULT current_timestamp,
updated_at timestamp NOT NULL DEFAULT current_timestamp);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE lists;


-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS inquiries (
id SERIAL PRIMARY KEY,
email varchar(255) NOT NULL,
name varchar(255) NOT NULL,
message text DEFAULT NULL,
created_at timestamp NOT NULL DEFAULT current_timestamp,
updated_at timestamp NOT NULL DEFAULT current_timestamp);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE inquiries;


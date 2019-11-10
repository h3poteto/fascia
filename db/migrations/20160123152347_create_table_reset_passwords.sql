
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS reset_passwords (
id SERIAL PRIMARY KEY,
user_id int NOT NULL,
token varchar(255) NOT NULL,
expires_at timestamp DEFAULT NULL,
created_at timestamp NOT NULL DEFAULT current_timestamp,
updated_at timestamp NOT NULL DEFAULT current_timestamp);

CREATE UNIQUE INDEX token_on_reset_passwords on reset_passwords (token);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX token_on_reset_password;
DROP TABLE reset_passwords;

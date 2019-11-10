
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS users (
id SERIAL PRIMARY KEY,
email varchar(255) NOT NULL,
password varchar(255) NOT NULL,
provider varchar(255) DEFAULT NULL,
oauth_token varchar(255) DEFAULT NULL,
uuid bigint DEFAULT NULL,
user_name varchar(255) DEFAULT NULL,
avatar_url varchar(255) DEFAULT NULL,
created_at timestamp NOT NULL DEFAULT current_timestamp,
updated_at timestamp NOT NULL DEFAULT current_timestamp);

CREATE UNIQUE INDEX email_on_users ON users (email);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX email_on_users;
DROP TABLE users;

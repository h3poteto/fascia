
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS reset_passwords (
id int(11) NOT NULL AUTO_INCREMENT,
user_id int(11) NOT NULL,
token varchar(255) NOT NULL DEFAULT "",
expires_at datetime DEFAULT NULL,
created_at datetime DEFAULT NULL,
updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
PRIMARY KEY (id),
UNIQUE INDEX index_on_token (token))
ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE reset_passwords;

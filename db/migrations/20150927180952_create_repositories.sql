
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS repositories (
id int(11) NOT NULL AUTO_INCREMENT,
project_id int(11) NOT NULL,
repository_id int(11) NOT NULL,
full_name varchar(255) DEFAULT NULL,
created_at datetime DEFAULT NULL,
updated_At timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
PRIMARY KEY (id))
ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE repositories;

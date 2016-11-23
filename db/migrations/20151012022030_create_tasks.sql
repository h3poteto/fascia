
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS tasks (
id int(11) NOT NULL AUTO_INCREMENT,
list_id int(11) NOT NULL,
title varchar(255) NOT NULL DEFAULT "",
created_at datetime DEFAULT NULL,
updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
PRIMARY KEY(id))
ENGINE=Innodb AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE tasks;

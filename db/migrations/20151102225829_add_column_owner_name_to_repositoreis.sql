
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE repositories DROP COLUMN full_name;
ALTER TABLE repositories ADD COLUMN owner varchar(255) DEFAULT NULL AFTER repository_id;
ALTER TABLE repositories ADD COLUMN name  varchar(255) DEFAULT NULL AFTER owner;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE repositories ADD COLUMN full_name varchar(255) DEFAULT NULL AFTER repository_id;
ALTER TABLE repositories DROP COLUMN owner;
ALTER TABLE repositories DROP COLUMN name;

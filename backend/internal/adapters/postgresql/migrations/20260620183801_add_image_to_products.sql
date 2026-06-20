-- +goose Up
ALTER TABLE products ADD COLUMN image_url VARCHAR(255);

-- +goose Down
ALTER TABLE products DROP COLUMN image_url;

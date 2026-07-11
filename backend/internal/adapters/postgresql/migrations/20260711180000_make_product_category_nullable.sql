-- +goose Up
ALTER TABLE products ALTER COLUMN category_id DROP NOT NULL;

-- +goose Down
ALTER TABLE products ALTER COLUMN category_id SET NOT NULL;

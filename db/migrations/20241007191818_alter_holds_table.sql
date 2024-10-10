-- +goose Up
ALTER TABLE holds
	ADD COLUMN image_url VARCHAR(255) NOT NULL DEFAULT '';

-- +goose Down
DROP TABLE holds;

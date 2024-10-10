-- +goose Up
ALTER TABLE holds
	ADD COLUMN image_width INT NOT NULL DEFAULT 24,
	ADD COLUMN image_length INT NOT NULL DEFAULT 24;

-- +goose Down
DROP TABLE holds;

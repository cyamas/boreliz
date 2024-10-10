-- +goose Up
ALTER TABLE holds
	ADD COLUMN angle INT NOT NULL DEFAULT 0;
-- +goose Down
DROP table holds;

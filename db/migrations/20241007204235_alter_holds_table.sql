-- +goose Up
ALTER TABLE holds
	ADD COLUMN texture INT NOT NULL DEFAULT 0;
ALTER TABLE hold_edges
	DROP column texture;
-- +goose Down
DROP TABLE holds;

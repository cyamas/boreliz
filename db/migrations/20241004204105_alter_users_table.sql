-- +goose Up
ALTER TABLE users
    ALTER COLUMN vert_reach SET DATA TYPE FLOAT,
    ALTER COLUMN vert_reach SET DEFAULT 0.0,
    ALTER COLUMN height SET DATA TYPE FLOAT,
    ALTER COLUMN height SET DEFAULT 0.0,
    ALTER COLUMN wingspan SET DATA TYPE FLOAT,
    ALTER COLUMN wingspan SET DEFAULT 0.0;

-- +goose Down
DROP TABLE users;-- +goose Up

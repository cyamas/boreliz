-- +goose Up
ALTER TABLE users
ADD vert_reach DECIMAL(5,2) NOT NULL DEFAULT 0.00,
ALTER COLUMN height TYPE DECIMAL(5, 2),
ALTER COLUMN wingspan TYPE DECIMAL(5, 2);

-- +goose Down
DROP TABLE users;

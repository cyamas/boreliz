-- +goose Up
CREATE TABLE walls (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL DEFAULT 'unknown',
	rows INT NOT NULL DEFAULT 0,
	cols INT NOT NULL DEFAULT 0,
	spacing INT NOT NULL DEFAULT 0,
	is_adjustable BOOLEAN DEFAULT FALSE, 
	angle INT NOT NULL DEFAULT 90
);
CREATE TABLE holds (
	id INT PRIMARY KEY,
	manufacturer VARCHAR(255) NOT NULL DEFAULT 'unknown',
	model VARCHAR(255) NOT NULL DEFAULT 'unknown',
	type VARCHAR(255) NOT NULL DEFAULT 'unknown',
	color VARCHAR(255) NOT NULL DEFAULT 'unknown',
	wall_id int NOT NULL DEFAULT 0,
	row int NOT NULL DEFAULT -1,
	col int NOT NULL DEFAULT -1,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (wall_id) REFERENCES walls(id)
);
CREATE TABLE hold_edges (
	id SERIAL PRIMARY KEY,
	hold_id INT NOT NULL,
	angle INT NOT NULL DEFAULT 0,
	width INT NOT NULL DEFAULT 0,
	incut INT NOT NULL DEFAULT 0,
	depth DECIMAL(4, 3) NOT NULL DEFAULT 0.000,
	texture INT NOT NULL DEFAULT 0,
	FOREIGN KEY (hold_id) REFERENCES holds(id)
);


-- +goose Down
DROP TABLE holds;


-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE clips (
	id SERIAL PRIMARY KEY, 
	name VARCHAR(100) NOT NULL,
	url VARCHAR(100) NOT NULL,
	slug VARCHAR(100) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by VARCHAR(100)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE clips;


-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE clips ADD COLUMN is_ready BOOLEAN NOT NULL DEFAULT TRUE;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE clips DROP COLUMN is_ready;

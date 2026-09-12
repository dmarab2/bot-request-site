-- +goose Up
ALTER TABLE tag_aliases
ALTER COLUMN tag_id
DROP NOT NULL;


-- +goose Down
ALTER TABLE tag_aliases
ALTER COLUMN tag_id
SET NOT NULL;
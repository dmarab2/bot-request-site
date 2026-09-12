-- +goose Up
CREATE TYPE sandbox_status AS ENUM ('pending', 'approved');
ALTER TABLE tag_aliases
ADD alias_status sandbox_status NOT NULL DEFAULT 'pending';

-- +goose Down
ALTER TABLE tag_aliases
DROP COLUMN alias_status;
DROP TYPE sandbox_status;
-- +goose Up
CREATE TYPE request_status AS ENUM ('open', 'in_progress', 'fulfilled', 'cancelled');
ALTER TABLE requests
ADD status request_status NOT NULL DEFAULT 'open';

-- +goose Down
ALTER TABLE requests
DROP COLUMN status;
DROP TYPE request_status;
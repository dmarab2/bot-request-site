-- name: ChangeRequestStatus :one
UPDATE requests
SET status = $1
WHERE id = $2
RETURNING *;
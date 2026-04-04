-- name: GetAllRequestsFiltered :many
SELECT * FROM requests
WHERE status = $1;
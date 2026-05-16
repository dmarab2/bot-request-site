-- name: GetSingleRequest :one
SELECT * FROM requests
WHERE id = $1;
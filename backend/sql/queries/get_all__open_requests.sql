-- name: GetAllOpenRequests :many
SELECT * FROM requests
WHERE status = 'open';
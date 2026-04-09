-- name: GetNextRequestPage :many
SELECT * FROM requests
WHERE status = $1
AND id <= $2
ORDER BY created_at DESC
LIMIT 5;
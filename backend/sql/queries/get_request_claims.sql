-- name: GetRequestClaims :many
SELECT * FROM request_claims
WHERE expires_at IS NOT NULL
AND request_id <= $1
ORDER BY claimed_at DESC
LIMIT 5;
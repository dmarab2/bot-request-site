-- name: GetFirstPageCursor :one
SELECT id FROM requests ORDER BY created_at DESC LIMIT 1;
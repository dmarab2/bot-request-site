-- name: DoesTagExist :one
SELECT COUNT(*)
FROM tags
WHERE name = $1;
-- name: DoesTagExist :one
SELECT COUNT(*)
FROM tags
WHERE tag_name = $1;
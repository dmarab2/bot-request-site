-- name: DeleteTagById :exec
DELETE FROM tags
WHERE id = $1;
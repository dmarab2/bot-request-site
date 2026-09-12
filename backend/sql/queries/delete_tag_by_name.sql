-- name: DeleteTagByName :exec
DELETE FROM tags
WHERE name = $1;
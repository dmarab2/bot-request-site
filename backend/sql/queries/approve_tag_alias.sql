-- name: ApproveTagAlias :one
UPDATE tag_aliases
SET alias_status = 'approved'
WHERE id = $1
RETURNING *;
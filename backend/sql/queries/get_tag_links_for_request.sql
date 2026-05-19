-- name: GetTagLinksForRequest :many

SELECT *
FROM tags
JOIN request_tags ON tags.id = request_tags.tag_id
WHERE request_tags.request_id = $1;
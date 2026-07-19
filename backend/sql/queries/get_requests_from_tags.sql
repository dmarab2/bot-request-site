-- name: GetRequestsFromTags :many
SELECT *
FROM requests
WHERE $1 <@ (
    SELECT coalesce(array_agg(tags.name), '{}')
    FROM request_tags
    LEFT JOIN tags ON tags.id = request_tags.tag_id
    WHERE requests.id = request_tags.request_id
)
LIMIT 10;
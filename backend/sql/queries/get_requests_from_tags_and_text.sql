-- name: GetRequestsFromTagsAndText :many
SELECT *
FROM requests
WHERE $1::text[] <@ (
    SELECT coalesce(array_agg(tags.name), '{}')
    FROM request_tags
    LEFT JOIN tags ON tags.id = request_tags.tag_id
    WHERE requests.id = request_tags.request_id
) AND request_search_vector @@ websearch_to_tsquery('english', $2)
ORDER BY ts_rank(request_search_vector, websearch_to_tsquery('english', $2)) DESC
LIMIT 10;
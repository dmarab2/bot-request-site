-- name: ChangeFulfilledClaimExpireTimes :exec
UPDATE request_claims AS claim
SET expires_at = NOW()
FROM requests AS request
WHERE request.id = claim.request_id AND request.status IN ('fulfilled', 'cancelled');
-- name: CreateOpportunity :one
INSERT INTO opportunities (name, description, stage, amount, close_date, probability, lead_id, account_id, owner_id)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
RETURNING *;

-- name: GetOpportunity :one
SELECT * FROM opportunities WHERE id = $1;

-- name: ListOpportunities :many
SELECT * FROM opportunities
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateOpportunity :one
UPDATE opportunities
SET stage=$2, amount=$3, probability=$4, updated_at=CURRENT_TIMESTAMP
WHERE id=$1
RETURNING *;

-- name: DeleteOpportunity :exec
DELETE FROM opportunities WHERE id = $1;

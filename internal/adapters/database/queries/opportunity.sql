-- name: CreateOpportunity :one
INSERT INTO opportunities (name, description, stage, amount, close_date, probability, lead_id, account_id, owner_id)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
RETURNING *;

-- name: GetOpportunity :one
SELECT * FROM opportunities WHERE id = $1 LIMIT 1;

-- name: ListOpportunities :many
SELECT *
FROM opportunities
WHERE ($1::int = 0 OR owner_id = $1)
ORDER BY created_at DESC;

-- name: UpdateOpportunity :one
UPDATE opportunities
SET stage=$2, amount=$3, probability=$4, updated_at=CURRENT_TIMESTAMP
WHERE id=$1
RETURNING *;

-- name: UpdateOpportunitySelective :one
UPDATE opportunities
SET
  name        = COALESCE(sqlc.narg(name), name),
  description = COALESCE(sqlc.narg(description), description),
  stage       = COALESCE(sqlc.narg(stage), stage),
  amount      = COALESCE(sqlc.narg(amount), amount),
  close_date  = COALESCE(sqlc.narg(close_date), close_date),
  probability = COALESCE(sqlc.narg(probability), probability),
  lead_id     = COALESCE(sqlc.narg(lead_id), lead_id),
  account_id  = COALESCE(sqlc.narg(account_id), account_id),
  owner_id    = COALESCE(sqlc.narg(owner_id), owner_id),
  updated_at  = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteOpportunity :exec
DELETE FROM opportunities WHERE id = $1;

-- name: CreateLead :one
INSERT INTO leads (first_name, last_name, email, phone, status, assigned_to, organization_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING.*;

-- name: GetLeadById :one
SELECT * FROM leads WHERE id = $1;

-- name: GetLeadByEmail :one
SELECT * FROM leads WHERE email = $1 LIMIT 1;

-- name: GetAll :many
SELECT * FROM leads
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateLead :one
UPDATE leads
SET status=$2, assigned_to=$3, updated_at=CURRENT_TIMESTAMP
WHERE id=$1
RETURNING *;

-- name: DeleteLead :exec
DELETE FROM leads WHERE id = $1;
